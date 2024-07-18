# sqlite3-filelock-test
showcase file locking mechanism of the go-sqlite3 library in action

## Usage

Build and run `first` and `second` from different CLIs in parallel.
They should write to the DB 10k entries each strictly one after the other.

In the end, they read the resulting DB and produce files named `result.first.txt` and `result.second.txt`.
They:

1. Should be euqal
2. Can be verified with `go run verify.go <file_name>` which checks that all the expected entries were read from the DB by both of the apps.
