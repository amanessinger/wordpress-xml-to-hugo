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

func TestHandleComments(t *testing.T) {
	item := wp.Item{}
	item.Comments = []wp.Comment{}

	c_0_1 := wp.Comment{}
	c_0_1.Parent = 0
	c_0_1.Id = 1
	item.Comments = append(item.Comments, c_0_1)

	c_0_2 := wp.Comment{}
	c_0_2.Parent = 0
	c_0_2.Id = 2
	item.Comments = append(item.Comments, c_0_2)

	c_1_3 := wp.Comment{}
	c_1_3.Parent = 1
	c_1_3.Id = 3
	item.Comments = append(item.Comments, c_1_3)

	c_3_4 := wp.Comment{}
	c_3_4.Parent = 3
	c_3_4.Id = 4
	item.Comments = append(item.Comments, c_3_4)

	c_2_5 := wp.Comment{}
	c_2_5.Parent = 2
	c_2_5.Id = 5
	item.Comments = append(item.Comments, c_2_5)

	if err := HandleComments("comments/post/2018/09/1000-some-title", item,
		func(comment wp.Comment, commentFileName string, indentLevel int) error {
			fmt.Printf("%d: indent: %d, %s\n", comment.Id, indentLevel, commentFileName)
			return nil
		}); err != nil {
		t.Errorf("%v", err)
	}

}

func TestEliminateAmazonAds(t *testing.T) {
	in := `<div class="container">
<div class="center"><a target="_blank" href="http://manessinger.com/display.php/1024x1024/2013/20131006_180455_lr.jpg"><img src="http://manessinger.com/images/0600x0600/2013/20131006_180455_lr.jpg" /></a></div>
</div>
<br />

Here's another Underground image of one of Vienna's more interesting stations. It is at the crossing of two lines, one of them very deep under ground, because it also crosses below a canal.

<iframe src="http://rcm-na.amazon-adsystem.com/e/cm?lt1=_blank&bc1=000000&IS2=1&bg1=000000&fc1=FFFFFF&lc1=99AADD&t=thedailphotof-20&o=1&p=8&l=as4&m=amazon&f=ifr&ref=ss_til&asins=B001HDVGIY" style="margin: 0pt 0px 0pt 10px; float: right; width:120px;height:240px;" scrolling="no" marginwidth="0" marginheight="0" frameborder="0"></iframe> The Song of the Day is "Bottom Below" from the 2008 album "Dirt Don't Hurt" by Holly Golightly &amp; The Brokeoffs. Hear it on <a href="http://www.youtube.com/watch?v=Kl5kZOHoJIg" target="_blank">YouTube</a>.`
	expected := `<div class="container">
<div class="center"><a target="_blank" href="http://manessinger.com/display.php/1024x1024/2013/20131006_180455_lr.jpg"><img src="http://manessinger.com/images/0600x0600/2013/20131006_180455_lr.jpg" /></a></div>
</div>
<br />

Here's another Underground image of one of Vienna's more interesting stations. It is at the crossing of two lines, one of them very deep under ground, because it also crosses below a canal.

 The Song of the Day is "Bottom Below" from the 2008 album "Dirt Don't Hurt" by Holly Golightly &amp; The Brokeoffs. Hear it on <a href="http://www.youtube.com/watch?v=Kl5kZOHoJIg" target="_blank">YouTube</a>.`
	out := EliminateAmazonAds(in)
	if out != expected {
		t.Errorf("Expected [%s], got [%s]", expected, out)
	}
}
