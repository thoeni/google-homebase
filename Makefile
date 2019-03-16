.PHONY: test
test:
	go test ./... -count=1

lint:
	gometalinter --vendor ./...

.PHONY: zip
zip:
	GOOS=linux GOARCH=amd64 go build -o main ./cmd/aws-lambda/ && zip deployment.zip main && rm main

push:
	aws lambda update-function-code --region eu-west-2 --function-name fmip --zip-file fileb://./deployment.zip

clean:
	rm deployment.zip

.PHONY: deploy
deploy: zip push clean