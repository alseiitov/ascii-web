package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var templates *template.Template
var indexLogo []byte
var fonts struct {
	Standard   []string
	Shadow     []string
	Thinkertoy []string
}
var Send struct {
	Input  string
	Font   string
	Result string
}

func main() {
	//Static files server and handler
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Alem logo on index page
	indexLogo, _ = ioutil.ReadFile("./static/indexlogo.txt")

	//Parse templates
	templates = template.Must(template.ParseGlob("./static/index.html"))

	//Read font files to memory
	fonts.Standard = readToMemory("standard")
	fonts.Shadow = readToMemory("shadow")
	fonts.Thinkertoy = readToMemory("thinkertoy")

	//Main handler
	http.HandleFunc("/", asciiWeb)

	//Handle custom PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Start server
	fmt.Printf("Listening server at port %v\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
