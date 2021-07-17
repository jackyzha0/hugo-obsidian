# Obsidian Link Scrapper
This repository comes to you in two parts.

1. GitHub Action (scrapes links into a `.yml` file)
2. Hugo Partial (turns `.yml` file into graphs and tables)

## GitHub Action
GitHub action and binary to scrape [Obsidian](http://obsidian.md/) vault for links and exposes them as a `.yml` file for easy consumption by [Hugo](https://gohugo.io/).
### Example Usage (Binary)
Read Markdown from the `/content` folder and place the resulting `linkIndex.yaml` into `/data`

```shell
# Installation
go install github.com/jackyzha0/hugo-obsidian

# Run
hugo-obsidian -input=content -output=data
```

### Example Usage (GitHub Action)



## Hugo Partial


### Configuration
```yaml
enableLegend: false
enableDrag: true
enableZoom: false
base:
  node: "#284b63"
  activeNode: "#f28482"
  inactiveNode: "#a8b3bd"
  hoverNode: "#afbfc9"
  link: "#aeb4b8"
  activeLink: "#5a7282"
paths:
  - /toc: "#4388cc"
  - /newsletters: "#e0b152"
  - /posts: "#42c988"
```