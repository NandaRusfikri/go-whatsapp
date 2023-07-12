# Set the compiler to use for building the program
CC=go

# Set the options to use when compiling the program
GOFLAGS=-ldflags="-s -w"

# Set the name of the output binary
OUTPUT=crm-tickets

# Set the list of source files to compile

# Set the GOOS environment variable to "linux"
export GOOS=linux
export CGO_ENABLED=0
# Define the build target
build:
	$(CC) build -a -installsuffix cgo $(GOFLAGS) -o $(OUTPUT) main.go

# Define the clean target
clean:
	rm -f $(OUTPUT)

swag:
	swag init -g main.go
