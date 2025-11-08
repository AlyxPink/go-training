package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// Simple C functions for demonstration
int add_numbers(int a, int b) {
    return a + b;
}

// String manipulation in C
char* to_uppercase(const char* str) {
    if (str == NULL) return NULL;

    size_t len = strlen(str);
    char* result = (char*)malloc(len + 1);
    if (result == NULL) return NULL;

    for (size_t i = 0; i < len; i++) {
        if (str[i] >= 'a' && str[i] <= 'z') {
            result[i] = str[i] - 32;
        } else {
            result[i] = str[i];
        }
    }
    result[len] = '\0';
    return result;
}

// Struct example
typedef struct {
    int id;
    double value;
    char name[50];
} DataPoint;

DataPoint* create_datapoint(int id, double value, const char* name) {
    DataPoint* dp = (DataPoint*)malloc(sizeof(DataPoint));
    if (dp == NULL) return NULL;

    dp->id = id;
    dp->value = value;
    strncpy(dp->name, name, 49);
    dp->name[49] = '\0';
    return dp;
}

void free_datapoint(DataPoint* dp) {
    if (dp != NULL) {
        free(dp);
    }
}

// Array processing
int sum_array(int* arr, int len) {
    int sum = 0;
    for (int i = 0; i < len; i++) {
        sum += arr[i];
    }
    return sum;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Add two integers using C code
// Demonstrates basic CGO function calls and type conversion
func Add(a, b int) int {
	// Convert Go int to C int
	ca := C.int(a)
	cb := C.int(b)

	// Call C function
	result := C.add_numbers(ca, cb)

	// Convert C int back to Go int
	return int(result)
}

// ToUppercase converts string to uppercase using C code
// Demonstrates string conversion and memory management
func ToUppercase(s string) (string, error) {
	// Convert Go string to C string
	// C.CString allocates memory that we must free
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))

	// Call C function
	cResult := C.to_uppercase(cStr)
	if cResult == nil {
		return "", errors.New("C function returned NULL")
	}
	// Important: Free the memory allocated by the C function
	defer C.free(unsafe.Pointer(cResult))

	// Convert C string back to Go string
	result := C.GoString(cResult)

	return result, nil
}

// DataPoint represents a C struct in Go
type DataPoint struct {
	ID    int
	Value float64
	Name  string
}

// NewDataPoint creates a DataPoint using C code
// Demonstrates struct handling and memory management
func NewDataPoint(id int, value float64, name string) (*DataPoint, error) {
	// Convert Go string to C string
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	// Call C function to create struct
	cDataPoint := C.create_datapoint(C.int(id), C.double(value), cName)
	if cDataPoint == nil {
		return nil, errors.New("failed to create datapoint")
	}
	// Free the C struct when done
	defer C.free_datapoint(cDataPoint)

	// Convert C struct to Go struct
	// We copy the data so we can safely free the C memory
	dp := &DataPoint{
		ID:    int(cDataPoint.id),
		Value: float64(cDataPoint.value),
		Name:  C.GoString(&cDataPoint.name[0]),
	}

	return dp, nil
}

// SumArray sums an integer array using C code
// Demonstrates array/slice conversion
func SumArray(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	// Convert Go int slice to C int array
	// Go's int is 64-bit but C's int is 32-bit, so we need to convert
	cArray := make([]C.int, len(numbers))
	for i, v := range numbers {
		cArray[i] = C.int(v)
	}

	// Call C function with converted array
	result := C.sum_array(&cArray[0], C.int(len(cArray)))

	return int(result)
}

// StringProcessor demonstrates CGO string handling
type StringProcessor struct {
	// No fields needed - stateless operations
}

// NewStringProcessor creates a new string processor
func NewStringProcessor() *StringProcessor {
	return &StringProcessor{}
}

// Process a string using C functions
// This demonstrates composition of CGO operations
func (sp *StringProcessor) Process(input string) (string, error) {
	// Use the ToUppercase function which handles CGO
	return ToUppercase(input)
}

// Helper function to convert C string to Go string safely
func cStringToGo(cstr *C.char) (string, error) {
	if cstr == nil {
		return "", errors.New("null C string")
	}
	return C.GoString(cstr), nil
}

func main() {
	// Example usage
	fmt.Println("CGO Examples:")

	// Test Add
	result := Add(5, 3)
	fmt.Printf("5 + 3 = %d\n", result)

	// Test ToUppercase
	upper, err := ToUppercase("hello world")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Uppercase: %s\n", upper)
	}

	// Test DataPoint
	dp, err := NewDataPoint(1, 42.5, "test")
	if err != nil {
		fmt.Printf("Error creating datapoint: %v\n", err)
	} else {
		fmt.Printf("DataPoint: ID=%d, Value=%.2f, Name=%s\n", dp.ID, dp.Value, dp.Name)
	}

	// Test SumArray
	numbers := []int{1, 2, 3, 4, 5}
	sum := SumArray(numbers)
	fmt.Printf("Sum of %v = %d\n", numbers, sum)

	// Test StringProcessor
	processor := NewStringProcessor()
	processed, err := processor.Process("test string")
	if err != nil {
		fmt.Printf("Error processing: %v\n", err)
	} else {
		fmt.Printf("Processed: %s\n", processed)
	}
}
