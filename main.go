package main

import (
	"github.com/tdaira/gocrawl"
	"net/http"
	"time"
	"regexp"
	"github.com/PuerkitoBio/goquery"
	"crypto/sha1"
	"os"
	"bytes"
		"encoding/base64"
)

// Only enqueue the yahoo news domain.
var rxOk = regexp.MustCompile(`http://news\.yahoo\.co\.jp.*`)

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type ExampleExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
}

// Override Visit for our need.
func (x *ExampleExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(res.Body)
	body := bufBody.Bytes()
	url := ctx.NormalizedURL()

	sha := sha1.New()
	sha.Write([]byte(url.String()))
	shaStr := base64.URLEncoding.EncodeToString(sha.Sum(nil))
	file, err := os.Create("./data/" + shaStr)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(body)

	return nil, true
}

// Override Filter for our need.
func (x *ExampleExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func main() {
	// Set custom options
	opts := gocrawl.NewOptions(new(ExampleExtender))

	// should always set your robot name so that it looks for the most
	// specific rules possible in robots.txt.
	opts.RobotUserAgent = "Example"
	// and reflect that in the user-agent string used to make requests,
	// ideally with a link so site owners can contact you if there's an issue
	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogAll

	opts.MaxVisits = 20

	// Create crawler and start at root of duckduckgo
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://news.yahoo.co.jp/")
}
