package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	velocity "c/_DEV/GitHub/koho/velocitylimit"
)

// main
// - parses an input file for data
// - encodes each line as json request
// - each request is run through business logic via velocitylimit package
// - each request is assessed as pass or fail to be ingested to the system
func main() {
	// open file for reading
	file, _ := os.Open("input.txt")
	// read the file
	scanner := bufio.NewScanner(file)
	// create file for writing
	f, err := os.Create("data.txt")
	if err != nil {
		panic(err)
	}

	// parse through the input file
	for scanner.Scan() {
		readLine := scanner.Text() // read the line

		// apply velocity limit check
		if transaction := velocity.Limit(readLine); transaction != "" {
			w := bufio.NewWriter(f)                    // buffer for writing to file
			strTrx := []string{transaction, "\n"}      // append new line character
			strTransaction := strings.Join(strTrx, "") // concatenate the string
			_, err2 := w.WriteString(strTransaction)   // write line to file
			if err2 != nil {
				panic(err2)
			}

			w.Flush() // clear buffer
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// close files
	defer f.Close() // close write file
	file.Close()    // close read file
}
