package main

import (
	"fmt"
	"github.com/adjust/redismq"
	"github.com/tdaira/gocrawl"
	"github.com/tdaira/test/app/crawler"
	"regexp"
	"time"
)

func main() {
	workQueue := redismq.CreateQueue(
		"localhost", "6379", "", 0, "crawl_site")
	consumer, err := workQueue.AddConsumer("crawler_service_consumer")
	if err != nil {
		panic(err)
	}
	err = consumer.ResetWorking()
	if err != nil {
		panic(err)
	}
	for {
		pkg, err := consumer.Get()
		if err != nil {
			panic(err)
		}
		siteInfo, err := crawler.StringToSiteInfo(pkg.Payload)
		if err != nil {
			fmt.Println("Unmarshal error: " + err.Error())
		}
		pkg.Ack()
		go func() {
			crawl(siteInfo)
		}()
	}
}

func crawl(siteInfo *crawler.SiteInfo) {
	opts := gocrawl.NewOptions(&crawler.MirageBotExtender{
		ValidURLRegex: regexp.MustCompile(siteInfo.ValidUrlRegex)})

	opts.RobotUserAgent = "MirageBot"
	opts.UserAgent = "Mozilla/5.0 (compatible; MirageBot/1.0; +http://miragebot.com)"

	opts.LogFlags = gocrawl.LogAll

	opts.CrawlDelay = time.Duration(siteInfo.CrawlDelay) * time.Second
	opts.SameHostOnly = siteInfo.SameHostOnly
	opts.MaxVisits = siteInfo.MaxVisits

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(siteInfo.Seeds)
}
