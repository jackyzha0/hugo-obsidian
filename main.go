package main

import (
	"flag"
	"fmt"
	"github.com/gernest/front"
	md "github.com/nikitavoloboev/markdown-parser"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type Link struct {
	Source string
	Target string
	Text   string
}

type LinkTable = map[string][]Link
type Index struct {
	Links LinkTable
	Backlinks LinkTable
}

type Content struct {
	Title string
	Content string
}

type ContentIndex = map[string]Content

func trim(source, prefix, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(source, suffix), prefix)
}

func hugoPathTrim(source string) string {
	return strings.TrimSuffix(strings.TrimSuffix(source, "/index"), "_index")
}

func processTarget(source string) string {
	if !isInternal(source) {
		return source
	}
	if strings.HasPrefix(source, "/") {
		return strings.TrimSuffix(source, ".md")
	}
	return "/" + strings.TrimSuffix(source, ".md")
}

func isInternal(link string) bool {
	return !strings.HasPrefix(link, "http")
}

// parse single file for links
func parse(dir, pathPrefix string) []Link {
	// read file
	bytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("%s\n", trim(dir, pathPrefix, ".md"))
	for text, target := range md.GetAllLinks(string(bytes)) {
		target := strings.Split(processTarget(target), "#")[0]
		fmt.Printf("  %s\n", target)
		links = append(links, Link{
			Source: filepath.ToSlash(hugoPathTrim(trim(dir, pathPrefix, ".md"))),
			Target: target,
			Text: text,
		})
	}
	return links
}

func getText(dir string) string {
	// read file
	bytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// recursively walk directory and return all files with given extension
func walk(root, ext string, index bool) (res []Link, i ContentIndex) {
	println(root)
	i = make(ContentIndex)

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	nPrivate := 0

	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			res = append(res, parse(s, root)...)
			if index {
				text := getText(s)

				frontmatter, body, err := m.Parse(strings.NewReader(text))
				if err != nil {
					frontmatter = map[string]interface{}{}
					body = text
				}

				var title string
				if parsedTitle, ok := frontmatter["title"]; ok {
					title = parsedTitle.(string)
				} else {
					title = "Untitled Page"
				}

				// check if page is private
				if parsedPrivate, ok := frontmatter["draft"]; !ok || !parsedPrivate.(bool) {
					adjustedPath := hugoPathTrim(trim(s, root, ".md"))
					i[adjustedPath] = Content{
						Title: title,
						Content: body,
					}
				} else {
					nPrivate++
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ignored %d private files \n", nPrivate)
	fmt.Printf("Parsed %d total links \n", len(res))
	return res, i
}

// filter out certain links (e.g. to media)
func filter(links []Link) (res []Link) {
	for _, l := range links {
		// filter external and non-md
		isMarkdown := filepath.Ext(l.Target) == "" || filepath.Ext(l.Target) == ".md"
		if isInternal(l.Target) && isMarkdown {
			res = append(res, l)
		}
	}
	fmt.Printf("Removed %d external and non-markdown links\n", len(links) - len(res))
	return res
}

// constructs index from links
func index(links []Link) (index Index) {
	linkMap := make(map[string][]Link)
	backlinkMap := make(map[string][]Link)
	for _, l := range links {
		// backlink (only if internal)
		if _, ok := backlinkMap[l.Target]; ok {
			backlinkMap[l.Target] = append(backlinkMap[l.Target], l)
		} else {
			backlinkMap[l.Target] = []Link{l}
		}

		// regular link
		if _, ok := linkMap[l.Source]; ok {
			linkMap[l.Source] = append(linkMap[l.Source], l)
		} else {
			linkMap[l.Source] = []Link{l}
		}
	}
	index.Links = linkMap
	index.Backlinks = backlinkMap
	return index
}

const message = "# THIS FILE WAS GENERATED USING github.com/jackyzha0/hugo-obsidian\n# DO NOT EDIT\n"
func write(links []Link, contentIndex ContentIndex, toIndex bool, out string) error {
	index := index(links)
	resStruct := struct{
		Index Index
		Links []Link
	}{
		Index: index,
		Links: links,
	}
	marshalledIndex, mErr := yaml.Marshal(&resStruct)
	if mErr != nil {
		return mErr
	}

	writeErr := ioutil.WriteFile(path.Join(out, "linkIndex.yaml"), append([]byte(message), marshalledIndex...), 0644)
	if writeErr != nil {
		return writeErr
	}

	if toIndex {
		marshalledContentIndex, mcErr := yaml.Marshal(&contentIndex)
		if mcErr != nil {
			return mcErr
		}

		writeErr = ioutil.WriteFile(path.Join(out, "contentIndex.yaml"), append([]byte(message), marshalledContentIndex...), 0644)
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}

func main() {
	in := flag.String("input", ".", "Input Directory")
	out := flag.String("output", ".", "Output Directory")
	index := flag.Bool("index", false, "Whether to index the content")
	flag.Parse()
	l, i := walk(*in, ".md", *index)
	f := filter(l)
	err := write(f, i, *index, *out)
	if err != nil {
		panic(err)
	}
}
