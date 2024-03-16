package endpoint

//Contains all data structure which is used in endpoint
type Blog struct {
	ID      string `json:"id"`
	Author  string `json:"author" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Upload  string `json:"upload_time"`
	Views   int    `json:"views"`
	Coments []string `json:"comment"`
}

type Author_info struct {
	Id           string `json:"author_id"`
	Author_Name  string `json:"author_name"`
	Author_views int    `json:"view_count"`
}

type Comments struct {
	Blog_title string `json:"blog_title"`
	Comment string `json:"comment"`
}
