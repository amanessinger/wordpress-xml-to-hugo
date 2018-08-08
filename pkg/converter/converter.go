// opinionated conversion from WordPress to Hugo
package converter

import (
	// a fork of github.com/grokify/wordpress-xml-go with parsing of comments added
	wp "github.com/amanessinger/wordpress-xml-go"

	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// parse a WordPress XML export
func Parse(path string) (error, *wp.WpXml) {
	var wpXml = wp.NewWordpressXml()
	var err = wpXml.ReadXml(path)
	if err != nil {
		return err, nil
	}
	return nil, &wpXml
}

// convert all items
func Convert(items []wp.Item, targetBaseDir string) {

	postBaseDir := CreateSubPath(targetBaseDir, "content/post")
	commentBaseDir := CreateSubPath(targetBaseDir, "comments/post")

	for _, item := range items {
		if isPost(item) {
			if err := convertItem(item, postBaseDir, commentBaseDir); err != nil {
				panic(err)
			}
		}
	}
}

// convert an item according to a template
func convertItem(item wp.Item, itemBaseDir string, commentBaseDir string) error {
	// make replacements
	item.Title = QuotesReplacer.Replace(item.Title)
	item.Link = UrlReplacer2.Replace(item.Link)
	item.Content = UrlReplacer1.Replace(item.Content)
	item.Content = UrlReplacer2.Replace(item.Content)
	item.Content = EmojiReplacer.Replace(item.Content)

	// construct and make the target directory
	targetPath := itemBaseDir
	// TODO: not everybody will have ".html"
	itemPath := strings.TrimSuffix(strings.Join(strings.Split(item.Link, "/"), string(filepath.Separator)), ".html")
	itemSubDirPath := filepath.Dir(itemPath)
	CreateSubPath(targetPath, itemSubDirPath)

	// construct target file path
	itemFullPath := targetPath +
		string(filepath.Separator) + itemPath + ".md"

	// open target file
	f, err := os.OpenFile(itemFullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	err = PostTemplate.Execute(f, item)
	if err != nil {
		return err
	}

	// if we have comments, create a comment directory named after the post
	if len(item.Comments) > 0 {
		commentDir := CreateSubPath(commentBaseDir, itemPath)
		if err = HandleComments(commentDir, item, convertComment); err != nil {
			return err
		}
	}

	return nil
}

// takes a func as handler to make it testable
func HandleComments(commentDir string, item wp.Item, handler func(wp.Comment, string, int) error) error {
	// capture replyTo relationships
	repliesTo := make(map[int]int)
	for _, c := range item.Comments {
		repliesTo[c.Id] = c.Parent
	}
	// determine names and write files
	for _, c := range item.Comments {
		commentFileName, indentLevel := GetCommentFileNameAndIndentLevel(repliesTo, c, commentDir)
		err := handler(c, commentFileName, indentLevel)
		if err != nil {
			return err
		}

	}
	return nil
}

// construct comment filename reflecting replyTo relationship, determine indent level
func GetCommentFileNameAndIndentLevel(repliesTo map[int]int, c wp.Comment, commentDir string) (string, int) {
	id := c.Id
	name := fmt.Sprintf("_%d", id)
	loop := true
	depth := -1
	for loop {
		parentId := repliesTo[id]
		name = fmt.Sprintf("_%d%s", parentId, name)
		if parentId == 0 {
			loop = false
		} else {
			id = parentId
		}
		depth++
	}
	commentFileName := commentDir +
		string(filepath.Separator) +
		fmt.Sprintf("comment%s.json", name)
	return commentFileName, depth
}

// write the comment
func convertComment(comment wp.Comment, commentFileName string, indentLevel int) error {
	// set indentation
	comment.IndentLevel = indentLevel

	// own comments may need replacements
	comment = FixCommentAuthor(comment)
	comment.AuthorUrl = UrlReplacer1.Replace(comment.AuthorUrl)
	comment.AuthorUrl = UrlReplacer2.Replace(comment.AuthorUrl)
	comment.Content = QuotesReplacer.Replace(comment.Content)
	comment.Content = EmojiReplacer.Replace(comment.Content)

	// open comment file
	f, err := os.OpenFile(commentFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	err = CommentTemplate.Execute(f, comment)
	if err != nil {
		return err
	}

	return nil
}
