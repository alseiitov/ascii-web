package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template

func main() {
	fs := http.FileServer(http.Dir("styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	templates = template.Must(template.ParseGlob("*.html"))

	http.HandleFunc("/", asciiWeb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	fmt.Printf("Starting server...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func asciiWeb(w http.ResponseWriter, r *http.Request) {
	indexLogo := "\n\n\n\n            @@@@@@@@@.    &@@@@@@@@@@ \n       &@@@@@@@@@@@@@@@@@@&@@@@@@@@@@ \n     @@@@@@@@@@@@@@@@@@@@@  @@@@@@@@@ \n   @@@@@@@@@@@@@@@@@@@@@@@    @@@@@@@ \n  @@@@@@@@@@@@@@@@@@@@@@@@     @@@@@@ \n  @@@@@@@@@@@@@@@@@@@@@@@@     *@@@@@ \n  @@@@@@@@@@@@@@@@@@@@@@@@      @@@@@ \n  @@@@@@@@@@@@@@@@@@@@@@@@     @@@@@@ \n   @@@@@@@@@@@@@@@@@@@@@@@    *@@@@@@ \n    @@@@@@@@@@@@@@@@@@@@@@   @@@@@@@@ \n      @@@@@@@@@@@@@@@@@@@@ @@@@@@@@@@ \n         #@@@@@@@@@@@@@@  &@@@@@@@@@@ "

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		if err := templates.ExecuteTemplate(w, "index.html", indexLogo); err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		}

	case "POST":
		var lines []string
		word := r.FormValue("textToPrint")
		font := r.FormValue("font")
		file, err := os.Open("fonts/" + font + ".txt")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		file.Close()
		result := generator(word, lines)

		if err := templates.ExecuteTemplate(w, "index.html", result); err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func generator(word string, lines []string) string {
	res := ""
	num := 0
	var newLine bool
	for i := 0; i < 8; i++ {
		for iWord, sWord := range word {
			if word[iWord] == '\\' && iWord+1 < len(word) {
				if word[iWord+1] == 'n' {
					num = iWord
					newLine = true
					break
				}
			}
			for iSym := 32; iSym <= 126; iSym++ {
				if sWord == rune(iSym) {
					res = res + lines[(iSym-32)*8+i]
				}
			}
		}
		res = res + "\n"
	}
	if newLine == true {
		res = res + generator(word[num+2:len(word)], lines)
	}
	return res
}
