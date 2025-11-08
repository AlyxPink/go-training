package main

import "fmt"

func Ping(ping chan<- string, msg string) {
	ping <- msg
}

func Pong(ping <-chan string, pong chan<- string) {
	msg := <-ping
	pong <- msg
}

func BufferedWriter(ch chan<- int, values []int) {
	for _, v := range values {
		ch <- v
	}
	close(ch)
}

func SafeReceiver(ch <-chan int) []int {
	var results []int
	for v := range ch {
		results = append(results, v)
	}
	return results
}

func main() {
	fmt.Println("Channels Basics - Solution")
}
