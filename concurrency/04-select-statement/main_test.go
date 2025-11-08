package main

import (
	"testing"
	"time"
)

func TestMultiplex(t *testing.T) {
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 1
	ch1 <- 2
	ch2 <- 10
	close(ch1)
	close(ch2)

	count := 0
	for range Multiplex(ch1, ch2) {
		count++
	}

	if count != 3 {
		t.Errorf("Multiplex produced %d values, want 3", count)
	}
}

func TestTimeout(t *testing.T) {
	ch := make(chan string, 1)
	ch <- "data"

	if _, ok := Timeout(ch, 1*time.Second); !ok {
		t.Error("Should have received value before timeout")
	}

	chSlow := make(chan string)
	if _, ok := Timeout(chSlow, 10*time.Millisecond); ok {
		t.Error("Should have timed out")
	}
}

func TestNonBlockingSend(t *testing.T) {
	ch := make(chan int, 1)

	if !NonBlockingSend(ch, 42) {
		t.Error("Should have sent to empty channel")
	}

	if NonBlockingSend(ch, 43) {
		t.Error("Should not have sent to full channel")
	}
}
