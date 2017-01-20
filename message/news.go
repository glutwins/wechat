package message

//News 图文消息
type News struct {
	CommonToken

	ArticleCount int        `xml:"ArticleCount"`
	Articles     []*article `xml:"Articles>item,omitempty"`
}

//NewNews 初始化图文消息
func NewNews() *News {
	news := new(News)
	news.MsgType = MsgTypeNews
	return news
}

//Article 单篇文章
type article struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicURL      string `xml:"PicUrl,omitempty"`
	URL         string `xml:"Url,omitempty"`
}

func (msg *News) AddArticle(title, description, picURL, url string) {
	if msg.ArticleCount >= 8 {
		return
	}
	art := new(article)
	art.Title = title
	art.Description = description
	art.PicURL = picURL
	art.URL = url

	msg.Articles = append(msg.Articles, art)
	msg.ArticleCount = len(msg.Articles)
}
