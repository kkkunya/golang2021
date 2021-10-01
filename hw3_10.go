package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeMarker struct {
	v   map[string]bool
	mux sync.Mutex
}

func (m *SafeMarker) Mark(key string) {
	m.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	m.v[key] = true
	m.mux.Unlock()
}

func (m *SafeMarker) Already(key string) bool {
	m.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer m.mux.Unlock()
	return m.v[key]
}

var m = SafeMarker{v: make(map[string]bool)}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, g *sync.WaitGroup) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.
	defer g.Done()

	if depth <= 0 {
		return
	}
	m.Mark(url)
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if m.Already(u){
			continue
		}
		g.Add(1)
		go Crawl(u, depth-1, fetcher, g)
	}
	return
}

func main() {
	// for goroutine with recursion,
	// the caller (in this context which is main) should
	// wait for all the goroutines
	g := &sync.WaitGroup{}
    g.Add(1)
	Crawl("https://golang.org/", 4, fetcher, g)
    g.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}