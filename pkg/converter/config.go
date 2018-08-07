// CONFIGURATION
package converter

// text/template for output of Hugo markdown
const PostTemplateSrc = `---
title: "{{.Title}}"
url: {{.Link}}
publishDate: {{.PubDate}}
date: {{.PostDate}}
categories:{{range .Categories}}{{if eq .Domain "category"}}
  - "{{.UrlSlug}}"{{end}}{{end}}
tags:{{range .Categories}}{{if eq .Domain "post_tag"}}
  - "{{.UrlSlug}}"{{end}}{{end}}
---
{{.Content}}
`

// parsed post template
var PostTemplate = MakeParsedTemplate("post_template", PostTemplateSrc)

// URL replacements pass 1. Can't do it in one pass because of overlap
var UrlReplacements1 = []Replacement{
	{"http://www.manessinger.com", "http://manessinger.com"},
	{"href=\"manessinger", "href=\"http://manessinger"},
}

// URL replacements pass 2 in order of execution
var UrlReplacements2 = []Replacement{
	// img src URLs
	{"http://manessinger.com/images", "https://d25zfm9zpd7gm5.cloudfront.net"},
	// img target URLs
	{"http://manessinger.com/display.php/1024x1024", "https://d25zfm9zpd7gm5.cloudfront.net/1200x1200"},
	// URL of post, URLs of links to other posts: make it all server-absolute
	{"http://manessinger.com", ""},
}

// ready to use replacer
var UrlReplacer1 = MakeReplacer(UrlReplacements1...)
var UrlReplacer2 = MakeReplacer(UrlReplacements2...)

// Replacements in Title (because we quote it in FrontMatter)
var TitleReplacements = []Replacement{
	{"\"", "\\\""},
}

// ready to use replacer
var TitleReplacer = MakeReplacer(TitleReplacements...)
