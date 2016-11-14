PWD := $(shell pwd)
PKG := github.com/conseweb/common
MINT_PWD := /Users/mint/develop/gopath/src
IMAGE := ckeyer/obc:dev
INNER_GOPATH := /opt/gopath

dev:
	docker run --rm \
	 --net host \
	 --name commondev \
	 -v $(MINT_PWD):$(INNER_GOPATH)/src \
	 -w $(INNER_GOPATH)/src/github.com/conseweb/common \
	 -v /var/run/docker.sock:/var/run/docker.sock \
	 -it $(IMAGE) bash

test: 
	docker run --rm \
	 --net host \
	 --name common-testing \
	 -v $(PWD):$(INNER_GOPATH)/src/$(PKG) \
	 -v /var/run/docker.sock:/var/run/docker.sock \
	 -w $(INNER_GOPATH)/src/$(PKG) \
	 $(IMAGE) make test-unit

test-unit: govendor
	go test $$(go list ./... |grep -v "vendor")

govendor:
	which govendor||go get -u github.com/kardianos/govendor
	govendor sync