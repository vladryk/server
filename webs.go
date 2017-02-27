package main

import ("io"
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"regexp"
	"strings"
	"encoding/json"
	"github.com/gorilla/mux"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Page not found\n")
}

func getTitle(url string, ch chan<-map[string]string) {
	urlTitle := make(map[string]string)
	var (
		resp *http.Response
		err error
	)

	resp, err = http.Get(url)
	if err != nil {
		fmt.Printf("Page not found %s", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`<TITLE>(?s)(.*)</TITLE>`)

	matches := reg.FindString(string(body))



	matches = strings.Replace(matches, "<TITLE>", "", 1) // FIXME need to use re groups
	matches = strings.Replace(matches, "</TITLE>", "", 1)

	if len(matches) == 0 {
		fmt.Printf("Tried to parse %s, got empty match\n\n", string(body))
	}


	urlTitle["url"] = url
	urlTitle["title"] = matches
	ch <- urlTitle
}

func analyzeHandler(w http.ResponseWriter, r *http.Request)  {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
    		log.Fatal(err)
	}
	reg := regexp.MustCompile(`(https?://[^\s]+)`)
	matches := reg.FindAllString(string(body), -1)
	mainDict := make(map[string][]map[string]string)

	ch := make(chan map[string]string, len(matches))

	for _, el := range matches {
		if strings.HasSuffix(el, "."){  // FIXME: We need use list of symbols
			el = el[:len(el)-1]
		}

		go getTitle(el, ch)
	}

	for range matches{
		data := <- ch
		mainDict["links"] = append(mainDict["links"], data)
	}

	jsonString, err := json.Marshal(mainDict)
	 w.Write(jsonString)
}

func getServer() *http.Server{
	r := mux.NewRouter()
	r.HandleFunc("/", notFoundHandler)
	r.HandleFunc("/analyze", analyzeHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
    	}
	return srv
}

func main()  {
	serv := getServer()
	log.Fatal(serv.ListenAndServe())
}
