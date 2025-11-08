package main
import "strconv"

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func main() {}
