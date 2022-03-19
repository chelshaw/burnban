package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	ttlcache "github.com/ReneKroon/ttlcache/v2"
	"github.com/gin-gonic/gin"
)

const CACHE_DURATION = time.Duration(10 * time.Minute)
type CountyData struct {
	Name  string  `json:"name"`
	Source  string  `json:"source"`
	Fetcher func(string)(string, error)
}
type Counties map[string]CountyData

func supportedCounties() (db Counties) {
	db = make(map[string]CountyData,4) // Bump number when we add new ones
	db["comal"] = CountyData{
		Name: "Comal", 
		Source: "https://www.co.comal.tx.us/Fire_Marshal.htm", 
		Fetcher: Comal,
	}
	db["hays"] = CountyData{
		Name: "Hays", 
		Source: "https://hayscountytx.com/law-enforcement/fire-marshal/", 
		Fetcher: Hays,
	}
	db["travis"] = CountyData{
		Name: "Travis", 
		Source: "https://www.traviscountytx.gov/fire-marshal/burn-ban", 
		Fetcher: Travis,
	}
	db["presidio"] = CountyData{
		Name: "Presidio", 
		Source: "https://www.co.presidio.tx.us/", 
		Fetcher: Presidio,
	}
	return db
}

func setupRouter(db Counties) *gin.Engine {
	cache := ttlcache.NewCache()
	cache.SetTTL(CACHE_DURATION)
	
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.LoadHTMLGlob("templates/*")
	// r.LoadHTMLFiles("assets/")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	// r.GET("/template", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "off.tmpl", gin.H{
	// 		"link": "https://google.com",
	// 		"county": "Here",
	// 	})
	// })

	r.GET("/", func(c *gin.Context) {
		// TODO: return list of all counties
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"counties": db,
		})
	})

	// r.Static("/assets", "assets")
	r.Static("/assets", "./assets")
	
	r.GET("/county/:county", func(c *gin.Context) {
		county := strings.ToLower(c.Params.ByName("county"))
		value, ok := db[county]

		// If county doesn't exist, return error
		if !ok {
			// TODO: Swap out for "request county" page
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"county": county,
				"error": "We don't support that county yet",
			})
			return;
		}

		if cachedBan, exists := cache.Get(county); exists == nil {
			log.Printf("Got cached value: %v\n", value)
			if (cachedBan != "") {
				var template = "on.tmpl"
				if cachedBan == "OFF" {
					template = "off.tmpl"
				}
				c.HTML(http.StatusOK, template, gin.H{
					"county": value.Name,
					"link": value.Source,
				})
				return;
			}
		}
		
		ban, err := value.Fetcher(value.Source)
		if err != nil {
			// TODO: Log potential bad url
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"error": err,
				"county": value.Name,
				"link": value.Source,
			})
			return;
		}

		if ban == "" {
			// TODO: Log, no information found from content
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"error": errors.New("could not find an answer"),
				"county": value.Name,
				"link": value.Source,
			})
			return;
		}

		var template = "on.tmpl"
		if ban == "OFF" {
			template = "off.tmpl"
		}
		c.HTML(http.StatusOK, template, gin.H{
			"county": value.Name,
			"link": value.Source,
		})
		cache.Set(county, ban)
	})
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
	})
	return r
}

func main() {
	log.Println("Locating ..")
	db := supportedCounties()
	log.Println("Preparing routes...")
	r := setupRouter(db)
	
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r.Run(":" + port);
}


/** COUNTIES */


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
