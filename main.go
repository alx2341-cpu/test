package main

import (
	"bufio"
	"os"
	"fmt"
	"runtime"
	"strings"
	"net/http"
    "io/ioutil"
    //"syscall"
)

func main() {

	var numCPU = runtime.NumCPU()

	var exit_ chan int = make(chan int,numCPU) 
	var queue_url chan string = make(chan string)
	var search_pattern string = "Go"

	//CREATE GOROUTINES
	for i := 0; i < numCPU; i++ {

		//chain var copy by REF
		go func(queue_url chan string,i int) {

			//REDUCE CPU CASH FLUSH BY FORBID FOR A GOROUTINE TO EXECUTE IN ANOTHER THREAD
			runtime.LockOSThread()

			sum := 0

			//until chan queue_url close
			for url_req := range queue_url {
				response, err := http.Get(url_req)
			    if err != nil {
			        fmt.Printf("%s", err)
			    } else {
			    	defer response.Body.Close()
			        body_byte, err := ioutil.ReadAll(response.Body)
			        if err != nil {
			            fmt.Printf("%s", err)
					}
					body_str := string(body_byte[:])
					//fmt.Printf("syscall.Gettid routine%d: %v\n",i,syscall.Gettid())
					a := strings.Count(body_str, search_pattern)
					fmt.Printf("Count for " + url_req + ": %d\n",a)
					sum += a
				}
			}
			exit_ <- sum
		}(queue_url,i)
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
			//fmt.Printf("scan " + scanner.Text()+"\n")
			queue_url <- scanner.Text()
	}

	//break for queue cycles in each goroutines
	close(queue_url)
	val := 0
	for i := 0; i < numCPU; i++ {
		valTmp := <- exit_
		val += valTmp
	}

	fmt.Printf("Total: %d\n",val)
}
