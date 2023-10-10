
OUTPUT=go-whatsapp
export GOOS=linux
export CGO_ENABLED=0
# Define the build target
build:
	go build  -o $(OUTPUT) main.go

# Define the clean target
clean:
	rm -f $(OUTPUT)

swag:
	swag init -g main.go
