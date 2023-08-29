package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"html/template"

)

type Page struct {
	Title  	string
	Body   	[]byte
}

func LoadPage(title string) (*Page, error) {
	fmt.Println("hi")
	body, err := ioutil.ReadFile(title)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func LoadMain(w http.ResponseWriter, r *http.Request) {
	cd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(cd)
	body, err := ioutil.ReadFile("web/main.html")
	if err != nil {
	}
	fmt.Fprintf(w, string(body))
	fmt.Println(string(body))
}



