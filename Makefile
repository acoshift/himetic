GO=go
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE=acoshift/himetic
K8S-RESOURCE=deploy/himetic
K8S-CONTAINER=himetic
BUCKET=bucket

help:
	# himetic
	#
	# make setup -- setup project
	# make start -- start server
	# make dev -- start live reload on port 8000
	# make clean -- remove built result
	# make assets -- build assets
	# make build -- build server
	# make docker -- build then build docker image
	# make deploy -- build, build docker image, then patch k8s resource with new image

setup:
	$(GO) get -u github.com/codegangsta/gin
	npm -g install node-sass, gulp
	yarn install

start:
	$(GO) run main.go

dev:
	gin -p 8000 -a 8080 -x vendor --all

clean:
	rm -f himetic
	rm -f assets/app-*.css
	rm -f assets/app-*.js
	rm -rf .build/

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

deploy-style:
	mkdir -p .build
	node-sass --output-style compressed style/main.scss > .build/style.css
	$(eval style := style.$(shell cat .build/style.css | md5).css)
	mv .build/style.css .build/$(style)
	gsutil -h "Cache-Control:public, max-age=31536000" cp -Z .build/$(style) gs://$(BUCKET)/$(style)
	echo "style.css: https://storage.googleapis.com/$(BUCKET)/$(style)" > .build/static.yaml

deploy: docker deploy-style patch
