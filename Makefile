generateAPI:
	rm -rf ./generated
	swagger generate server -P models.Principal -f ./be.yaml -s ./generated/swagger/restapi -m ./generated/swagger/models --exclude-main

build:
	go build -o be ./main.go
