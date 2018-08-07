package converter

import (
	"fmt"
	wp "github.com/amanessinger/wordpress-xml-go"
	"os"
	"strings"
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

// PARSE TESTS - just to make sure parsing did no break due to changes in the imported project
func TestParsedItemNumberCorrect(t *testing.T) {
	want := 9
	got := len(parsed.Channel.Items)
	if got != want {
		t.Errorf("Expected %d items, got %d", want, got)
	}
	fmt.Printf("Got %d items\n", got)
}

func TestParsedFirstItemIsNoPost(t *testing.T) {
	if isPost(parsed.Channel.Items[0]) {
		t.Errorf("Expected first It to be a page, not a post")
	}
}

// REPLACEMENT TESTS
func TestMakeReplacer(t *testing.T) {
	input := "__aaa__aa__XXaaaXX__"
	expect := "__222__33__yy111yy__"
	r := MakeReplacer([]Replacement{
		{"XXaaaXX", "yy111yy"},
		{"aaa", "222"},
		{"aa", "33"},
	}...)
	result := r.Replace(input)
	if result != expect {
		t.Errorf("Expected %s, got %s", expect, result)
	}
}

func TestUrlReplacements(t *testing.T) {
	item := parsed.Channel.Items[8]
	result := UrlReplacer1.Replace(item.Content)
	result = UrlReplacer2.Replace(result)
	if strings.Index(result, "manessinger.com") != -1 {
		t.Errorf("Result shouldn't contain manessinger.com\n\n%s", result)
	}
}
