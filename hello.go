// from https://zetcode.com/golang/net-html/
package burnban

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


func scrapeSite(url string, bodySelector string) (content string, err error){
	// TODO: Check that string starts with https://
	if bodySelector == "" || url == "" {
		return "", errors.New("URL and selector are both required")
	}
	res, getError := http.Get(url);
	if getError != nil || res.StatusCode != 200 {
		return "", errors.New("requested URL is unavailable")
	}
	
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	doc.Find(bodySelector).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		contentString := s.Text()
		if contentString != "" {
			content = content + " " + strings.TrimSpace(contentString)
		}
	})
	if content == "" {
		return "", errors.New("no content found at URL")
	}
	return;
}

func Travis(url string) (ban string, err error) {
	content, err := scrapeSite(url, "#burnban div");
	if err != nil {
		return "", err
	}
	// Check for existence of key phrases
	if strings.Contains(strings.ToLower(content), "burn ban is off") {
		ban = "OFF"
	} else if strings.Contains(strings.ToLower(content), "burn ban is in effect") || strings.Contains(strings.ToLower(content), "burn ban is currently:on")  {
		ban = "ON"
	}
	return
}

func Hays(url string) (ban string, err error) {
	content, err := scrapeSite(url, ".entry-content p");
	if err != nil {
		return "", err
	}
	
	// Check for existence of key phrases
	if strings.Contains(strings.ToLower(content), "burn ban is currently:off") {
		ban = "OFF"
	} else if strings.Contains(strings.ToLower(content), "burn ban in effect") || strings.Contains(strings.ToLower(content), "burn ban is currently:on")  {
		ban = "ON"
	}
	return
}

func Comal(url string) (ban string, err error) {
	content, err := scrapeSite(url, "ul#menu-v li");
	if err != nil {
		return "", err
	}
	
	// Check for existence of key phrases
	if strings.Contains(strings.ToLower(content), "burn ban is off") {
		ban = "OFF"
	} else if strings.Contains(strings.ToLower(content), "burn ban is on") || strings.Contains(strings.ToLower(content), "burn ban is currently:on")  {
		ban = "ON"
	}
	return
}
func Presidio(url string) (ban string, err error) {
	content, err := scrapeSite(url, "#ContentPlaceHolder4_ContentRepeater4_WidgetBox_3 span");
	if err != nil {
		return "", err
	}
	
	// Check for existence of key phrases
	if strings.Contains(strings.ToLower(content), "burn ban is off") {
		ban = "OFF"
	} else if strings.Contains(strings.ToLower(content), "burn ban in effect") || strings.Contains(strings.ToLower(content), "burn ban is currently:on")  {
		ban = "ON"
	}
	return
}
