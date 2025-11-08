package main

import "testing"

func TestPingPong(t *testing.T) {
	ping := make(chan string)
	pong := make(chan string)
	
	go Ping(ping, "hello")
	go Pong(ping, pong)
	
	result := <-pong
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestBufferedChannel(t *testing.T) {
	ch := make(chan int, 5)
	values := []int{1, 2, 3, 4, 5}
	
	go BufferedWriter(ch, values)
	
	results := SafeReceiver(ch)
	if len(results) != len(values) {
		t.Errorf("Expected %d values, got %d", len(values), len(results))
	}
}
