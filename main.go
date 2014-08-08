package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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
var linesMutex sync.Mutex

func main() {
	readFilters("filters.yaml")

	// start reading adb output
	adbChan := make(chan string, 10000)
	crashIf(readAdb(adbChan))

	handler := func(w http.ResponseWriter, r *http.Request) {
		lineNum := 0
		var l string

		for {
			linesMutex.Lock()
			mustSend := lineNum < len(lines)
			if mustSend {
				l = lines[lineNum]
				lineNum++
			}
			linesMutex.Unlock()

			if !mustSend {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			fmt.Fprintln(w, l)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}

	server := http.Server{
		Addr:         ":10001",
		WriteTimeout: 4 * time.Hour,
		Handler:      http.HandlerFunc(handler),
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	for {
		line := <-adbChan
		linesMutex.Lock()
		lines = append(lines, line)
		linesMutex.Unlock()

		if mustPrint(line) {
			fmt.Println(highlightString(line))
		}
	}

}
