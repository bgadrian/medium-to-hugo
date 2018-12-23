# Medium to Hugo 
Tool for the migration of articles from a medium.com account to a static website generator with markdown like Hugo.

## Features:
* transform HTML posts to markdown
* ignores comments and empty articles
* keeps all the metadata and adds most Hugo front Matter, including the alias with the old URL
* even if one article fails it keeps going
* marks the drafts as "draft_"
* Fetch the article TAGS (which are not included in the Medium exporter)

## Usage 

1. Download your medium data
2. Unzip it
3. Download our binary (see Releases)
4. `./mediumtohugo /path/to/export/posts /path/to/hugo/content/`


### Build and run (with Go)
1. You need Bash, Go 1.11+

```bash
git clone git@github.com:bgadrian/medium-to-hugo.git
cd medium-to-hugo
env src=~/Documents/medium/posts dest=./output/ make run
```