package main

import (
	"flag"
	wikilink "github.com/abhinav/goldmark-wikilink"
	"github.com/yuin/goldmark"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(&wikilink.Extender{}),
	)
}

type Link struct {
	Source string
	Target string
	Text   string
}

type LinkTable = map[string][]Link
type Index struct {
	Links     LinkTable
	Backlinks LinkTable
}

type Content struct {
	Title   string
	Content string
}

type ContentIndex = map[string]Content

type IgnoredFiles struct {

}

func getIgnoredFiles() {

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
