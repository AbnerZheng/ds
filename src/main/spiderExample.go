package main

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"sync"
)

var urls []string
var urlsProcessed int
var fullText string
var totalURLCount int
var wg sync.WaitGroup


func readURLs(statusChannel chan int, textChannel chan string) {

	time.Sleep(time.Millisecond * 1)
	fmt.Println("Grabbing", len(urls), "urls")
	for i := 0; i < totalURLCount; i++ {

		fmt.Println("Url", i, urls[i])
		resp, _ := http.Get(urls[i])
		text, err := ioutil.ReadAll(resp.Body)

		textChannel <- string(text)

		if err != nil {
			fmt.Println("No HTML body")
		}

		statusChannel <- 0
	}
	fmt.Println("exit read")
	wg.Done()
}

func addToScrapedText(textChannel chan string, processChannel chan bool) {
	for {
		select {
		case pC := <-processChannel:
			if pC == true {
				//hang on
			}
			if pC == false {
				close(textChannel)
				close(processChannel)
				wg.Done()
				fmt.Println("exit add")
				return
			}
		case tC := <-textChannel:
			fullText += tC
		}
	}
}

func evaluateStatus(statusChannel chan int, applicationChannel chan bool, processChannel chan bool) {
	for {
		select {
		case status := <-statusChannel:
			fmt.Println(urlsProcessed, totalURLCount)
			urlsProcessed++
			if status == 0 {
				fmt.Println("Got url")
			}
			if status == 1 {
				close(statusChannel)
			}
			if urlsProcessed == totalURLCount {
				fmt.Println("Read all top-level URLS")
				processChannel <- false
				applicationChannel <- false
				wg.Done()
				fmt.Println("exit eval")
				return
			}
		}
	}
}

func main() {
	applicationChannel := make(chan bool)
	statusChannel := make(chan int)
	textChannel := make(chan string)
	processChannel := make(chan bool)
	totalURLCount = 0

	urls = append(urls, "http://www.mastergoco.com/index1.html")
	urls = append(urls, "http://www.mastergoco.com/index2.html")
	urls = append(urls, "http://www.mastergoco.com/index3.html")
	urls = append(urls, "http://www.mastergoco.com/index4.html")
	urls = append(urls, "http://www.mastergoco.com/index5.html")
	fmt.Println("Start spider")

	urlsProcessed = 0
	totalURLCount = len(urls)

	go evaluateStatus(statusChannel, applicationChannel, processChannel)
	go readURLs(statusChannel, textChannel)
	go addToScrapedText(textChannel, processChannel)
	wg.Add(3)

	select {
	case sC := <-applicationChannel:
		if sC == false {
			fmt.Println(fullText)
			fmt.Println("Done!")
			break
		}
	}
	wg.Wait()
	fmt.Println("exit main")
}
