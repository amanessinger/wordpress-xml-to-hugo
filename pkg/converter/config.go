// CONFIGURATION
package converter

import (
	wp "github.com/amanessinger/wordpress-xml-go"
	"strings"
)

// text/template for posts
const PostTemplateSrc = `---
title: "{{ .Title }}"
url: {{ .Link }}
publishDate: {{ .PubDate }}
date: {{ .PostDate }}
categories: {{ range .Categories }}{{ if eq .Domain "category" }}
  - "{{ .UrlSlug }}"{{ end }}{{ end }}
tags: {{ range .Categories }}{{ if eq .Domain "post_tag" }}
  - "{{ .UrlSlug }}"{{ end }}{{ end }}
---
{{ .Content }}
`

// parsed post template
var PostTemplate = MakeParsedTemplate("post_template", PostTemplateSrc)

// text/template for comments
const CommentTemplateSrc = `---
author: "{{ .Author }}"
author_email: {{ .AuthorEmail }}
author_url: {{ .AuthorUrl }}
date: {{ .DateGmt }}
indent_level: {{ .IndentLevel }}
---
{{ .Content }}
`

// parsed comment template
var CommentTemplate = MakeParsedTemplate("comment_template", CommentTemplateSrc)

// URL replacements pass 1. Can't do it in one pass because of overlap
var UrlReplacements1 = []Replacement{
	{"http://www.manessinger.com", "http://manessinger.com"},
	{"href=\"manessinger", "href=\"http://manessinger"},
}

// ready to use replacer
var UrlReplacer1 = MakeReplacer(UrlReplacements1...)

// URL replacements pass 2
var UrlReplacements2 = []Replacement{
	// img src URLs
	{"http://manessinger.com/images", "https://d25zfm9zpd7gm5.cloudfront.net"},
	// img target URLs
	{"http://manessinger.com/display.php/1024x1024", "https://d25zfm9zpd7gm5.cloudfront.net/1200x1200"},
	// URL of post, URLs of links to other posts: make it all server-absolute
	{"http://manessinger.com", ""},
}

// ready to use replacer
var UrlReplacer2 = MakeReplacer(UrlReplacements2...)

// Replacements in Title (because we quote it in FrontMatter)
var TitleReplacements = []Replacement{
	{"\"", "\\\""},
}

// ready to use replacer
var TitleReplacer = MakeReplacer(TitleReplacements...)

// maybe not for everybody, but this author needs to be unified
func FixCommentAuthor(comment wp.Comment) wp.Comment {
	if comment.Author == "advman" ||
		(comment.Author == "andreas" && (strings.Index(comment.AuthorUrl, "manessinger.com") != -1)) {
		comment.Author = "andreas"
		comment.AuthorEmail = "info@andreas.manessinger.info"
		comment.AuthorUrl = "https://manessinger.com/"
	}
	return comment
}
