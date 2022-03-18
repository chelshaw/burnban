package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chelshaw/burnban"

	"github.com/gin-gonic/gin"
)

type county struct {
	county		string
	state			string
	link			string
	selector 	string
}
type CountyData struct {
	Name  string  `json:"name"`
	Source  string  `json:"source"`
	Fetcher func(string)(string, error)
	LastGet time.Time `json:"last_get"`
	CachedBan string `json:"cached_ban"`
}

func setupRouter() *gin.Engine {
	// Seed data
	seedTime := time.Date(2022,3,17,0,0,0,0,time.UTC)
	db := make(map[string]CountyData,10)
	db["comal"] = CountyData{
		Name: "Comal", 
		Source: "https://www.co.comal.tx.us/Fire_Marshal.htm", 
		Fetcher: burnban.Comal,
		LastGet: seedTime,
		CachedBan: "",
	}
	db["hays"] = CountyData{
		Name: "Hays", 
		Source: "https://hayscountytx.com/law-enforcement/fire-marshal/", 
		Fetcher: burnban.Hays,
		LastGet: seedTime,
		CachedBan: "",
	}
	db["travis"] = CountyData{
		Name: "Travis", 
		Source: "https://www.traviscountytx.gov/fire-marshal/burn-ban", 
		Fetcher: burnban.Travis,
		LastGet: seedTime,
		CachedBan: "",
	}
	db["presidio"] = CountyData{
		Name: "Presidio", 
		Source: "https://www.co.presidio.tx.us/", 
		Fetcher: burnban.Presidio,
		LastGet: seedTime,
		CachedBan: "",
	}

	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/template", func(c *gin.Context) {
		c.HTML(http.StatusOK, "off.tmpl", gin.H{
			"link": "https://google.com",
			"county": "Here",
		})
	})

	r.GET("/", func(c *gin.Context) {
		// TODO: return list of all counties
		c.HTML(http.StatusOK, "notfound.tmpl", gin.H{})
	})
	
// 	var counties = []countyData{
//     {ID: "comal", Name: "Comal", Source: "John Coltrane"},
//     {ID: "hays", Name: "Hays", Source: "John Coltrane"},
    
// }

	// counties["hays"] = 
	r.GET("/county/:county", func(c *gin.Context) {
		now := time.Now().UTC()
		county := c.Params.ByName("county")
		// If county doesn't exist, return error
		value, ok := db[county]
		log.Println("Getting county ", county, ok)
		if !ok {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"county": county,
				"error": "Some error occurred",
			})
			return;
		}

		ban, err := value.Fetcher(value.Source)
		if err != nil || ban == "" {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"error": err,
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
	})

	// https://github.com/gin-gonic/examples/blob/master/basic/main.go
	// Get user value
	// var db = make(map[string]string)
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := db[user]
	// 	if ok {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")
		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	
	return r
}

func main() {
	log.Println("Running program")
	
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r.Run(":" + port);
	// OLD
	// router := gin.Default()
	// router.GET("/", )

	// router.Run("localhost:8080")
}
