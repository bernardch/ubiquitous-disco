#!/bin/bash

echo "Run Go server in background"
go run . &
SERVER_PID=$!
echo "⏳ Waiting for server to initialize..."
sleep 2
trap "echo -e '\n Stopping Go server (PID: $SERVER_PID)...'; kill $SERVER_PID" EXIT

# FILES AND ENDPOINTS HERE
FILES=("invalid/jagged-matrix.csv" "invalid/nonint-matrix.csv" "invalid/nonint2-matrix.csv" "invalid/nonsquare-matrix.csv" "valid/empty-matrix.csv" "valid/matrix.csv" "valid/valid-matrix.csv" "valid/valid2-matrix.csv" "valid/valid3-matrix.csv")
# "valid/large-matrix.csv"
ENDPOINTS=("echo" "invert" "flatten" "sum" "multiply")

echo "======================================="
echo "🧪 STARTING TESTS"
echo "======================================="

for file in "${FILES[@]}"; do
    # Quick check to make sure you actually created the dummy file
    if [[ ! -f "$file" ]]; then
        echo "File '$file' not found. Skipping."
        continue
    fi

    echo -e "\nTESTING FILE: $file"
    echo "---------------------------------------"
    
    for endpoint in "${ENDPOINTS[@]}"; do
        echo "/localhost:8080/$endpoint"
        curl -s -w "\n" -F "file=@$file" "localhost:8080/$endpoint"
    done
done

echo -e "\nAll tests complete."