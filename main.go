package main

import (
	"fmt"
	"os"
)

func crashIf(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

// We store the adb output. In the future we'd like
// to be able to save/process the whole adb history
var lines []string

func main() {
	readFilters("filters.yaml")

	// start reading adb output
	adbChan := make(chan string, 10000)
	crashIf(readAdb(adbChan))

	for {
		line := <-adbChan
		lines = append(lines, line)
		if mustPrint(line) {
			fmt.Println(highlightString(line))
		}
	}
}
