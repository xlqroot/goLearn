package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	path2 "path"
	"strconv"
	"time"
)

func main()  {
	http.HandleFunc("/",upload)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func upload(w http.ResponseWriter,r *http.Request)  {
	r.ParseMultipartForm(r.ContentLength)

	fmt.Println(r.MultipartForm)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	dir := "./img/"
	_,err = os.Stat(dir)

	if err != nil {
		os.MkdirAll(dir,066)
	}

	//path := dir+handler.Filename
	path :=dir + strconv.FormatInt(time.Now().Unix(),10) + path2.Ext(handler.Filename)

	fmt.Println(path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//f.Close()

	//fmt.Println(os.RemoveAll("./test"))
}
