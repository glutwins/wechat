package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

//PostFile 上传文件
func PostFile(fieldname, filename, uri string) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(fields, uri)
}

//MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := http.Post(uri, contentType, bodyBuf)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}
