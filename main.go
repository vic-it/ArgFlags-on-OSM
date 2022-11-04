package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFile(filename string) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		panic(ferr)
	}

	scanner := bufio.NewScanner(file)
	ctr := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%d: %s\n", ctr, line)
		ctr++
	}
}

func main() {
	readFile("abc.osm")
}
