APP=iec104
BASE_VERSION=0.0.2
BASE_DEPLOY_VERSION=0.0.1
DOCKER_DIR=build/package
IMAGE_TAG=$(shell git log --pretty=format:"%ad_%h" -1 --date=short)

SERVER_HOST=192.168.0.105
SERVER_PORT=2404
SUB_SERVER_HOST=192.168.0.104
SUB_SERVER_PORT=2404
DEBUG=true

deps:
	go mod download
dev: example/client/main.go
	ENV SERVER_HOST=$(SERVER_HOST) SERVER_PORT=$(SERVER_PORT) SUB_SERVER_HOST=$(SUB_SERVER_HOST) SUB_SERVER_PORT=$(SUB_SERVER_PORT) DEBUG=$(DEBUG) go run example/client/main.go
test:
	go test -cover ./...
	
#docker
base: $(DOCKER_DIR)/base/Dockerfile
	docker build -t 9d77v/$(APP):base-$(BASE_VERSION) -f $(DOCKER_DIR)/base/Dockerfile .
	docker push 9d77v/$(APP):base-$(BASE_VERSION)
base-deploy: $(DOCKER_DIR)/base/Dockerfile.deploy
	docker build -t 9d77v/$(APP):base-deploy-$(BASE_DEPLOY_VERSION) -f $(DOCKER_DIR)/base/Dockerfile.deploy .
	docker push 9d77v/$(APP):base-deploy-$(BASE_DEPLOY_VERSION)
deploy: $(DOCKER_DIR)/Dockerfile test
	docker build -t 9d77v/iec104:$(IMAGE_TAG) -f $(DOCKER_DIR)/Dockerfile .
	docker push 9d77v/iec104:$(IMAGE_TAG)