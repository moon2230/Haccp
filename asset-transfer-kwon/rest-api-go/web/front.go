package web

import (
	"net/http"
	"io/ioutil"
)


func (setup *OrgSetup) PageMain(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadFile("web/main.html")
	w.Write(body)
}