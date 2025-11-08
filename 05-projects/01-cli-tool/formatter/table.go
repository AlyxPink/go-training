package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(data interface{}) (string, error) {
	// TODO: Implement table formatting
	// Hint: Use text/tabwriter, handle array of objects
	arr, ok := data.([]interface{})
	if !ok {
		return "", fmt.Errorf("table format requires array, got %T", data)
	}

	if len(arr) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Extract headers from first object
	first, ok := arr[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("table format requires array of objects")
	}

	// Collect keys for consistent ordering
	var keys []string
	for k := range first {
		keys = append(keys, k)
	}

	// Print header
	fmt.Fprintln(w, strings.Join(keys, "\t"))
	
	// Print separator
	seps := make([]string, len(keys))
	for i := range seps {
		seps[i] = "---"
	}
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Print rows
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		
		values := make([]string, len(keys))
		for i, k := range keys {
			values[i] = fmt.Sprintf("%v", obj[k])
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	w.Flush()
	return buf.String(), nil
}
