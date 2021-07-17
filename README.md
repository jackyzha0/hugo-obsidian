# Obsidian Link Scrapper
GitHub action and binary to scrape [Obsidian](http://obsidian.md/) vault for links and exposes them as a `.yml` file for easy consumption by [Hugo](https://gohugo.io/).

## Installation
`go install github.com/jackyzha0/hugo-obsidian`

### Example Usage
Read Markdown from the `/content` folder and place the resulting `linkIndex.yaml` into `/data`

```shell
hugo-obsidian -input=content -output=data
```