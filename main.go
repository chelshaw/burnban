package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ttlcache "github.com/ReneKroon/ttlcache/v2"
	"github.com/chelshaw/burnban/counties"
	"github.com/gin-gonic/gin"
)

const CACHE_DURATION = time.Duration(10 * time.Hour)
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
		Fetcher: counties.Comal,
	}
	db["hays"] = CountyData{
		Name: "Hays", 
		Source: "https://hayscountytx.com/law-enforcement/fire-marshal/", 
		Fetcher: counties.Hays,
	}
	db["travis"] = CountyData{
		Name: "Travis", 
		Source: "https://www.traviscountytx.gov/fire-marshal/burn-ban", 
		Fetcher: counties.Travis,
	}
	db["presidio"] = CountyData{
		Name: "Presidio", 
		Source: "https://www.co.presidio.tx.us/", 
		Fetcher: counties.Presidio,
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
	log.Println("Locating counties...")
	db := supportedCounties()
	log.Println("Preparing routes...")
	r := setupRouter(db)
	
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r.Run(":" + port);
}
