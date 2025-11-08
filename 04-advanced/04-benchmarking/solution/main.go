package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

func StringConcat(parts []string) string {
	result := ""
	for _, part := range parts {
		result += part
	}
	return result
}

func StringsBuilder(parts []string) string {
	var sb strings.Builder
	for _, part := range parts {
		sb.WriteString(part)
	}
	return sb.String()
}

func BytesBuffer(parts []string) string {
	var buf bytes.Buffer
	for _, part := range parts {
		buf.WriteString(part)
	}
	return buf.String()
}

func MapLookup(m map[int]string, keys []int) []string {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		if val, ok := m[key]; ok {
			result = append(result, val)
		}
	}
	return result
}

func SliceScan(items []string, targets []string) []string {
	result := make([]string, 0)
	for _, target := range targets {
		for _, item := range items {
			if item == target {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

type LargeStruct struct {
	Data [1024]byte
	ID   int
	Name string
}

func ProcessStructByValue(s LargeStruct) int {
	sum := 0
	for _, b := range s.Data {
		sum += int(b)
	}
	return sum
}

func ProcessStructByPointer(s *LargeStruct) int {
	sum := 0
	for _, b := range s.Data {
		sum += int(b)
	}
	return sum
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func MarshalJSON(users []User) ([]byte, error) {
	return json.Marshal(users)
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func MarshalJSONOptimized(users []User) ([]byte, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(users); err != nil {
		return nil, err
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

func main() {
	parts := []string{"Hello", " ", "World", "!"}

	result1 := StringConcat(parts)
	result2 := StringsBuilder(parts)
	result3 := BytesBuffer(parts)

	fmt.Printf("StringConcat: %s\n", result1)
	fmt.Printf("StringsBuilder: %s\n", result2)
	fmt.Printf("BytesBuffer: %s\n", result3)

	m := map[int]string{1: "one", 2: "two", 3: "three"}
	keys := []int{1, 2, 3}
	fmt.Printf("Map lookup: %v\n", MapLookup(m, keys))

	large := LargeStruct{ID: 1, Name: "test"}
	fmt.Printf("By value: %d\n", ProcessStructByValue(large))
	fmt.Printf("By pointer: %d\n", ProcessStructByPointer(&large))

	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30},
	}
	data, _ := MarshalJSON(users)
	fmt.Printf("JSON: %s\n", data)

	dataOpt, _ := MarshalJSONOptimized(users)
	fmt.Printf("JSON Optimized: %s\n", dataOpt)
}
