package main

import (
	"bufio"
	"fmt"
	"os"
)

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
