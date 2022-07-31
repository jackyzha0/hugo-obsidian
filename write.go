package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func write(links []Link, contentIndex ContentIndex, toIndex bool, out string, root string) error {
	index := index(links)
	resStruct := struct {
		Index Index  `json:"index"`
		Links []Link `json:"links"`
	}{
		Index: index,
		Links: links,
	}
	marshalledIndex, mErr := json.MarshalIndent(&resStruct, "", "  ")
	if mErr != nil {
		return mErr
	}

	writeErr := ioutil.WriteFile(path.Join(out, "linkIndex.json"), marshalledIndex, 0644)
	if writeErr != nil {
		return writeErr
	}

	// check whether to index content
	if toIndex {
		marshalledContentIndex, mcErr := json.MarshalIndent(&contentIndex, "", "  ")
		if mcErr != nil {
			return mcErr
		}

		writeErr = ioutil.WriteFile(path.Join(out, "contentIndex.json"), marshalledContentIndex, 0644)
		if writeErr != nil {
			return writeErr
		}

		// write linkmap
		writeErr = writeLinkMap(&contentIndex, root)
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}

func writeLinkMap(contentIndex *ContentIndex, root string) error {
	fp := path.Join(root, "static", "linkmap")
	file, err := os.OpenFile(fp, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)
	for path := range *contentIndex {
		if path == "/" {
			_, _ = datawriter.WriteString("/index.html /\n")
		} else {
			_, _ = datawriter.WriteString(path + "/index.{html} " + path + "/\n")
		}
	}
	datawriter.Flush()
	file.Close()

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
