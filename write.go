package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

const message = "# THIS FILE WAS GENERATED USING github.com/jackyzha0/hugo-obsidian\n# DO NOT EDIT\n"
func write(links []Link, contentIndex ContentIndex, toIndex bool, out string) error {
	index := index(links)
	resStruct := struct {
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



