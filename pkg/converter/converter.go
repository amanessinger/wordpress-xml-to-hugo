package converter

import (
	"fmt"
	wp "github.com/amanessinger/wordpress-xml-go"
	"os"
	"text/template"
)

func Parse(path string) (error, *wp.WpXml) {
	var wp = wp.NewWordpressXml()
	var err = wp.ReadXml(path)
	if err != nil {
		return err, nil
	}
	return nil, &wp
}

func Convert(items []wp.Item) {
	for _, item := range items {
		if isApplicable(item) {
			convertPost(item)
		}
	}
}

func isApplicable(item wp.Item) bool {
	return item.PostType == "post"
}

func convertPost(item wp.Item) error {
	t := template.New("post_template")
	tp, err := t.Parse(
		`---
title: {{.Title}}
url: {{.Link}}
publishDate: {{.PubDate}}
date: {{.PostDate}}
categories:{{range .Categories}}{{if eq .Domain "category"}}
  - {{.UrlSlug}}{{end}}{{end}}
tags:{{range .Categories}}{{if eq .Domain "post_tag"}}
  - {{.UrlSlug}}{{end}}{{end}}
---
`)
	if err != nil {
		fmt.Printf("%q", err)
		return err
	}
	err = tp.Execute(os.Stdout, item)
	if err != nil {
		return err
	}
	return nil
}
