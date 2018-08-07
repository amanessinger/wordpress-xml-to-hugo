// wordpress-xml-to-hugo parses an XML export from WordPress and generates Markdown files for Hugo

package main

import (
	"fmt"
	c "github.com/amanessinger/wordpress-xml-to-hugo/pkg/converter"
	"os"
)

// commandline processing only. Everything else is in pkg/converter
func main() {
	if len(os.Args) != 3 {
		Usage()
	}

	inFilePath := os.Args[1]
	if _, err := os.Stat(inFilePath); os.IsNotExist(err) {
		fmt.Printf("Wordpress XML export does not exist")
		Usage()
	}
	targetBaseDir := os.Args[2]
	dirInfo, err := os.Stat(targetBaseDir)
	if os.IsNotExist(err) {
		fmt.Printf("output base directory does not exist")
		Usage()
	}
	if !dirInfo.IsDir() {
		fmt.Printf("output base is no directory")
		Usage()
	}

	err, parsed := c.Parse(inFilePath)
	if err != nil {
		fmt.Printf("Parse error: %q", err)
		Usage()
	}

	c.Convert(parsed.Channel.Items, targetBaseDir)

	fmt.Printf("parsed and converted a file with %d items\n", len(parsed.Channel.Items))
	os.Exit(0)
}

func Usage() {
	fmt.Printf("Usage: %s <path-to-wp-export> <target-base-dir>", os.Args[0])
	os.Exit(1)
}
