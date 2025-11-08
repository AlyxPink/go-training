package main
import "fmt"

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusCompleted
)

func (s Status) String() string {
	return [...]string{"Pending", "Active", "Completed"}[s]
}

func main() {
	fmt.Println(StatusPending, StatusActive, StatusCompleted)
}
