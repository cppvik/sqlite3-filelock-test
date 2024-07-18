# Makefile

# Define Go compiler
GO := go

# Define source files and output binaries
FIRST_SRC := first.go
SECOND_SRC := second.go
VERIFY_SRC := verify.go

FIRST_BIN := first
SECOND_BIN := second
VERIFY_BIN := verify

# Build targets
.PHONY: all build_first build_second build_verify run verify clean

all: build_first build_second build_verify

build_first:
	$(GO) build -o $(FIRST_BIN) $(FIRST_SRC)

build_second:
	$(GO) build -o $(SECOND_BIN) $(SECOND_SRC)

build_verify:
	$(GO) build -o $(VERIFY_BIN) $(VERIFY_SRC)

# Run targets
run: build_first build_second clean
	./$(FIRST_BIN) &
	./$(SECOND_BIN)

verify:
	./$(VERIFY_BIN) result.first.txt
	./$(VERIFY_BIN) result.second.txt

clean:
	rm -f *.txt *.db | true
# Clean up binaries
clean-all:
	rm -f $(FIRST_BIN) $(SECOND_BIN) $(VERIFY_BIN) *.txt *.db
