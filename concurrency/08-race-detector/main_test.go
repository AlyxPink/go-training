package main

import (
	"sync"
	"testing"
)

func TestBuggyCounter(t *testing.T) {
	// This test will fail with -race
	c := &BuggyCounter{}
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()
	
	// Note: Even without race detector, count may be wrong
	t.Logf("Counter value: %d (should be 100 if fixed)", c.Value())
}

func TestBuggyMapWriter(t *testing.T) {
	// This will likely panic or show race with -race
	m := BuggyMapWriter()
	
	if len(m) != 10 {
		t.Logf("Map has %d entries, expected 10 (races may cause loss)", len(m))
	}
}

func TestBuggySliceAppend(t *testing.T) {
	// Race detector will catch this
	s := BuggySliceAppend()
	
	if len(s) != 10 {
		t.Logf("Slice has %d elements, expected 10 (races may cause loss)", len(s))
	}
}

func TestBuggyLoopCapture(t *testing.T) {
	// Run with -race to see the issue
	BuggyLoopCapture()
}
