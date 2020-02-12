package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var templates *template.Template
var indexLogo []byte
var fonts struct {
	Standard   []string
	Shadow     []string
	Thinkertoy []string
}

func main() {
	fs := http.FileServer(http.Dir("styles")) //Serving static files
	// downloadDir := http.FileServer(http.Dir("download"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	http.Handle("/download/", http.StripPrefix("/download/", fs))

	indexLogo, _ = ioutil.ReadFile("./styles/indexlogo.txt") // Alem logo on index page

	templates = template.Must(template.ParseGlob("*.html"))

	fonts.Standard = readToMemory("standard")
	fonts.Shadow = readToMemory("shadow")
	fonts.Thinkertoy = readToMemory("thinkertoy")

	http.HandleFunc("/", asciiWeb)

	fmt.Printf("Listening server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func asciiWeb(w http.ResponseWriter, r *http.Request) {
	//Not found status handler
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	//Request methods handler
	switch r.Method {
	case "GET":
		if err := templates.ExecuteTemplate(w, "index.html", string(indexLogo)); err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		}
	case "POST":
		var input string
		var font string

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		query, err := url.ParseQuery(string(body))
		if err != nil {
			fmt.Println(err.Error())
		}

		for i, v := range query {
			switch i {
			case "textToPrint":
				input = v[0]
			case "font":
				font = v[0]
			default:
				http.Error(w, "400 Bad request", 400)
				return
			}
		}

		if font != "standard" && font != "shadow" && font != "thinkertoy" {
			http.Error(w, "400 Bad request", 400)
			return
		}

		//Writing art to template
		result := generator(input, font)
		file, err := os.Create("./download/art.pdf")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = file.Write([]byte(result))
		if err != nil {
			fmt.Println(err.Error())
		}
		file.Close()
		if err := templates.ExecuteTemplate(w, "index.html", result); err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//ASCII art generator
func generator(input, font string) string {
	var lines []string
	var res string

	switch font {
	case "standard":
		lines = fonts.Standard
	case "shadow":
		lines = fonts.Shadow
	case "thinkertoy":
		lines = fonts.Thinkertoy
	}

	words := strings.Split(input, "\\n")
	for _, word := range words {
		for i := 0; i < 8; i++ {
			for _, char := range word {
				if char > 31 && char < 127 {
					res = res + lines[(int(char)-32)*8+i]
				}
			}
			res += "\n"
		}
	}
	return res
}

//Read fonts to memory
func readToMemory(font string) []string {
	var lines []string
	file, err := os.Open("fonts/" + font + ".txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
