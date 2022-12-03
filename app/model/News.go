package model

// News
// 新闻表
type News struct {
	Id          int    `json:"id"`
	CId         int    `json:"c_id"`
	Lang        string `json:"lang"`
	Title       string `json:"title"`
	Recommend   int    `json:"recommend"`
	Audit       int    `json:"audit"`
	Display     int    `json:"display"`
	Discuss     int    `json:"discuss"`
	Author      string `json:"author"`
	BrowseGrant int    `json:"browse_grant"`
	Keyword     string `json:"keyword"`
	Abstract    string `json:"abstract"`
	Content     string `json:"content"`
	Views       int    `json:"views"`
	CreateTime  int    `json:"create_time"`
	UpdateTime  int    `json:"update_time"`
	Thumbnail   string `json:"thumbnail"`
	Cover       string `json:"cover"`
	Sorts       int    `json:"sorts"`
}

// QueryNewsListPage
// 分页获取新闻列表
func (n *News) QueryNewsListPage(page, limit int) {

}
