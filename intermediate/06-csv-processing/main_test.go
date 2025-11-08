package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCSV(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "test.csv")
	
	// Create test CSV
	content := "Name,Age,Salary\nAlice,30,75000\nBob,25,65000\n"
	os.WriteFile(tmpfile, []byte(content), 0644)
	
	records, err := ParseCSV(tmpfile)
	if err != nil {
		t.Fatalf("ParseCSV error: %v", err)
	}
	
	if len(records) != 2 {
		t.Errorf("Got %d records, want 2", len(records))
	}
	
	if records[0].Name != "Alice" {
		t.Errorf("Name = %q, want \"Alice\"", records[0].Name)
	}
	
	if records[0].Salary != 75000 {
		t.Errorf("Salary = %f, want 75000", records[0].Salary)
	}
}

func TestWriteCSV(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "output.csv")
	
	records := []Record{
		{"Alice", 30, 75000},
		{"Bob", 25, 65000},
	}
	
	if err := WriteCSV(tmpfile, records); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}
	
	// Verify by reading back
	parsed, err := ParseCSV(tmpfile)
	if err != nil {
		t.Fatalf("ParseCSV error: %v", err)
	}
	
	if len(parsed) != len(records) {
		t.Errorf("Got %d records, want %d", len(parsed), len(records))
	}
}

func TestFilterRecords(t *testing.T) {
	records := []Record{
		{"Alice", 30, 75000},
		{"Bob", 25, 65000},
		{"Charlie", 35, 85000},
	}
	
	filtered := FilterRecords(records, func(r Record) bool {
		return r.Age > 26
	})
	
	if len(filtered) != 2 {
		t.Errorf("Filtered count = %d, want 2", len(filtered))
	}
}

func TestAverageSalary(t *testing.T) {
	records := []Record{
		{"Alice", 30, 60000},
		{"Bob", 25, 80000},
		{"Charlie", 35, 100000},
	}
	
	avg := AverageSalary(records)
	want := 80000.0
	
	if avg != want {
		t.Errorf("AverageSalary = %f, want %f", avg, want)
	}
}
