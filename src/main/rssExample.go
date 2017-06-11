package main

import (
	"sync"
	"github.com/ajstarks/svgo"
	rss "github.com/jteeuwen/go-pkg-rss"
	"time"
	"fmt"
	"os"
	"net/http"
	"strconv"
	"runtime"
	"log"
)

type Feed struct {
	url           string
	status        int
	itemCount     int
	complete      bool
	itemsComplete bool
	index         int
}

type FeedItem struct {
	feedIndex int
	complete  bool
	url       string
}

var feeds []Feed
var height int
var width int
var colors []string
var startTime int64
var timeout int
var feedSpace int

var wg sync.WaitGroup

func grabFeed(feed *Feed, feedChan chan bool, osvg *svg.SVG) {
	startGrab := time.Now().Unix()
	startGrabSeconds := startGrab - startTime

	fmt.Println("Grabbing feed", feed.url, " at", startGrabSeconds, "second mark")

	if feed.status == 0 {
		fmt.Println("Feed not yet read")
		feed.status = 1

		startX := int(startGrabSeconds * 33)
		startY := feedSpace * (feed.index)
		fmt.Println(startY)
		wg.Add(1)

		rssFeed := rss.New(timeout, true, channelHandler, itemsHandler)

		if err:= rssFeed.Fetch(feed.url, nil); err != nil{
			fmt.Fprintf(os.Stderr, "[e] %s: %s", feed.url, err)
			return
		}else{
			endSec := time.Now().Unix()
			endX := int(endSec - startGrab)
			if endX == 0{
				endX =1
			}
			fmt.Println("Read feed in",endX, "seconds")
			osvg.Rect(startX, startY, endX, feedSpace, "fill:#000000;opacity:.4")
			wg.Wait()
			endGrab := time.Now().Unix()
			endGrabSeconds := endGrab - startTime
			feedEndx := int(endGrabSeconds * 33)
			osvg.Rect(feedEndx, startY, 1, feedSpace, "fill:#ff0000;opacity:.9")
			feedChan <- true
		}
	}else if feed.status == 1{
		fmt.Println("Feed already in progress")
	}
}

func channelHandler(feed *rss.Feed, newChannels []*rss.Channel){

}

func itemsHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item){
	fmt.Println("Found",len(newItems),"items in", feed.Url)

	for i:= range newItems{
		url := *newItems[i].Guid
		fmt.Println(url)
	}

	wg.Done()
}

func getRSS(rw http.ResponseWriter, req *http.Request) {
	startTime  = time.Now().Unix()
	rw.Header().Set("Content-Type", "image/svg+xml")
	outputSVG := svg.New(rw)
	outputSVG.Start(width, height)
	feedSpace = (height - 20) / len(feeds)

	for i:=0; i< 30000; i++{
		timeText := strconv.FormatInt(int64(i/10), 10)
		if i %1000 == 0{
			outputSVG.Text(i/30, 390, timeText, "text-anchor:middle;font-size:10px;file:#000000")
		}else if i%4 == 0{
			outputSVG.Circle(i, 390, 1, "fill:#cccccc;stroke:none")
		}

		if i % 10 == 0{
			outputSVG.Rect(i, 0, 1, 400, "fill: #dddddd")
		}

		if i % 50 == 0{
			outputSVG.Rect(i, 0, 1, 400, "fill:#cccccc")
		}
	}

	feedChan := make(chan bool, 2)
	for i:= range feeds{
		outputSVG.Rect(0, i * feedSpace, width, feedSpace, "fill:"+colors[i] + ";stroke:none;")
		feeds[i].status = 0
		go grabFeed(&feeds[i], feedChan, outputSVG)
		<- feedChan
	}

	outputSVG.End()
}

func main() {
	runtime.GOMAXPROCS(2)
	timeout = 1000

	width = 1000
	height = 400
	feeds = append(feeds, Feed{index: 0, url:"https://groups.google.com/forum/feed/golang-nuts/msgs/rss_v2_0.xml?num=50", status: 0, itemCount: 0, complete: false, itemsComplete: false})
	feeds = append(feeds, Feed{index: 1, url:"https://groups.google.com/forum/feed/golang-dev/msgs/rss_v2_0.xml?num=50", status: 0, itemCount: 0, complete: false, itemsComplete: false })

	colors = append(colors,"#ff9999")
	colors = append(colors,"#9999ff")

	http.Handle("/getrss", http.HandlerFunc(getRSS))
	if err:=http.ListenAndServe(":8099", nil);err!= nil{
		log.Fatal("ListenAnd Server: ",err)

	}

}