// CONFIGURATION
package converter

import (
	wp "github.com/amanessinger/wordpress-xml-go"
	"regexp"
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
const CommentTemplateSrc = `{
    "id": "{{ .Id }}",
    "author": "{{ .Author }}",
    "author_url": "{{ .AuthorUrl }}",
    "date": "{{ .DateGmt }}",
    "indent_level": {{ .IndentLevel }},
    "content": "{{ .Content }}"
}
`

// parsed comment template
var CommentTemplate = MakeParsedTemplate("comment_template", CommentTemplateSrc)

// URL replacements pass 1. Can't do it in one pass because of overlap
var urlReplacements1 = []Replacement{
	{"http://www.manessinger.com", "http://manessinger.com"},
	{"href=\"manessinger", "href=\"http://manessinger"},
}
var UrlReplacer1 = MakeReplacer(urlReplacements1...)

// URL replacements pass 2
var urlReplacements2 = []Replacement{
	// img src URLs
	{"http://manessinger.com/images", "https://d25zfm9zpd7gm5.cloudfront.net"},
	// img target URLs
	{"http://manessinger.com/display.php/1024x1024", "https://d25zfm9zpd7gm5.cloudfront.net/1200x1200"},
	// URL of post, URLs of links to other posts: make it all server-absolute
	{"http://manessinger.com", ""},
}

// ready to use replacer
var UrlReplacer2 = MakeReplacer(urlReplacements2...)

// Replacements in Title and content of comments (because we quote it in FrontMatter)
var quotesReplacements = []Replacement{
	{"\"", "\\\""},
	{"\n", "\\n"},
	{"\t", " "},
}

// ready to use replacer
var QuotesReplacer = MakeReplacer(quotesReplacements...)

// Replacements in content of posts and comments
var emojiReplacements = []Replacement{
	{":)", "🙂"},
	{":(", "☹️"},
	{":p", "😛"},
	{":P", "😛"},
	{":D", "😄"},
	{";)", "😉"},
	{":-)", "🙂"},
	{":-(", "☹️"},
	{":-\\", "😏"},
	{":roll:", "🙄"},
}

// ready to use replacer
var EmojiReplacer = MakeReplacer(emojiReplacements...)

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

// certainly not for everybody: eliminate Amazon ads embedded as iframes
var amazonAdRegexp = regexp.MustCompile("<iframe [^>]+></iframe>")

func EliminateAmazonAds(content string) string {
	return string(amazonAdRegexp.ReplaceAllLiteral([]byte(content), []byte("")))
}
