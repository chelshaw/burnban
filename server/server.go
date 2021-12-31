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

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// Ping test
	r.GET("/comal", func(c *gin.Context) {
		found := burnban.Comal()
		// found := burnban.Comal()
		
		var resultString = "Sorry, we couldn't find an answer"
		if (len(found) > 0) {
			fmt.Println("Result was found")
			fmt.Println(found)
			resultString = found
		} 
		c.String(http.StatusOK, resultString)
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
