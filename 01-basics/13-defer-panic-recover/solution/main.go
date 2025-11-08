package main
import (
	"errors"
	"fmt"
)

func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	if b == 0 {
		panic("division by zero")
	}
	return a / b, nil
}

func main() {}
