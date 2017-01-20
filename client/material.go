package client

import (
	"encoding/json"

	"github.com/glutwins/wechat/util"
)

//Article 永久图文素材
type Article struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
}

//AddNews 新增永久图文素材
func (c *Client) AddNews(articles []*Article) (string, error) {
	var res Material
	if err := c.postJsonUrlFormat(WachatReq{"articles": articles}, &res, addNewsURL); err != nil {
		return "", err
	}

	return res.MediaID, res.Error()
}

// AddMaterial 上传永久性素材（处理视频需要单独上传）
func (c *Client) AddMaterial(mediaType MediaType, filename string) (string, string, error) {
	url, err := c.formatUrlWithAccessToken(addMaterialURL, mediaType)
	if err != nil {
		return "", "", err
	}

	response, err := util.PostFile("media", filename, url)
	if err != nil {
		return "", "", err
	}
	var resMaterial Material
	if err = json.Unmarshal(response, &resMaterial); err != nil {
		return "", "", err
	}
	return resMaterial.MediaID, resMaterial.URL, resMaterial.Error()
}

//AddVideo 永久视频素材文件上传
func (c *Client) AddVideo(filename, title, introduction string) (string, string, error) {
	req := WachatReq{"title": title, "introduction": introduction}
	uri, err := c.formatUrlWithAccessToken(addMaterialURL, MediaTypeVideo)
	if err != nil {
		return "", "", err
	}

	fieldValue, _ := json.Marshal(req)

	fields := []util.MultipartFormField{
		{
			IsFile:    true,
			Fieldname: "video",
			Filename:  filename,
		},
		{
			IsFile:    false,
			Fieldname: "description",
			Value:     fieldValue,
		},
	}

	response, err := util.PostMultipartForm(fields, uri)
	if err != nil {
		return "", "", err
	}

	var resMaterial Material
	if err = json.Unmarshal(response, &resMaterial); err != nil {
		return "", "", err
	}
	return resMaterial.MediaID, resMaterial.URL, resMaterial.Error()
}

//DeleteMaterial 删除永久素材
func (c *Client) DeleteMaterial(mediaID string) error {
	var result CommonResp
	if err := c.postJsonUrlFormat(WachatReq{"media_id": mediaID}, &result, delMaterialURL); err != nil {
		return err
	}
	return result.Error()
}

//MediaUpload 临时素材上传
func (c *Client) MediaUpload(mediaType MediaType, filename string) (Material, error) {
	var media Material
	if uri, err := c.formatUrlWithAccessToken(mediaUploadURL, mediaType); err != nil {
		return media, err
	} else if response, err := util.PostFile("media", filename, uri); err != nil {
		return media, err
	} else if err = json.Unmarshal(response, &media); err != nil {
		return media, err
	}

	return media, media.Error()
}

//GetMediaURL 返回临时素材的下载地址供用户自己处理
//NOTICE: URL 不可公开，因为含access_token 需要立即另存文件
func (c *Client) GetMediaURL(mediaID string) (string, error) {
	return c.formatUrlWithAccessToken(mediaGetURL, mediaID)
}

//ImageUpload 图片上传
func (c *Client) ImageUpload(filename string) (string, error) {
	var image Material

	if uri, err := c.formatUrlWithAccessToken(mediaUploadImageURL); err != nil {
		return "", err
	} else if response, err := util.PostFile("media", filename, uri); err != nil {
		return "", err
	} else if err = json.Unmarshal(response, &image); err != nil {
		return "", err
	}
	return image.URL, image.Error()
}
