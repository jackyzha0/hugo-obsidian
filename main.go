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
	fmt.Printf("%s \n", trim(dir, pathPrefix, ".md"))
	for text, target := range md.GetAllLinks(string(bytes)) {
		fmt.Printf("  %s -> %s \n", text, target)
		links = append(links, Link{
			Source: trim(dir, pathPrefix, ".md"),
			Target: target,
			Text: text,
		})
	}
	return links
}

// recursively walk directory and return all files with given extension
func find(root, ext string) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			parse(s, root)
		}
		return nil
	})
}

func main() {
	find("../www/content", ".md")
}
