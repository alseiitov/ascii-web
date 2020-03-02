package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func asciiWeb(w http.ResponseWriter, r *http.Request) {
	//Not found status handler
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	//Request methods handler
	switch r.Method {
	case "GET":
		Send.Input = ""
		Send.Font = "standard"
		Send.Result = string(indexLogo)
		if err := templates.ExecuteTemplate(w, "index.html", Send); err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		}
	case "POST":
		var input string
		var font string
		var genOrDown string
		// Parse request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		query, err := url.ParseQuery(string(body))
		if err != nil {
			fmt.Println(err.Error())
		}
		//Check request for invalid keys
		for i, v := range query {
			switch i {
			case "textToPrint":
				input = v[0]
			case "font":
				font = v[0]
			case "genOrDown":
				genOrDown = v[0]
			default:
				http.Error(w, "400 Bad request", 400)
				return
			}
		}
		//Check request for invalid keys #2
		if font != "standard" && font != "shadow" && font != "thinkertoy" {
			http.Error(w, "400 Bad request", 400)
			return
		}
		art := generator(input, font) //Generate art
		//Generate and send or serve to download
		switch genOrDown {
		case "generate":
			//Writing art to template
			Send.Input = input
			Send.Font = font
			Send.Result = art
			if err := templates.ExecuteTemplate(w, "index.html", Send); err != nil {
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			}
		case "download":
			//Serving to download
			file := strings.NewReader(art)
			fileSize := strconv.FormatInt(file.Size(), 10)
			w.Header().Set("Content-Disposition", "attachment; filename=art.txt")
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", fileSize)
			file.Seek(0, 0)
			io.Copy(w, file)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
