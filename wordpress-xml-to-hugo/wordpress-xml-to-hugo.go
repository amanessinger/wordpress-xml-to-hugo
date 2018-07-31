// wordpress-xml-to-hugo parses an XML export from WordPress and generates Markdown files for Hugo
//
//

package main

import (
	"fmt"
	p "github.com/amanessinger/wordpress-xml-to-hugo/pkg/parse"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <path-to-wp-export>", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	err, parsed := p.Parse(path)
	if err != nil {
		fmt.Printf("Parse error: %q", err)
		os.Exit(1)
	}

	fmt.Printf("parsed a file with %d items", len(parsed.Channel.Items))
	os.Exit(0)
}
