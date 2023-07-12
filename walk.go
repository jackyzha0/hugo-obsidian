package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	wikilink "github.com/abhinav/goldmark-wikilink"
	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

type Front struct {
	Title string   `yaml:"title"`
	Draft bool     `yaml:"draft"`
	Tags  []string `yaml:"tags"`
}

type fileIndex struct {
	index map[string]string
}

// resolve takes a link path and attempts to canonicalize it according to this
// fileIndex. If a matching canonical link cannot be found, the link is returned
// untouched.
func (i *fileIndex) resolve(path string) string {
	if !isInternal(path) {
		return path
	}

	trimmedPath := strings.TrimLeft(path, "/")

	// If the path has any degree of nesting built into it, treat it as an absolute path
	if strings.Contains(trimmedPath, "/") {
		return path
	}

	resolved, ok := i.index[trimmedPath]
	if ok {
		target := processTarget(resolved)
		return target
	}

	return path
}

func buildFileIndex(root, ext string, ignorePaths map[string]struct{}) (fileIndex, error) {
	index := map[string]string{}
	err := filepath.WalkDir(root, func(fp string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		// path normalize fp
		s := filepath.ToSlash(fp)
		s = strings.ReplaceAll(s, " ", "-")
		if _, ignored := ignorePaths[s]; ignored {
			return nil
		} else if filepath.Ext(d.Name()) == ext {
			base := filepath.Base(strings.TrimSuffix(s, ".md"))
			index[base] = strings.TrimPrefix(s, root)
		}

		return nil
	})
	if err != nil {
		return fileIndex{}, err
	}

	return fileIndex{index: index}, nil
}

// recursively walk directory and return all files with given extension
func walk(root, ext string, index bool, ignorePaths map[string]struct{}) (res []Link, i ContentIndex) {
	fmt.Printf("Scraping %s\n", root)
	i = make(ContentIndex)

	nPrivate := 0

	formats := []*frontmatter.Format{
		frontmatter.NewFormat("---", "---", yaml.Unmarshal),
	}

	start := time.Now()

	md := goldmark.New(
		goldmark.WithExtensions(&wikilink.Extender{}))

	fileIndex, err := buildFileIndex(root, ext, ignorePaths)
	if err != nil {
		panic(err)
	}

	err = filepath.WalkDir(root, func(fp string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		// path normalize fp
		s := filepath.ToSlash(fp)
		if _, ignored := ignorePaths[s]; ignored {
			fmt.Printf("[Ignored] %s\n", d.Name())
			nPrivate++
		} else if filepath.Ext(d.Name()) == ext {
			if index {
				text := getText(s)

				var matter Front
				raw_body, err := frontmatter.Parse(strings.NewReader(text), &matter, formats...)
				body := string(raw_body)
				if err != nil {
					matter = Front{
						Title: "Untitled Page",
						Draft: false,
						Tags:  []string{},
					}
					body = text
				}
				// check if page is private
				if !matter.Draft {
					info, _ := os.Stat(s)
					source := processSource(trim(s, root, ".md"))

					// default title
					title := matter.Title
					if title == "" {
						fileName := d.Name()
						title = strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
					}

					// default tags
					if matter.Tags == nil {
						matter.Tags = []string{}
					}

					// add to content and link index
					i[source] = Content{
						LastModified: info.ModTime(),
						Title:        title,
						Content:      body,
						Tags:         matter.Tags,
					}
					res = append(res, parse(md, fileIndex, s, root)...)
				} else {
					fmt.Printf("[Ignored] %s\n", d.Name())
					nPrivate++
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	end := time.Now()

	fmt.Printf("[DONE] in %s\n", end.Sub(start).Round(time.Millisecond))
	fmt.Printf("Ignored %d private files \n", nPrivate)
	fmt.Printf("Parsed %d total links \n", len(res))
	return res, i
}

func getText(dir string) string {
	// read file
	fileBytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	return string(fileBytes)
}
