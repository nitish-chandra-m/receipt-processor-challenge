run_local: 
	go mod tidy && go run cmd/web/*.go

run_docker:
	docker build -t receipts . && docker run -d -p 8080:8080 receipts