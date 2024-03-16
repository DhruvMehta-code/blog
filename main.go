package main

//all Imported Libraries
import (
	"e/blog/endpoint"
	
	"e/blog/middleware"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()      //intialize Gin Default Engine to run API
	grp := router.Group("/blog") // Define Path Prefix

	// All Request with their Function
	grp.GET("/:author", endpoint.GetOneBlog)
	grp.POST("/create", endpoint.CreateBlog)

	grp.PUT("/comments/:title", endpoint.AddComments)
	grp.Use(middleware.AnonymousLoginMiddleware()).GET("/all-blogs",endpoint.GetBlogs)

	//Run API with Localhost port
	router.Run("localhost:8080")

}
