// package Backend_Challenge
package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		defer file.Close()
		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		var response string
		for _, row := range records {
			response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
		}
		fmt.Fprint(w, response)
	})
	http.HandleFunc("/invert", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		defer file.Close()
		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		numRows := len(records)
		numCols := len(records[0])
		var response strings.Builder         // Use string builder for more efficient memory usage
		for col := 0; col < numCols; col++ { // perform column-major traversal to invert matrix
			for row := 0; row < numRows; row++ {
				response.WriteString(records[row][col])
				if row < numRows-1 {
					response.WriteString(",") // Add commas to all rows except last row
				}
			}
			response.WriteString("\n")
		}
		fmt.Fprint(w, response.String())
	})
	http.HandleFunc("/flatten", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		defer file.Close()
		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
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
		response.WriteString("\n")
		fmt.Fprint(w, response.String())
	})
	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		defer file.Close()
		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		numRows := len(records)
		numCols := len(records[0])
		var sum int = 0
		for row := 0; row < numRows; row++ {
			for col := 0; col < numCols; col++ {
				num, _ := strconv.Atoi(records[row][col])
				sum += num
			}
		}
		var response string = fmt.Sprintf("%d\n", sum)
		fmt.Fprint(w, response)
	})
	http.HandleFunc("/multiply", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		defer file.Close()
		records, err := csv.NewReader(file).ReadAll()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}
		numRows := len(records)
		numCols := len(records[0])
		var product int = 1
		for row := 0; row < numRows; row++ {
			for col := 0; col < numCols; col++ {
				num, _ := strconv.Atoi(records[row][col])
				product *= num
			}
		}
		var response string = fmt.Sprintf("%d\n", product)
		fmt.Fprint(w, response)
	})
	http.ListenAndServe(":8080", nil)
}

func checkValidMatrix(records [][]string) {
	numRows := len(records)
	numCols := len(records)

}
