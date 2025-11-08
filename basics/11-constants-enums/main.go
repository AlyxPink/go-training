package main

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusCompleted
)

// String returns string representation of Status
// TODO: Implement String() method for Status type
func (s Status) String() string {
	panic("not implemented")
}

func main() {}
