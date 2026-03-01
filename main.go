// package Backend_Challenge
package main

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {
	// Return incoming csv as a matrix
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		records := parseMatrix(w, r)
		if records == nil {
			return
		}
		var response string
		for _, row := range records {
			response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
		}
		fmt.Fprint(w, response)
	})
	// Return incoming csv as inverted matrix
	http.HandleFunc("/invert", func(w http.ResponseWriter, r *http.Request) {
		records := parseMatrix(w, r)
		if records == nil {
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
	})
	// Return incoming csv as single flattened array
	http.HandleFunc("/flatten", func(w http.ResponseWriter, r *http.Request) {
		records := parseMatrix(w, r)
		if records == nil {
			return
		}
		// if len(records) == 0 {
		// 	return
		// }
		var response strings.Builder
		numRows := len(records)
		for i, row := range records {
			response.WriteString(strings.Join(row, ","))
			if i < numRows-1 {
				response.WriteString(",")
			}
		}
		fmt.Fprint(w, response.String())
	})
	// Return sum of all values in csv
	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		records := parseMatrix(w, r)
		if records == nil {
			return
		}
		// valid empty matrix case
		if len(records) == 0 {
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
	})
	// Return product of all values in csv
	http.HandleFunc("/multiply", func(w http.ResponseWriter, r *http.Request) {
		records := parseMatrix(w, r)
		if records == nil {
			return
		}
		if len(records) == 0 {
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
	})
	fmt.Println("Server now running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// Parse matrix by forming file, and then reading file into records [][]string
// Validate matrix afterwards by calling checkValidMatrix
func parseMatrix(w http.ResponseWriter, r *http.Request) [][]string {
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "error: no file uploaded. please use the key \"file=@PATH_TO_YOUR_FILE\"", http.StatusBadRequest)
			return nil
		}
		http.Error(w, fmt.Sprintf("error: processing upload: %v", err), http.StatusBadRequest)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// Handle jagged matrices in checkValidMatrix function call
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: parsing csv: %v", err), http.StatusBadRequest)
		return nil
	}
	if !checkValidMatrix(w, records) {
		return nil
	}
	// Return empty string matrix, if records is empty
	if len(records) == 0 {
		return [][]string{}
	}
	return records
}

// Check valid matrix by ensuring:
// Matrix has square shape
// Matrix contains only int-like strings (validate using regex)
func checkValidMatrix(w http.ResponseWriter, records [][]string) bool {
	var intRegexString = regexp.MustCompile(`^-?\d+$`)
	numRows := len(records)
	for rowIndex, row := range records {
		if len(row) != numRows {
			http.Error(w, fmt.Sprintf("error: matrix is not square! row %d has length %d, but there are a total of %d rows!", rowIndex, len(row), numRows), http.StatusBadRequest)
			return false
		}
		for colIndex, cell := range row {
			// Trim whitespace in front of and behind each value
			cleanCell := strings.TrimSpace(cell)
			records[rowIndex][colIndex] = cleanCell
			if !intRegexString.MatchString(cleanCell) {
				http.Error(w, fmt.Sprintf("error: matrix contains non-integer character '%s' at [%d][%d]!", cell, rowIndex, colIndex), http.StatusBadRequest)
				return false
			}
		}
	}
	return true
}

// THINGS TO HANDLE
// EMPTY MATRIX with ALL INPUTS and MATH
// HUGE NUMBERS
// spaces between numbers (leading and trailing)
// HUGE CSV SIZE
// negative numbers, 0 numbers
// non number inputs
// invalid file
// missing character
