GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v test/ | grep -v pkg/apis/)

dep:
	@go get -u github.com/golang/dep/cmd/dep

build-deps: dep
	@dep ensure -v -vendor-only

test-deps: build-deps
	@go get github.com/stretchr/testify/assert
	@go get github.com/golang/lint/golint
	@go get github.com/pierrre/gotestcover
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/onsi/gomega/...

push:
	docker build -t $(HUB)/portieris:$(TAG) .
	docker push $(HUB)/portieris:$(TAG)

test: test-deps
	$(GOPATH)/bin/gotestcover -v -coverprofile=cover.out ${GOPACKAGES}

copyright:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright.sh

copyright-check:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright-check.sh

fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

lint:
	@set -e; for LINE in ${GOPACKAGES}; do golint -set_exit_status=true $${LINE} ; done

vet:
	@set -e; for LINE in ${GOPACKAGES}; do go vet $${LINE} ; done

helm.install.local: push
	-rm $$(pwd)/portieris-0.2.0.tgz
	helm package helm/portieris
	helm install -n portieris $$(pwd)/portieris-0.2.0.tgz --set image.host=$(HUB) --set image.tag=$(TAG)

helm.install:
	-rm $$(pwd)/portieris-0.2.0.tgz
	helm package helm/portieris
	helm install -n portieris $$(pwd)/portieris-0.2.0.tgz

helm.clean:
	-helm/cleanup.sh portieris

e2e:
	-helm package install/helm/portieris
	@go test -v ./test/e2e --helmChart $$(pwd)/portieris-0.2.0.tgz

e2e.local: helm.install.local e2e.quick

e2e.local.ics: helm.install.local e2e.quick.ics

e2e.quick: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.wildcards e2e.quick.generic
	- kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.ics: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.armada e2e.quick.wildcards e2e.quick.generic
	- kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.imagepolicy:
	@go test -v ./test/e2e --no-install --trust-image-policy --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.clusterimagepolicy:
	@go test -v ./test/e2e --no-install --trust-cluster-image-policy --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.wildcards:
	@go test -v ./test/e2e --no-install --wildcards-image-policy --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.armada:
	@go test -v ./test/e2e --no-install --armada --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.generic:
	@go test -v ./test/e2e --no-install --generic --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.helm:
	kubectl apply -f test/helm/tiller-rbac.yaml
	helm init --service-account tiller

e2e.clean: helm.clean


