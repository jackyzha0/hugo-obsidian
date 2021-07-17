package main

import (
	"fmt"
	md "github.com/nikitavoloboev/markdown-parser"
	"io/fs"
	"io/ioutil"
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

func trim(source, prefix, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(source, suffix), prefix)
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
	//fmt.Printf("%s \n", trim(dir, pathPrefix, ".md"))
	for text, target := range md.GetAllLinks(string(bytes)) {
		//fmt.Printf("  %s -> %s \n", text, target)
		links = append(links, Link{
			Source: trim(dir, pathPrefix, ".md"),
			Target: target,
			Text: text,
		})
	}
	return links
}

// recursively walk directory and return all files with given extension
func walk(root, ext string) (res []Link) {
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			res = append(res, parse(s, root)...)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return res
}

// constructs index from links
func index(links []Link) (index Index) {
	linkMap := make(map[string][]Link)
	backlinkMap := make(map[string][]Link)
	for _, l := range links {
		bl := Link{
			Source: l.Target,
			Target: l.Source,
			Text: l.Text,
		}

		// map has link
		if val, ok := backlinkMap[l.Target]; ok {
			val = append(val, bl)
		} else {
			backlinkMap[l.Target] = []Link{bl}
		}

		if val, ok := linkMap[l.Source]; ok {
			val = append(val, l)
		} else {
			linkMap[l.Target] = []Link{l}
		}
	}
	index.Links = linkMap
	index.Backlinks = backlinkMap
	return index
}

func main() {
	l := walk("../www/content", ".md")
	i := index(l)
	fmt.Printf("%+v", i.Links["/toc/cognitive-sciences"])
	fmt.Printf("%+v", i.Backlinks["/toc/cognitive-sciences"])
}
