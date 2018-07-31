package parse

import (
	"testing"
)

func TestParseXml(t *testing.T) {

	var err, parsed = Parse("../../testdata/short.xml")

	if err != nil {
		t.Errorf("Parse error: %q", err)
	}

	want := 7
	got := len(parsed.Channel.Items)
	if got != want {
		t.Errorf("Expected %d items, got %d", want, got)
	}
}
