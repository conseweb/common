#PWD := $(shell pwd)
PWD := /Users/mint/develop/gopath/src
IMAGE := ckeyer/obc:dev
INNER_GOPATH := /opt/gopath
dev:
	docker run --rm \
	 --net host \
	 --name commondev \
	 -v $(PWD):$(INNER_GOPATH)/src \
	 -w $(INNER_GOPATH)/src/github.com/conseweb/common \
	 -v /var/run/docker.sock:/var/run/docker.sock \
	 -it $(IMAGE) bash