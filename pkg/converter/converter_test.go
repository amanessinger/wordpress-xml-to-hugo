package converter

import (
	"fmt"
	wp "github.com/amanessinger/wordpress-xml-go"
	"os"
	"testing"
)

// one global parsed blog export for all tests to operate on
var (
	err    error
	parsed *wp.WpXml
)

// set up the global parsed export and run the tests
func TestMain(m *testing.M) {
	err, parsed = Parse("../../testdata/short.xml")
	if err != nil {
		panic(fmt.Sprintf("Parse error: %q", err))
	}
	os.Exit(m.Run())
}

func TestItemNumberCorrect(t *testing.T) {
	want := 8
	got := len(parsed.Channel.Items)
	if got != want {
		t.Errorf("Expected %d items, got %d", want, got)
	}
	fmt.Sprintf("Got %d items", got)
}

func TestFirstItemIsNoPost(t *testing.T) {
	if isApplicable(parsed.Channel.Items[0]) {
		t.Errorf("Expected first Item to be a page, not a post")
	}
}
