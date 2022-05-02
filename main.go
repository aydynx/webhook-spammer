package main

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	client fasthttp.Client

	errors int = 0
	sent   int = 0
	rps    int = 0
)

func sendMessage(webhook string, message string) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod("POST")
	req.SetRequestURI(webhook)
	req.SetBody([]byte(fmt.Sprintf("{\"content\": \"%s\"}", message)))

	req.Header.Set("Content-Type", "application/json")

	if err := client.Do(req, res); err != nil {
		errors++
		return
	}

	sent++
}

func rpsCounter() {
	for {
		before := sent
		time.Sleep(time.Second)
		after := sent

		rps = after - before
	}
}

func statusPrinter() {
	for {
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("[+] sent: %v | rps: %v | errors: %v\r", sent, rps, errors)
	}
}

func main() {
	var threads int
	var webhook string
	var message string

	fmt.Print("[>] threads: ")
	fmt.Scanln(&threads)

	fmt.Print("[>] webhook: ")
	fmt.Scanln(&webhook)

	fmt.Print("[>] message: ")
	fmt.Scanln(&message)

	fmt.Print("\n")

	go rpsCounter()
	go statusPrinter()

	for i := 0; i < threads; i++ {
		go func() {
			for {
				sendMessage(webhook, message)
			}
		}()
	}
	select {}
}
