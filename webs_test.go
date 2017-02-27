package main

import (
	"fmt"
	"testing"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"os"
)

const (
	requestNumber = 10
	url = "http://127.0.0.1:8000/analyze"
	testStr = " hello http://info.cern.ch/hypertext/WWW/TheProject.html"
	testOutput = `{"links":[{"title":"The World Wide Web project","url":"http://info.cern.ch/hypertext/WWW/TheProject.html"}]}`
)

func TestMain(m *testing.M) {
	// Run server
	go func() {
		serv := getServer()
		log.Fatal(serv.ListenAndServe())
	}()
	os.Exit(m.Run())
}

func TestTest(t *testing.T)  {
	ch := make(chan string, requestNumber)
	for i:=0; i<requestNumber; i++ {
		timer := time.NewTimer(200 * time.Millisecond)  // We need time-sleep, because often target server has no time to answer (or considering as a DDos attack)
		<-timer.C
		go MakeRequest(url, ch)
	}

	for i:=0; i<requestNumber; i++{
		responseData := <-ch
		bool1 := strings.EqualFold(testOutput, string(responseData))
		if bool1 == false {
			t.Error("Different strings:\n", responseData, "\nShould be:\n", testOutput)
		}
	}

}

func MakeRequest(url string, ch chan<-string) {
	resp, _ := http.Post(url, "application/text", bytes.NewBufferString(testStr))
	body, errs := ioutil.ReadAll(resp.Body)
	if errs != nil {
		log.Fatal(errs)
	}
	ch <- fmt.Sprintf(string(body))
}