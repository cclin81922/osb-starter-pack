# If the USE_SUDO_FOR_DOCKER env var is set, prefix docker commands with 'sudo'
ifdef USE_SUDO_FOR_DOCKER
	SUDO_CMD = sudo
endif

IMAGE ?= quay.io/osb-starter-pack/servicebroker
TAG ?= $(shell git describe --tags --always)
PULL ?= IfNotPresent

build: ## Builds the starter pack
	go build -i github.com/cclin81922/osb-starter-pack/cmd/servicebroker

test: ## Runs the tests
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux: ## Builds a Linux executable
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o servicebroker-linux --ldflags="-s" github.com/cclin81922/osb-starter-pack/cmd/servicebroker

image: linux ## Builds a Linux based image
	cp servicebroker-linux image/servicebroker
	$(SUDO_CMD) docker build image/ -t "$(IMAGE):$(TAG)"

clean: ## Cleans up build artifacts
	rm -f servicebroker
	rm -f servicebroker-linux
	rm -f image/servicebroker

push: image ## Pushes the image to dockerhub, REQUIRES SPECIAL PERMISSION
	$(SUDO_CMD) docker push "$(IMAGE):$(TAG)"

deploy-sc: ## Deploys service-catalog with helm
    helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
	helm install svc-cat/catalog --name catalog --namespace catalog

remove-sc: ## Removes service-catalog with helm
	helm delete --purge catalog
	kubectl delete ns catalog

deploy-broker: image ## Deploys broker with helm
	helm upgrade --install broker-skeleton --namespace broker-skeleton \
	charts/servicebroker \
	--set image="$(IMAGE):$(TAG)",imagePullPolicy="$(PULL)"

remove-broker: ## Removes broker with helm
	helm delete --purge broker-skeleton

create-ns: ## Creates a namespace
	kubectl create ns test-ns

remove-ns: ## Removes a namespace
	kubectl delete ns test-ns

provision: create-ns ## Provisions a service instance
	kubectl apply -f manifests/service-instance.yaml

unprovision: ## Removes a service instance
	kubectl delete -f manifests/service-instance.yaml

bind: ## Creates a binding
	kubectl apply -f manifests/service-binding.yaml

unbind: ## Removes a binding
	kubectl delete -f manifests/service-binding.yaml

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build test linux image clean push deploy-sc remove-sc deploy-broker remove-broker create-ns remove-ns provision unprovision bind unbind help
