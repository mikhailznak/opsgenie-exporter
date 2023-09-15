GO := GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go
IMAGE_NAME := opsgenie-exporter
BINARY := exporter
SRC_PATH := "cmd/exporter"
BACK_FROM_SRC_PATH := "../../"

.PHONY: build
build:
	go mod tidy
	go mod download
	$(GO) build -C $(SRC_PATH) -o $(BACK_FROM_SRC_PATH)$(BINARY)

.PHONY: clean
clean:
	rm $(BINARY)

.PHONY: build-docker
build-docker:
	docker build -t $(IMAGE_NAME):local .

.PHONY: build-docker-ext push-docker-ext pipe-docker-ext
build-docker-ext:
	docker build -t ${REPOSITORY}/$(IMAGE_NAME):${TAG} .
push-docker-ext:
	docker push ${REPOSITORY}/$(IMAGE_NAME):${TAG}
pipe-docker-ext: build-docker-ext push-docker-ext

.PHONY: run-docker
run-docker:
	docker run --rm --network host -e API_KEY=${API_KEY} $(IMAGE_NAME):local
