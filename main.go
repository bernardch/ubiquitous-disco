package main

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/invert", invertHandler)
	http.HandleFunc("/flatten", flattenHandler)
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/multiply", multiplyHandler)
	fmt.Println("Server now running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// Return incoming csv as a matrix
func echoHandler(w http.ResponseWriter, r *http.Request) {
	records, err := parseMatrix(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	fmt.Fprint(w, response)
}

// Return incoming csv as inverted matrix
func invertHandler(w http.ResponseWriter, r *http.Request) {
	records, err := parseMatrix(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	// No need to traverse empty matrix
	if len(records) == 0 {
		return
	}
	numRows := len(records)
	numCols := len(records[0])
	var response strings.Builder                        // Use string Builder for more efficient writing
	for colIndex := 0; colIndex < numCols; colIndex++ { // Column-major traversal to invert matrix
		for rowIndex := 0; rowIndex < numRows; rowIndex++ {
			response.WriteString(strings.TrimSpace(records[rowIndex][colIndex]))
			if rowIndex < numRows-1 {
				response.WriteString(",") // Add commas to all rows except last row
			}
		}
		response.WriteString("\n")
	}
	fmt.Fprint(w, response.String())
}

// Return incoming csv as single flattened array
func flattenHandler(w http.ResponseWriter, r *http.Request) {
	records, err := parseMatrix(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	var response strings.Builder
	numRows := len(records)
	for i, row := range records {
		response.WriteString(strings.Join(row, ","))
		if i < numRows-1 {
			response.WriteString(",")
		}
	}
	if numRows > 0 {
		response.WriteString("\n")
	}
	fmt.Fprint(w, response.String())
}

// Return sum of all values in csv
func sumHandler(w http.ResponseWriter, r *http.Request) {
	records, err := parseMatrix(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	// If empty matrix provided, sum should be 0
	if records == nil {
		fmt.Fprint(w, "0\n")
		return
	}
	numRows := len(records)
	numCols := len(records[0])
	// Use big for potentially massive numbers
	sum := big.NewInt(0)
	cnum := new(big.Int)
	for rowIndex := 0; rowIndex < numRows; rowIndex++ {
		for colIndex := 0; colIndex < numCols; colIndex++ {
			cnum.SetString(records[rowIndex][colIndex], 10)
			sum.Add(sum, cnum)
		}
	}
	var response string = fmt.Sprintf("%s\n", sum)
	fmt.Fprint(w, response)
}

// Return product of all values in csv
func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	records, err := parseMatrix(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
		return
	}
	// If empty matrix provided, product should be 0
	if records == nil {
		fmt.Fprint(w, "0\n")
		return
	}
	numRows := len(records)
	numCols := len(records[0])
	product := big.NewInt(1)
	cnum := new(big.Int)
	for rowIndex := 0; rowIndex < numRows; rowIndex++ {
		for colIndex := 0; colIndex < numCols; colIndex++ {
			cnum.SetString(records[rowIndex][colIndex], 10)
			product.Mul(product, cnum)
		}
	}
	var response string = fmt.Sprintf("%s\n", product)
	fmt.Fprint(w, response)
}

// Parse matrix by forming file, and then reading file into records [][]string
// Validate matrix afterwards by calling checkValidMatrix
func parseMatrix(r *http.Request) ([][]string, error) {
	// Option to limit maximum csv size
	// r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, fmt.Errorf("no file uploaded. please use the key \"file=@/path/matrix.csv\"")
		}
		return nil, fmt.Errorf("processing upload: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// Handle jagged matrices in checkValidMatrix function call for pretty error handling
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %v", err)
	}
	if err := checkValidMatrix(records); err != nil {
		return nil, err
	}
	return records, nil
}

// Check valid matrix by ensuring:
// Matrix has square shape
// Matrix contains only int-like strings (validate using regex)
func checkValidMatrix(records [][]string) error {
	intRegexString := regexp.MustCompile(`^-?\d+$`)
	numRows := len(records)

	for rowIndex, row := range records {
		if len(row) != numRows {
			return fmt.Errorf("matrix is not square! row %d has length %d, but there are a total of %d rows!", rowIndex, len(row), numRows)
		}
		for colIndex, cell := range row {
			// Trim whitespace in front of and behind each value
			cleanCell := strings.TrimSpace(cell)
			records[rowIndex][colIndex] = cleanCell
			if !intRegexString.MatchString(cleanCell) {
				return fmt.Errorf("matrix contains non-integer character '%s' at [%d][%d]!", cell, rowIndex, colIndex)
			}
		}
	}
	return nil
}
