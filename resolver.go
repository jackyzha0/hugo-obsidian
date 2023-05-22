package main

import (
	"fmt"
	"path/filepath"
	wikilink "github.com/abhinav/goldmark-wikilink"
)
var QuartzResolver Resolver = quartzResolver{}

type Resolver interface {
	ResolveWikilink(*wikilink.Node) (destination []byte, err error)
}

var _html = []byte(".html")
var _hash = []byte("#")

type quartzResolver struct{}

func (quartzResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	dest := make([]byte, len(n.Target)+len(_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], _html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		fmt.Println(n.Fragment)
		for f := 0; f < len(n.Fragment); f++ {
			fmt.Println(string(n.Fragment[f]))
			if string(n.Fragment[f]) != "." {
				dest[i] = n.Fragment[f]
				i++
			}
		}
		fmt.Println(dest[:i])
	}
	return dest[:i], nil
}
