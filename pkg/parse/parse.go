package parse

import (
	p "github.com/amanessinger/wordpress-xml-go"
)

func Parse(path string) (error, *p.WpXml) {
	var wp = p.NewWordpressXml()
	var err = wp.ReadXml(path)
	if err != nil {
		return err, nil
	}
	return nil, &wp
}
