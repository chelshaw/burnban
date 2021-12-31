package main

import (
	"fmt"
	"net/http"

	"example/burnban"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/template", func(c *gin.Context) {
		c.HTML(http.StatusOK, "off.tmpl", gin.H{
			"title": "Main website",
			"county": "Here",
		})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "💩")
	})
	
	r.GET("/comal", func(c *gin.Context) {
		found,on := burnban.Comal()
		
		if found != true {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
		})
		}

		if on {
			c.HTML(http.StatusOK, "on.tmpl", gin.H{
				"county": "Comal",
			})
			} else {
				c.HTML(http.StatusOK, "off.tmpl", gin.H{
					"county": "Comal",
				})
				
		}
		
		// c.HTML(http.StatusOK, "off.tmpl", gin.H{
		// 	"title": "Main website",
		// 	"county": "Comal",
		// })
	})
	
	r.GET("/travis", func(c *gin.Context) {
		found,on := burnban.Travis()
		
		if found != true {
			c.HTML(http.StatusNotFound, "notfound.tmpl", gin.H{
		})
		}

		if on {
			c.HTML(http.StatusOK, "on.tmpl", gin.H{
				"county": "Travis",
			})
			} else {
				c.HTML(http.StatusOK, "off.tmpl", gin.H{
					"county": "Travis",
				})
				
		}
	})
	

	// Get user value
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
	fmt.Println("Running program")
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080");
	// OLD
	// router := gin.Default()
	// router.GET("/", )

	// router.Run("localhost:8080")
}
