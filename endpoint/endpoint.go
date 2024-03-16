package endpoint

import (
	"context"
	"e/blog/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogService interface {
	GetBlogs(c *gin.Context)
	CreateBlog(c *gin.Context)
	GetOneBlog(c *gin.Context)
	AddComments(c *gin.Context)
}


var dab = database.FirebaseDB() //Database refrence which has firebase Client with established Connection

//Getting all blogs from Database
func GetBlogs(c *gin.Context) {
	fmt.Println("Inside")
	// c.Header("Content-Type", "application/json")
	var data map[string]Blog

	//Getting  all Data from Reference in Realtime Database in Firebase and Reference name is Data
	err := dab.NewRef("data").Get(context.Background(), &data)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	blogs := make([]Blog, 0, len(data))
	for _, v := range data {
		blogs = append(blogs, Blog{
			Author:  v.Author,
			Content: v.Content,
			ID:      v.ID,
			Title:   v.Title,
			Upload:  v.Upload,
			Views:   v.Views,
			Coments: v.Coments,
		})
	}
	c.IndentedJSON(http.StatusOK, blogs)
}

//Adding Blog in Database and Creating new Database reference name Author
func CreateBlog(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var add *Blog

	if err := c.ShouldBindJSON(&add); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	add.Upload = time.Now().Format("2006-01-02 15:04")
	blogID := uuid.New().String()
	add.ID = blogID
	// Creating Reference with data/blog Title
	fmt.Println(add)
	err := dab.NewRef("data/"+add.Title).Set(context.Background(), add)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error While uploading Blog": err.Error()})
		return
	}
	
	authorData := Author_info{
		Id:           blogID,
		Author_Name:  add.Author,
		Author_views: 0, // Initialize view count to 0
	}

	//Creating Reference with Author
	if err := dab.NewRef("authors/"+add.Author).Set(context.Background(), authorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save author info to Firebase"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Blog Created"})
}


// Getting One Blog at A time and increasing View Count
func GetOneBlog(c *gin.Context) {
	authorid := c.Param("author")
	fmt.Println(authorid)
	var blogData map[string]Blog
	err := dab.NewRef("data").Get(context.Background(), &blogData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var foundBlog *Blog
	for _, b := range blogData {
		if b.Author == authorid {
			foundBlog = &b
			break
		}
	}
	if foundBlog == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Blog not found for author"})
		return
	}
	//increwasing View and updating dataset
	foundBlog.Views++
	if err := dab.NewRef("data/"+foundBlog.Title+"/views").Set(context.Background(), foundBlog.Views); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
		return
	}

	c.IndentedJSON(http.StatusOK, foundBlog)
	//Increasing Author View Count
	if err := IncreaseAuthorCount(foundBlog.Author); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment author view count"})
		return
	}

}

func IncreaseAuthorCount(author string) error {
	var authorData Author_info
	err := dab.NewRef("authors/"+author).Get(context.Background(), &authorData)
	if err != nil {
		return err
	}
	authorData.Author_views++
	fmt.Println(authorData)
	if err := dab.NewRef("authors/"+author+"/view_count").Set(context.Background(), authorData.Author_views); err != nil {
		return err
	}

	return nil

}

// Adding Comment in Database with Put request
func AddComments(c *gin.Context) {
	upcmnt := c.Param("title")
	fmt.Println(upcmnt)
	var existingBlog Blog
	err := dab.NewRef("data/"+upcmnt).Get(context.Background(), &existingBlog)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(existingBlog)

	var updatedComments []Comments
	if err := c.ShouldBindJSON(&updatedComments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while updating": err.Error()})
		return
	}
	fmt.Println(updatedComments)
	for _, comment := range updatedComments {
		existingBlog.Coments = append(existingBlog.Coments, comment.Comment)
	}
	if err := dab.NewRef("data/"+upcmnt+"/comment").Set(context.Background(), existingBlog.Coments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comments"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Message": "Comment uploaded"})
}
