package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	// arbitrary large file to load in
	fileName = "lotsofdata.csv"
	// runtime tooling used to print memory usage
	mem runtime.MemStats
)

func main() {
	// register all HTTP requests to "/" to
	// trigger the serverResponder function
	http.HandleFunc("/", serverResponder)

	// begin printing memory usage each second
	// continue until the program exits
	go func() {
		for {
			printMemoryUsage()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// start the http server!
	fmt.Println("starting server!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// serverResponder is the function that is called when hitting the webserver!
//
// For each request, you're likely to use somewhere between 40-45kb of memory.
// This is especially important when considering parallel requests!
// Each request is about .04mbs of memory usage. So, if you had a ceiling of 20mbs
// of usage for this app, you're probably looking at about 500 parallel requests max!
// which is a ton for many cases.
func serverResponder(w http.ResponseWriter, r *http.Request) {
	// open the _large_ file
	// don't load it yet!
	f, err := os.Open(fileName)
	defer f.Close()

	if err != nil {
		panic("failed to open file!")
	}

	// by default io.Copy will read/write 32kb chunks at a time!
	// essentially it's copying 32kb from the file (f)
	// to the response (w).
	io.Copy(w, f)
}

// printMemoryUsage will call a garbage colletion (GC)
// GC cleans up all unreferenced memory, essential objects that
// are not longer being pointed to and thus will never be relevant
// to the application.
//
// then it prints the current memory used by the app in kb!
func printMemoryUsage() {
	runtime.GC()
	runtime.ReadMemStats(&mem)
	log.Printf("%s: %dkb", "mem usage", mem.HeapAlloc/1024)
}
