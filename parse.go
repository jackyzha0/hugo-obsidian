package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// parse single file for links
func parse(dir, pathPrefix string) []Link {
	// read file
	source, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("[Parsing note] %s\n", trim(dir, pathPrefix, ".md"))

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

		target = strings.Replace(target, "%20", " ", -1)
		target = strings.Split(processTarget(target), "#")[0]
		target = strings.TrimSpace(target)
		target = strings.Replace(target, " ", "-", -1)

		source := filepath.ToSlash(hugoPathTrim(trim(dir, pathPrefix, ".md")))

		fmt.Printf("  '%s' => %s\n", text, target)
		links = append(links, Link{
			Source: UnicodeSanitize(source),
			Target: UnicodeSanitize(target),
			Text:   text,
		})
		n++
	})
	fmt.Printf("  Found: %d links\n", n)

	return links
}
