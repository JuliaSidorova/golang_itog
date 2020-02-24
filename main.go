package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

var strToFind string = "go"

const k = 5

//----------------------getCount-------------------------------
func getCount(url string, wg *sync.WaitGroup, c chan int) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	kolvo := strings.Count(string(body), strToFind)
	fmt.Println("Count for", url, "-", kolvo)
	runtime.Gosched()
	c <- kolvo
}

//---------------------------main--------------------------
func main() {

	totalKolvo := 0
	wg := &sync.WaitGroup{}
	c := make(chan int, k)

	file, err := os.Open("url.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go getCount(scanner.Text(), wg, c)
		tmp := <-c
		totalKolvo = totalKolvo + tmp
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total: ", totalKolvo)
	//
	wg.Wait()
}

//---------------------------------------------------------------
