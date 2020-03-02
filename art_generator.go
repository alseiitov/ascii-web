package main

import "strings"

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
