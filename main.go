package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func makeRequest(waitGroup *sync.WaitGroup, requestType string, domain string) bool {
	testedURL := fmt.Sprintf("http%s://%s/.git/", requestType, domain)
	response, err := http.Get(testedURL)

	if err == nil && response.StatusCode >= 200 && response.StatusCode <= 299 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		if err == nil && strings.Contains(string(body), "<title>Index of /.git</title>") {
			fmt.Println(fmt.Sprintf("The path %s is apparently vulnerable to .git exposure.\n", testedURL))
			defer waitGroup.Done()
			return true
		}

	}

	defer waitGroup.Done()
	return false
}

func main() {
	var pagesPath string

	flag.StringVar(&pagesPath, "d", "", "Path to the pages' file")
	flag.Parse()

	if pagesPath == "" {
		log.Fatal("You need to specify a path to the pages' file")
		os.Exit(1)
	}

	file, err := os.Open(pagesPath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var waitGroup sync.WaitGroup

	for scanner.Scan() {
		domain := scanner.Text()

		waitGroup.Add(1)
		go makeRequest(&waitGroup, "", domain)
		waitGroup.Add(1)
		go makeRequest(&waitGroup, "s", domain)
	}

	waitGroup.Wait()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[*] All the given pages where verified")
}
