package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name   string
	Age    int
	Salary float64
}

func ParseCSV(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	
	// Read header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}
	
	var records []Record
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		
		age, _ := strconv.Atoi(row[1])
		salary, _ := strconv.ParseFloat(row[2], 64)
		
		records = append(records, Record{
			Name:   row[0],
			Age:    age,
			Salary: salary,
		})
	}
	
	return records, nil
}

func WriteCSV(path string, records []Record) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write header
	writer.Write([]string{"Name", "Age", "Salary"})
	
	// Write records
	for _, r := range records {
		writer.Write([]string{
			r.Name,
			strconv.Itoa(r.Age),
			fmt.Sprintf("%.0f", r.Salary),
		})
	}
	
	return nil
}

func FilterRecords(records []Record, pred func(Record) bool) []Record {
	var result []Record
	for _, r := range records {
		if pred(r) {
			result = append(result, r)
		}
	}
	return result
}

func AverageSalary(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}
	
	var sum float64
	for _, r := range records {
		sum += r.Salary
	}
	
	return sum / float64(len(records))
}

func main() {
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
