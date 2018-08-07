// opinionated conversion from WordPress to Hugo
package converter

import (
	// a fork of github.com/grokify/wordpress-xml-go with parsing of comments added
	wp "github.com/amanessinger/wordpress-xml-go"

	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	t := template.New("post_template")
	tp, err := t.Parse(PostTemplateSrc)
	if err != nil {
		panic(err)
	}

	postTargetPath := CreateSubPath(targetBaseDir, "post")

	for _, item := range items {
		if isPost(item) {
			if err = convertItem(item, tp, postTargetPath); err != nil {
				panic(err)
			}
		}
	}
}

// convert an item according to a template
func convertItem(item wp.Item, t *template.Template, itemBaseDir string) error {
	// make replacements
	item.Title = TitleReplacer.Replace(item.Title)
	item.Link = UrlReplacer2.Replace(item.Link)
	item.Content = UrlReplacer1.Replace(item.Content)
	item.Content = UrlReplacer2.Replace(item.Content)

	// construct and make the target directory
	targetPath := itemBaseDir
	itemPath := strings.Join(strings.Split(item.Link, "/"), string(filepath.Separator))
	itemSubDirPath := filepath.Dir(itemPath)
	CreateSubPath(targetPath, itemSubDirPath)

	// construct target file path
	itemFullPath := targetPath +
		string(filepath.Separator) +
		strings.TrimSuffix(itemPath, ".html") + // TODO: not everybody will have ".html"
		".md"

	// open target file
	os.Open(itemFullPath)
	f, err := os.OpenFile(itemFullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	err = t.Execute(f, item)
	if err != nil {
		return err
	}

	return nil
}
