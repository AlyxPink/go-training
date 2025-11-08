package main

import (
	// TODO: Uncomment for CSV operations
	// "encoding/csv"
	"fmt"
	// "io"
	"os"
	// "strconv"
)

// Record represents a CSV record
type Record struct {
	Name   string
	Age    int
	Salary float64
}

// ParseCSV reads CSV file and returns records
func ParseCSV(path string) ([]Record, error) {
	// TODO: Open file, create CSV reader, parse records
	return nil, nil
}

// WriteCSV writes records to CSV file
func WriteCSV(path string, records []Record) error {
	// TODO: Create file, CSV writer, write headers and records
	return nil
}

// FilterRecords returns records matching predicate
func FilterRecords(records []Record, pred func(Record) bool) []Record {
	// TODO: Filter records using predicate function
	return nil
}

// AverageSalary calculates average salary from records
func AverageSalary(records []Record) float64 {
	// TODO: Calculate average
	return 0
}

func main() {
	// Example: create sample CSV
	sample := []Record{
		{"Alice", 30, 75000},
		{"Bob", 25, 65000},
		{"Charlie", 35, 85000},
	}
	
	WriteCSV("sample.csv", sample)
	
	records, _ := ParseCSV("sample.csv")
	fmt.Printf("Loaded %d records\n", len(records))
	
	highEarners := FilterRecords(records, func(r Record) bool {
		return r.Salary > 70000
	})
	fmt.Printf("High earners: %d\n", len(highEarners))
	
	avg := AverageSalary(records)
	fmt.Printf("Average salary: %.2f\n", avg)
	
	os.Remove("sample.csv")
}
