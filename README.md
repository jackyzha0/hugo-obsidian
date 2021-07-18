# Obsidian Link Scrapper
Used by [Quartz](https://github.com/jackyzha0/quartz)

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

Add 'Build Link Index' as a build step in your workflow file (e.g. `.github/workflows/deploy.yaml`)
```yaml
...

jobs:
  deploy:
    runs-on: ubuntu-18.04
    steps:
      - uses: actions/checkout@v2
      - name: Build Link Index
        uses: jackyzha0/hugo-obsidian@v2.1
        with:
          input: content # input folder
          output: data   # output folder
      ...
```

## Hugo Partial
To then embed this information in your Hugo site, you can copy and use the provided partials in `/partials`. Graph provides a graph view of all nodes and links and Backlinks provides a list of pages that link to this page. 

To start, create a `graphConfig.yaml` file in `/data` in your Hugo folder. This will be our main point of configuration for the graph partial.

Then, in one of your Hugo templates, do something like the following to render the graph.

```html
<div id="graph-container">
    {{partial "graph_partial.html" .}}
</div>
```

### Configuration
Example:

```yaml
enableLegend: false
enableDrag: true
enableZoom: false
base:
  node: "#284b63"
  activeNode: "#f28482"
  inactiveNode: "#a8b3bd"
  link: "#babdbf"
  activeLink: "#5a7282"
paths:
  - /toc: "#4388cc"
  - /newsletters: "#e0b152"
  - /posts: "#42c988"
```
