redrive-sqs: cmd/*/*.go internal/*/*.go
	go build -o redrive-sqs cmd/redrive-sqs/main.go
