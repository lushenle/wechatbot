.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o wechatbot main.go

.PHONY: docker
docker:
	docker build . -t wechatbot:latest

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: run
run: fmt vet ## Run a wechatbot from your host
	go run ./main.go
