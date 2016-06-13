NAME := golang-sample
DESC := sample Golang Lambda Proc
PREFIX ?= usr/local
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDDATE := $(shell date -u +"%B %d, %Y")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL := "https://github.com/mchudgins/lambda_proc"
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.buildTime=$(BUILDTIME)' \
           -X 'main.builder=$(BUILDER)' \
           -X 'main.goversion=$(GOVERSION)'

test:
	go test $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

$(NAME).zip: lambda-proc
	-upx -9 lambda-proc
	zip -r $(NAME).zip lambda-proc index.js

deploy: $(NAME).zip
	aws lambda update-function-code --function $(NAME) --zip-file fileb://$(PWD)/$(NAME).zip

lambda-proc: main.go
	GOOS=linux go build -ldflags "$(LDFLAGS)" -o lambda-proc main.go

coverage:
   @-go test -v -coverprofile=cover.out $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)
   @-go tool cover -html=cover.out -o cover.html

benchmark:
	@echo "Running tests..."
	@go test -bench=. $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

create-stack:
	#aws cloudformation create-stack --stack-name sample-lambda-proc --template-body file://${PWD}/cloudFormation.json
	@echo use the 'create-lambda-proc' target

create-lambda-proc:
	aws lambda create-function --function-name $(NAME) \
		--role arn:aws:iam::362788145219:role/lambda_basic_execution \
		--handler index.handler \
		--description "$(DESC)" \
		--runtime nodejs4.3 \
		--zip-file fileb://$(PWD)/$(NAME).zip

invoke:
	aws lambda invoke --function-name $(NAME) \
		--payload '{ "greek" : "oopah" }' \
		/dev/stdout

clean:
	-@rm lambda-proc *.zip
