// from https://zetcode.com/golang/net-html/
package burnban

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func scrape(url string) *goquery.Document {
	res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }
	return doc
}

func ExampleScrape() string {
  // Request the HTML page.
  res, err := http.Get("https://www.co.comal.tx.us/Fire_Marshal.htm")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  // mystring := doc.Find("#menu-v li").First().Text()
	// fmt.Println(mystring)
	var stringFound string
	doc.Find("ul#menu-v li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Text()
		// fmt.Printf("Reading line %d", i)
		if (strings.Contains(strings.ToLower(title), "burn ban is")) {
			// fmt.Printf("Found %d: %s\n", i, title)
			stringFound = title
		}
	})
	return stringFound
}

func Comal() (found bool, ban bool) {
	doc := scrape("https://www.co.comal.tx.us/Fire_Marshal.htm")
	var stringFound string
	doc.Find("ul#menu-v li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		content := s.Text()
		if (strings.Contains(strings.ToLower(content), "burn ban is")) {
			stringFound = content
		}
	})
	return stringFound != "", strings.Contains(strings.ToLower(stringFound), "is on")
}

func Travis() (found bool, ban bool) {
	doc := scrape("https://www.traviscountytx.gov/fire-marshal/burn-ban")
	var stringFound string
	doc.Find("#burnban div").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		content := s.Text()
		if (strings.Contains(strings.ToLower(content), "burn ban is")) {
			stringFound = content
		}
	})
	return stringFound != "", strings.Contains(strings.ToLower(stringFound), "is in effect")
}

func Hays() (found bool, ban bool) {
	doc := scrape("https://hayscountytx.com/law-enforcement/fire-marshal/")
	var stringFound string
	var finished = false
	doc.Find("#et-boc p").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		content := s.Text()
		if stringFound != "" && !finished {
			stringFound = stringFound + content
			finished = true
		} else if (!finished && strings.Contains(strings.ToLower(content), "burn ban is")) {
			stringFound = content
		}
	})
	return stringFound != "", strings.Contains(strings.ToLower(stringFound), "ON")
}

// func FindCounty(name string) (bool, error) {
// 	if name == "" {
// 		return false, errors.New("No name provided")
// 	}
// 	var lowerName = strings.ToLower(name)
// 	var result bool
// 	switch lowerName {
// 		case "comal": 
// 			result = Comal()
// 		case "travis": 
// 			result = Travis()
// 		default: 
// 			fmt.Println("Default this")
// 	}
// 	// parse the result string
// 	return result, nil
// }

// func BanOnImage() int64err {
// 	// lekkewords := []string{
// 	// 	"https://images.unsplash.com/photo-1511027643875-5cbb0439c8f1",
// 	// }
// 	return rand.Seed(time.Now().UnixNano())
// }

func main() {
	fmt.Println("Getting data...")
  found := ExampleScrape()
	if (len(found) > 0) {
		fmt.Println("Result was found")
		fmt.Println(found)
	} else {
		fmt.Println("Sorry, we couldn't find an answer")
	}
}
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

// func main() {

//     webPage := "http://webcode.me/countries.html"
//     data, err := getHtmlPage(webPage)

//     if err != nil {
//         log.Fatal(err)
//     }
// 		// fmt.Printf(data)
//     parseAndShow(data)
// }
