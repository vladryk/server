package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"encoding/json"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page not found\n"))
}

func getTitle(url string, ch chan<- map[string]string) {
	urlTitle := make(map[string]string)
	var (
		resp *http.Response
		err  error
	)

	resp, err = http.Get(url)
	if err != nil {
		fmt.Printf("Page not found %s", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Can't read body: %s", err)
	}
	reg := regexp.MustCompile(`<(?:(?i)title>)(?P<title>.*)(?:</(?i)title>)`)
	match := reg.FindStringSubmatch(string(body))

	if len(match) == 0 {
		fmt.Printf("Tried to parse %s, got empty match\n\n", string(body))
	}

	urlTitle["url"] = url
	urlTitle["title"] = match[1]
	ch <- urlTitle
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	reg := regexp.MustCompile(`(https?://[^\s]+)`)
	matches := reg.FindAllString(string(body), -1)
	mainDict := make(map[string][]map[string]string)

	ch := make(chan map[string]string, len(matches))

	for _, el := range matches {
		for _, suf := range []string{".", ",", "?", ";"} {
			if strings.HasSuffix(el, suf) {
				el = el[:len(el)-1]
			}
		}
		go getTitle(el, ch)
	}

	for range matches {
		data := <-ch
		mainDict["links"] = append(mainDict["links"], data)
	}

	jsonString, err := json.Marshal(mainDict)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonString)
}
