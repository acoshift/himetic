GO=go
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE=acoshift/himetic
K8S-RESOURCE=deploy/himetic
K8S-CONTAINER=himetic

help:
	# himetic
	#
	# make setup -- setup project
	# make start -- start server
	# make dev -- start live reload on port 8000
	# make clean -- remove built result
	# make build -- build server
	# make docker -- build then build docker image
	# make deploy -- build, build docker image, then patch k8s resource with new image

setup:
	$(GO) get -u github.com/codegangsta/gin
	yarn global add gulp
	yarn install

start:
	$(GO) run main.go

dev:
	gin -p 8000 -a 8080 -x vendor --all

clean:
	rm -f himetic
	rm -f assets/app-*.css
	rm -f assets/app-*.js

.PHONY: assets
assets:
	gulp

build: clean assets
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o himetic -ldflags '-w -s' main.go

docker: build
	docker build -t $(IMAGE):$(COMMIT_SHA) .
	docker push $(IMAGE):$(COMMIT_SHA)

patch:
	kubectl set image $(K8S-RESOURCE) $(K8S-CONTAINER)=$(IMAGE):$(COMMIT_SHA)

deploy: docker patch

