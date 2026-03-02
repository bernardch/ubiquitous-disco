# League Backend Challenge

## Setup
Prior to beginning, go should already be installed.

To run the web server, please run the following commands:
```
go mod init league-backend-challenge
go run .
```

Now, you can freely interact with the web server by modifying the following command:
```
curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"
```
- replace `'file=@/path/matrix.csv'` with the path to a provided csv in the `valid/` or `invalid/` directory, or a path to your own csv
- replace the `"echo"` command with one of: `echo, invert, flatten, sum, multiply`
## Testing
A script for testing the server implementation against a range of core functionalities, as well as common and uncommon edge cases has been provided. To run this test script, please allow execute permissions on the test script `test.sh` and then run it like so:
```
chmod +x test.sh
./test.sh
```
The test script will:
1. Start the server
2. Test each specified csv against each endpoint available on the server
3. Compare the actual output to the expected output for each case, marking it as either PASS or FAIL

An additional large csv has been provided, but has been omitted from the test script for better clarity and readability of the test ouput.

### Edge Cases / Assumptions
Empty csv - this should be considered valid; and will either print nothing to output, or return 0 for sum/multiply

Large integer values - the web server should be able to return/print out very large integer values, and avoid potential integer overflow issues

Whitespace - input csvs with whitespace surrounding any of its values should be treated as valid, and the webserver should not error out when parsing a csv with this additional whitespace

Invalid matrices - input calls to the web server with input matrices that are not perfectly square, or contain non-integer values should receive a relevant/useful error message

Large csvs - An option to limit the request body size has been commented out in the parseMatrix function, in consideration of potential file/bandwidth processing limits. If processing a document that goes beyond this threshold, an appropriate processing error will be returned.
