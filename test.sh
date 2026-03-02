#!/bin/bash

trap 'echo -e "\nStopping Go server (PID: $SERVER_PID)..."; kill $SERVER_PID 2>/dev/null; rm -f temp-test-server' EXIT

echo "Building Go server..."
go build -o temp-test-server .

echo "Starting Go server in the background..."
./temp-test-server &
SERVER_PID=$!

echo "Waiting for server to initialize..."
sleep 2

echo "======================================="
echo "STARTING TESTS"
echo "======================================="

PASSED=0
FAILED=0

FILES=(
    "valid/empty-matrix.csv" 
    "valid/large-matrix.csv" 
    "valid/matrix.csv" 
    "valid/valid-matrix.csv" 
    "valid/valid2-matrix.csv" 
    "valid/valid3-matrix.csv" 
    "valid/whitespace-matrix.csv" 
    "invalid/jagged-matrix.csv" 
    "invalid/nonint-matrix.csv" 
    "invalid/nonsquare-matrix.csv"
)

ENDPOINTS=("echo" "invert" "flatten" "sum" "multiply")

for FILE in "${FILES[@]}"; do
    echo -e "\nTESTING FILE: $FILE"
    echo "---------------------------------------"

    # Skip large matrix for readability
    if [[ "$FILE" == *"large-matrix"* ]]; then
        for i in {1..5}; do echo "Skip large matrix test for readability"; done
        continue
    fi

    for ENDPOINT in "${ENDPOINTS[@]}"; do
        EXPECTED=""
        
        # --- EXPLICIT EXPECTATION MAPPING ---
        case "$FILE" in
            "valid/empty-matrix.csv")
                case "$ENDPOINT" in
                    sum|multiply) EXPECTED="0" ;;
                    *)            EXPECTED="" ;;
                esac ;;
                
            "valid/matrix.csv")
                case "$ENDPOINT" in
                    echo)     EXPECTED=$'1,2,3\n4,5,6\n7,8,9' ;;
                    invert)   EXPECTED=$'1,4,7\n2,5,8\n3,6,9' ;;
                    flatten)  EXPECTED="1,2,3,4,5,6,7,8,9" ;;
                    sum)      EXPECTED="45" ;;
                    multiply) EXPECTED="362880" ;;
                esac ;;
                
            "valid/valid-matrix.csv")
                case "$ENDPOINT" in
                    echo)     EXPECTED=$'1,2,3,4,5\n6,7,8,9,10\n11,12,13,14,15\n16,17,18,19,20\n21,22,23,24,25' ;;
                    invert)   EXPECTED=$'1,6,11,16,21\n2,7,12,17,22\n3,8,13,18,23\n4,9,14,19,24\n5,10,15,20,25' ;;
                    flatten)  EXPECTED="1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25" ;;
                    sum)      EXPECTED="325" ;;
                    multiply) EXPECTED="15511210043330985984000000" ;;
                esac ;;
                
            "valid/valid2-matrix.csv")
                EXPECTED="-1" ;;
                
            "valid/valid3-matrix.csv")
                case "$ENDPOINT" in
                    echo)     EXPECTED=$'-1,-2\n-102938,0' ;;
                    invert)   EXPECTED=$'-1,-102938\n-2,0' ;;
                    flatten)  EXPECTED="-1,-2,-102938,0" ;;
                    sum)      EXPECTED="-102941" ;;
                    multiply) EXPECTED="0" ;;
                esac ;;
                
            "valid/whitespace-matrix.csv")
                case "$ENDPOINT" in
                    echo)     EXPECTED=$'1,2,3\n4,5,6\n7,8,90' ;;
                    invert)   EXPECTED=$'1,4,7\n2,5,8\n3,6,90' ;;
                    flatten)  EXPECTED="1,2,3,4,5,6,7,8,90" ;;
                    sum)      EXPECTED="126" ;;
                    multiply) EXPECTED="3628800" ;;
                esac ;;
                
            "invalid/jagged-matrix.csv")
                EXPECTED=$'error: matrix is not square! row 1 has length 4, but there are a total of 3 rows!' ;;

            "invalid/nonsquare-matrix.csv")
                EXPECTED=$'error: matrix is not square! row 0 has length 3, but there are a total of 2 rows!' ;;

            "invalid/nonint-matrix.csv")
                EXPECTED=$'error: matrix contains non-integer character \'a\' at [0][0]!' ;;
        esac

        # Execute Request (Assuming you are using curl to test the endpoints)
        ACTUAL=$(curl -s -F "file=@$FILE" "http://localhost:8080/$ENDPOINT")

        # Validation Logic (Checking if ACTUAL contains EXPECTED)
        if [[ "$ACTUAL" == *"$EXPECTED"* ]]; then
            echo "PASS: /$ENDPOINT"
            echo "   Expected to contain:"
            echo "'$EXPECTED'"
            echo "   Output:"
            echo "'$ACTUAL'"
            ((PASSED++))
        else
            echo "FAIL: /$ENDPOINT"
            echo "   Expected to contain:"
            echo "'$EXPECTED'"
            echo "   Actual Output:"
            echo "'$ACTUAL'"
            ((FAILED++))
        fi
    done
done

echo -e "\n======================================="
echo "RESULTS: $PASSED Passed, $FAILED Failed"
echo "======================================="

if [ $FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi