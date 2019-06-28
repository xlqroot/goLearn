package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, targetUrl string) error {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fmt.Println(filename)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	field,_:=bodyWriter.CreateFormField("aa")
	field.Write([]byte("我爱中国"))

	contentType := bodyWriter.FormDataContentType()

	fmt.Println(contentType)

	bodyWriter.Close()
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)

	fmt.Println(string(resp_body))
	return nil
}
// sample usage
func main() {

	target_url := "http://localhost:8080/upload"

	filename := "./img/a.jpg"

	postFile(filename, target_url)

}