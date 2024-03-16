package endpoint_test

import (
	"bytes"
	"e/blog/database"
	"e/blog/endpoint"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestCreateBlog(t *testing.T) {
	w := httptest.NewRecorder()

	router := gin.New()
	c, _ := gin.CreateTestContext(w)

	router.POST("/blog/create", func(c *gin.Context) {
		endpoint.CreateBlog(c)
	})
	c.Request, _ = http.NewRequest(http.MethodPost, "/blog/create", bytes.NewBufferString(`{"author":"Test Author","content":"Sample content","title":"My Blog"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, c.Request)
	endpoint.CreateBlog(c)
	assert.Equal(t, http.StatusCreated, w.Code)

}

func TestGetoneBlog(t *testing.T) {
	
	w := httptest.NewRecorder()

	router := gin.New()
	c, _ := gin.CreateTestContext(w)
	router.GET("/blog/:author", endpoint.GetOneBlog)

	mockBlogData := map[string]endpoint.Blog{
        "1": {Author: "Test Author", Content: "Sample content", ID: "1", Title: "My Blog", Upload: "2022-01-01", Views: 10, Coments: []string{"Comment 1"}},
    }

    c.Params = append(c.Params, gin.Param{Key: "Test Author"})

	db := &database.FireDb{} // Assuming your FirebaseDB struct is exported
	db.SetMockData(mockBlogData)
	endpoint.GetOneBlog(c)
	router.ServeHTTP(w, c.Request)

	// Assert on the HTTP response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
	}

}

func TestGetBlogs(t *testing.T) {
	
	w := httptest.NewRecorder()

	router := gin.New()

    // Set up a route for the GetBlogs handler
    router.GET("/blog/all-blogs", endpoint.GetBlogs)

	req, err := http.NewRequest(http.MethodGet, "/blog/all-blogs", nil)
    if err != nil {
        t.Fatal(err)
    }

	
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}



func TestAddComments(t *testing.T) {

	reqBody := bytes.NewBufferString(`[
		{"blog_title": "Test Blog", "comment": "First comment"},
		{"blog_title": "Test Blog", "comment": "Second comment"}
	]`)
	req, _ := http.NewRequest(http.MethodPost, "/comments/Test Blog", reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	endpoint.AddComments(c)

	// Assert the HTTP response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response message
	expectedResponse := gin.H{"Message": "Comment uploaded"}
	assert.Equal(t, expectedResponse, decodeJSON(w.Body.Bytes()))
	
}
func decodeJSON(data []byte) gin.H {
	var result gin.H
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}
