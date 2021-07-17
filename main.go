package main

import (
	"fmt"
	md "github.com/nikitavoloboev/markdown-parser"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type Link struct {
	Source string
	Target string
	Text   string
}

func parse(dir string) []Link {
	// read file
	bytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("in %s \n", dir)
	for text, target := range md.GetAllLinks(string(bytes)) {
		fmt.Printf("found link: %s -> %s \n", text, target)
		links = append(links, Link{
			Source: dir,
			Target: target,
			Text: text,
		})
	}
	return links
}

func find(root, ext string) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			parse(s)
		}
		return nil
	})
}

func main() {
	find("../www/content", ".md")
}
