
.PHONY: build
IMAGE:=mlcp:v1alpha1

build: build-binary build-image

build-binary:
	go build -o mlcp

build-image:
	docker build -t $(IMAGE)
