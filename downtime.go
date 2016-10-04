package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Parse command line arguments
	outFilePtr := flag.String("f", "", "Output file path, Logs to stdout if not provided")
	pingFreqPtr := flag.Int("p", 5, "Ping frequency (Seconds)")
	flag.Parse()
	args := flag.Args()

	// Check url provided
	if len(args) < 1 || args[0] == "" || len(args[0]) < 4 {
		fmt.Println("No url provided or url is too short")
		return
	}

	// Format url in neccessary
	url := args[0]
	if url[:4] != "http" {
		url = "http://" + url
	}

	// Get our writing and checking components
	writer := getWriter(*outFilePtr)
	httpChecker := getHTTPChecker(writer, url)
	writer("INFO :: Starting downtime check of url: " + url)

	// Start loop to check http every so often
	quitChan := make(chan bool)
	t := time.NewTicker(time.Second * time.Duration(*pingFreqPtr))
	for {
		select {
		case <-t.C:
			go httpChecker()
		case <-quitChan:
			t.Stop()
			return
		}
	}
}

// Get a http checker function that tracks details such as
// fail start and fail end.
func getHTTPChecker(w writer, url string) func() {
	status := true
	failStart := time.Now()
	timeout := time.Duration(4 * time.Second)
	client := http.Client{Timeout: timeout}

	return func() {
		resp, err := client.Head(url)
		failing := (err != nil || resp.StatusCode > 250)

		// Failure start
		if failing && status {
			status = false
			failStart = time.Now()
			message := "WARN :: Requests failing. "
			if resp != nil {
				message += fmt.Sprintf("Status code %d. ", resp.StatusCode)
			}
			if err != nil {
				message += fmt.Sprintf("Error Message: ", err.Error())
			}
			w(message)
			return
		}

		// Success start after failure
		if !failing && !status {
			status = true
			seconds := int(time.Now().Sub(failStart).Seconds())
			hours := seconds / 3600
			seconds -= hours * 3600
			mins := seconds / 60
			seconds -= mins * 60
			message := "RESULT :: Failed to connect to %s for %d hours, %d minutes and %d seconds."
			w(fmt.Sprintf(message, url, hours, mins, seconds))
			return
		}

	}
}

func checkHTTP(w writer, url string) {

}

// writer is an interface for functions than can write output
type writer func(string)

// Standardise the output format with a timestamp
func formatLogText(text string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return "[" + t + "] " + text
}

// getWriter returns a new writer depending on the given user options.
// A stdout writer is returned if no outfile path is provided
// otherwise a file is written to
func getWriter(outputFile string) writer {

	outputFileSet := (outputFile != "")

	if !outputFileSet {
		return func(text string) {
			fmt.Println(formatLogText(text))
		}
	}

	path, err := filepath.Abs(outputFile)
	checkErr(err)

	return func(text string) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0666)
		checkErr(err)
		defer file.Close()
		writer := bufio.NewWriter(file)
		fmt.Fprintln(writer, formatLogText(text))
		writer.Flush()
	}

}

// cherrErr checks an error item and stops the application on error.
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
