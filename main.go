package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	fs := http.FileServer(http.Dir("styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	// Alem logo on index page
	indexLogo, _ = ioutil.ReadFile("./styles/indexlogo.txt")

	//Parse temp;ates
	templates = template.Must(template.ParseGlob("*.html"))

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
