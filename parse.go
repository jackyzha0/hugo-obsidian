package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yuin/goldmark"
)

// parse single file for links
func parse(md goldmark.Markdown, fileIndex fileIndex, dir, pathPrefix string) []Link {
	// read file
	source, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("[Parsing note] %s => ", trim(dir, pathPrefix, ".md"))

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	var n int
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		target, ok := s.Attr("href")
		if !ok {
			target = "#"
		}

		target = processTarget(target)
		source := processSource(trim(dir, pathPrefix, ".md"))

		target = fileIndex.resolve(target)

		// fmt.Printf("  '%s' => %s\n", source, target)
		if !strings.HasPrefix(text, "^") {
			links = append(links, Link{
				Source: source,
				Target: target,
				Text:   text,
			})
			n++
		}
	})
	fmt.Printf("found: %d links\n", n)

	return links
}
