package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"runtime"
	"time"
)

type result struct{
	url string
	urls []string
	err error
	depth int
}

var fetched map[string]bool

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.NumCPU())
	fetched = make(map[string]bool)
	now := time.Now()
	CrawlFast("http://jac-staging.jacarandafm.com", 2)
	fmt.Println("The time taken is", time.Since(now))
}

func Crawl(url string, depth int){
	if depth < 0{
		return
	}

	urls, err := findUrls(url)

	if err != nil{
		return
	}

	fmt.Printf("Found %s\n", url)
	fetched[url] = true

	for _, u := range urls{
		if !fetched[u]{
			Crawl(u, depth-1)
		}
	}
	return
}

func CrawlFast(url string, depth int){
	ch := make(chan *result)
	defer close(ch)

	fetch := func(url string, depth int) {
		urls, err := findUrls(url)
		ch <- &result{url, urls, err, depth}
	}

	go fetch(url, depth)
	fetched[url] = true

	for fetching:=1; fetching>0; fetching--{
		res := <-ch
		if res.err != nil{
			continue
		}
		fmt.Printf("found: %s\n", res.url)

		if res.depth > 0{
			for _, u := range res.urls{
				if !fetched[u]{
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}
}

func findUrls(url string)([]string, error){
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil{
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node)[] string{
	if n.Type == html.ElementNode && n.Data == "a"{
		for _, a := range n.Attr{
			if a.Key == "href"{
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling{
		links = visit(links, c)
	}
	return links
}