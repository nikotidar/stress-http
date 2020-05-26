package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var requests int64
var startingTime int64

func kill(target string) {
	for {
		resp, err := http.Get(target)
		if err != nil {
			fmt.Printf("\n[-] unable to perform GET request on target %s\n", target)
			fmt.Printf("\t[+] hint: server is down or you're not connected to the internet\n")
			fmt.Printf("\terror code: %s\n\n", err.Error())
		}

		io.Copy(ioutil.Discard, resp.Body)
		// requests += 1
		requests++
		currentTime := time.Now().Unix()
		if requests%5 == 0 {
			fmt.Printf("\rrequest per second: %d | status code: %d\r", requests/(currentTime-startingTime), resp.StatusCode)
		}
	}
}

func killThreaded(target string, threads int) {
	for x := 0; x < threads; x++ {
		go kill(target)
	}
}

func main() {
	requests = 0
	fmt.Println("[+] initializing http stress test!")
	if len(os.Args) > 2 {
		fmt.Printf("[+] target=%s\n", os.Args[1])
		fmt.Printf("[+] threads=%s\n", os.Args[2])
		fmt.Println("press Ctrl+C to sigint/kill process")
		threadCount, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("[-] second argument, '# of threads' must be an integer")
		}

		startingTime = time.Now().Unix()
		time.Sleep(time.Second * 1)
		killThreaded(os.Args[1], threadCount-1)
		kill(os.Args[1])
	} else {
		fmt.Println("[-] Invalid usage!")
		fmt.Println("[+] {target} {# of threads}")
		fmt.Println("\tnote: include scheme and path, as a GET request will be performed on the given target and path")
	}
}
