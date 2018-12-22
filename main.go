// A quick script for converting Medium HTML files to Markdown, suitable for use in a static file generator such as Hugo or Jekyll
//A fork of https://gist.github.com/clipperhouse/010d4666892807afee16ba7711b41401
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lunny/html2md"
)

type post struct {
	Title, Author, Body   string
	Date, Lastmod         string
	Subtitle, Description string
	Canonical             string
	Draft                 bool
	IsComment             bool
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("usage: path/to/medium-export-folder/posts/ path/to/hugo/content")
		os.Exit(1)
	}
	// Location of exported, unzipped Medium HTML files
	var src = os.Args[1] // "/Users/mwsherman/medium-export"

	// Destination for Markdown files, perhaps the content folder for Hugo or Jekyll
	var dest = os.Args[2] //"/Users/mwsherman/tmp"

	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d articles.\n", len(files))

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".html") || f.IsDir() {
			fmt.Printf("Ignoring (ext) %s\n", f.Name())
			continue
		}

		inpath := filepath.Join(src, f.Name())
		doc, err := read(inpath)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		post, err := process(doc)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		post.Draft = strings.HasPrefix(f.Name(), "draft_")

		if post.IsComment {
			fmt.Printf("Ignoring (comment) %s\n", f.Name())
			continue
		}

		if len(post.Title) == 0 || len(post.Body) == 0 {
			fmt.Printf("Ignoring (empty) %s\n", f.Name())
			continue
		}
		prefix := "draft_"
		if post.Draft == false {
			//datetime ISO 2018-09-25T14:13:46.823Z
			//we only keep the date for simplicity
			pieces := strings.Split(post.Date, "T")
			prefix = pieces[0]
		}

		outpath := fmt.Sprintf("%s/%s_%s.md", dest, prefix, slug(post.Title))

		fmt.Printf("Processing %s => %s\n", f.Name(), outpath)
		write(post, outpath)
	}
}

func nbsp(r rune) rune {
	if r == '\u00A0' {
		return ' '
	}
	return r
}

func process(doc *goquery.Document) (p post, err error) {
	defer func() {
		if mypanic := recover(); mypanic != nil {
			err = mypanic.(error)
		}
	}()

	p = post{}
	p.Lastmod = time.Now().Format(time.RFC3339)
	p.Title = doc.Find("title").Text()
	p.Date, _ = doc.Find("time").Attr("datetime")
	p.Author = doc.Find(".p-author.h-card").Text()

	tmp := doc.Find(".p-summary[data-field='subtitle']")
	if tmp != nil {
		p.Subtitle = tmp.Text()
	}
	tmp = doc.Find(".p-summary[data-field='description']")
	if tmp != nil {
		p.Description = tmp.Text()
	}

	//if there are no subtitle and description we presume that it is a comment
	//to another story. Medium treats comments/replies as posts
	//I presume you do not want the comments as posts
	p.IsComment = len(p.Subtitle) == 0 && len(p.Description) == 0

	tmp = doc.Find(".p-canonical")
	if tmp != nil {
		//https://coder.today/a-b-tests-developers-manual-f57f5c1a492
		canonicalUrl, exists := tmp.Attr("href")
		if exists && len(canonicalUrl) > 0 {
			pieces := strings.Split(canonicalUrl, "/")
			if len(pieces) > 2 {
				//a-b-tests-developers-manual-f57f5c1a492
				canonicalUrl = pieces[len(pieces)-1] //we only need the last part
			}
		}
		p.Canonical = canonicalUrl
	}

	body := ""
	doc.Find("div.section-inner").Each(func(i int, s *goquery.Selection) {
		h, _ := s.Html()
		body += html2md.Convert(h)
	})
	body = strings.Map(nbsp, body)

	redundant := fmt.Sprintf("### %s", p.Title) // post body shouldn't repeat the title

	if strings.HasPrefix(body, redundant) {
		body = body[len(redundant):]
	}
	p.Body = strings.TrimSpace(body)

	return
}

func read(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Load the HTML document
	return goquery.NewDocumentFromReader(f)
}

func write(post post, path string) {
	os.Remove(path)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, post)
	if err != nil {
		panic(err)
	}
}

var spaces = regexp.MustCompile(`[\s]+`)
var notallowed = regexp.MustCompile(`[^\p{L}\p{N}.\s]`)
var athe = regexp.MustCompile(`^(a\-|the\-)`)

func slug(s string) string {
	result := s
	result = strings.Replace(result, "%", " percent", -1)
	result = strings.Replace(result, "#", " sharp", -1)
	result = notallowed.ReplaceAllString(result, "")
	result = spaces.ReplaceAllString(result, "-")
	result = strings.ToLower(result)
	result = athe.ReplaceAllString(result, "")

	return result
}

var tmpl = template.Must(template.New("").Parse(`---
title: "{{ .Title }}"
author: "{{ .Author }}"
date: {{ .Date }}
lastmod: {{ .Lastmod }}

{{ if eq .Draft true }}draft: {{ .Draft }}{{end}}
description: "{{ .Description }}"
slug: "{{ .Canonical }}"
aliases: [
    - "{{ .Canonical }}""
]
---

{{ .Body }}
`))
