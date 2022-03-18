package main

import (
	"log"
	"net/http"
	"os"

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
}



func setupRouter() *gin.Engine {
	// Seed data
	db := make(map[string]CountyData,10)
	db["comal"] = CountyData{Name: "Comal", Source: "https://www.co.comal.tx.us/Fire_Marshal.htm", Fetcher: burnban.Comal}
	// db["hays"] = CountyData{Name: "Hays", Source: "blah", Fetcher: burnban.Hays}

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
	
	// r.GET("/comal", func(c *gin.Context) {
	// 	found,on := burnban.Comal()
		
	// 	if found != true {
	// 		c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{})
	// 	}

	// 	var template = "off.tmpl"
	// 	if on {
	// 		template = "on.tmpl"
	// 	}
		
	// 	c.HTML(http.StatusOK, template, gin.H{
	// 		"county": "Comal",
	// 		"link": "https://www.co.comal.tx.us/Fire_Marshal.htm",
	// 	})
	// })
	
	r.GET("/travis", func(c *gin.Context) {
		found,on := burnban.Travis()
		
		if found != true {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{})
		}
		var template = "off.tmpl"
		if on {
			template = "on.tmpl"
		}
		
		c.HTML(http.StatusOK, template, gin.H{
			"county": "Travis",
			"link": "https://www.traviscountytx.gov/fire-marshal/burn-ban",
		})
	})
	
	r.GET("/hays", func(c *gin.Context) {
		ban, url, err := burnban.Hays()
		if err != nil || ban == "" {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"error": err,
				"county": "Hays",
				"link": url,
			})
			return;
		}
		var template = "on.tmpl"
		if ban == "OFF" {
			template = "off.tmpl"
		}
		
		c.HTML(http.StatusOK, template, gin.H{
			"county": "Hays",
			"link": url,
		})
	})
	
	r.GET("/presidio", func(c *gin.Context) {
		found,on := burnban.Presidio()
		
		if found != true {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{})
		}
		var template = "off.tmpl"
		if on {
			template = "on.tmpl"
		}
		
		c.HTML(http.StatusOK, template, gin.H{
			"county": "Presidio",
			"link": "http://www.co.presidio.tx.us/",
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
			log.Print(err)
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
				"error": err.Error(),
				"county": "Hays",
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
