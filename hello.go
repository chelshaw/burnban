// from https://zetcode.com/golang/net-html/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func getHtmlPage(webPage string) (string, error) {

    resp, err := http.Get(webPage)

    if err != nil {
        return "", err
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {

        return "", err
    }

    return string(body), nil
}

func parseAndShow(text string) {

    tkn := html.NewTokenizer(strings.NewReader(text))

    var isTd bool
    var n int

    for {

        tt := tkn.Next()
				fmt.Println(tt)
        switch {

        case tt == html.ErrorToken:
            return

        case tt == html.StartTagToken:

            t := tkn.Token()
            isTd = t.Data == "td"

        case tt == html.TextToken:

            t := tkn.Token()

            if isTd {

                fmt.Printf("%s ", t.Data)
                n++
            }

            if isTd && n % 3 == 0 {

                fmt.Println()
            }

            isTd = false
        }
    }
}

func main() {

    webPage := "http://webcode.me/countries.html"
    data, err := getHtmlPage(webPage)

    if err != nil {
        log.Fatal(err)
    }
		// fmt.Printf(data)
    parseAndShow(data)
}
