

# Go-Whatsapp

Send Message Whatsapp with HTTP JSON using Golang



## Installation

Install Dependency Golang

```bash
  go get .
```


## Usage/Examples

Run Program

```shell
go run main.go
```

Scan Barcode result Output

POST http://localhost:33133/send_message

```json
{
    "destination_number" : "6281212121212",
    "message" :"test send WA"
}
```


