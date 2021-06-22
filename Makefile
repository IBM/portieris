GOFILES=$(shell find . -type f -name '*.go' -not -path "./code-generator/*")
GOPACKAGES=$(shell go list ./... | grep -v test/ | grep -v pkg/apis/)

VERSION=v0.10.3
TAG=$(VERSION)
GOTAGS='containers_image_openpgp'

.PHONY: test

image: 
	docker build --build-arg VERSION=$(VERSION) -t portieris:$(TAG) .

push:
	docker tag portieris:$(TAG) $(HUB)/portieris:$(TAG)
	docker push $(HUB)/portieris:$(TAG)

test-deps:
	@go get golang.org/x/lint/golint

alltests: test-deps fmt lint vet copyright-check test

test:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/makeTest.sh "${GOPACKAGES}" ${GOTAGS}

copyright:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright.sh

copyright-check:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright-check.sh

fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

lint:
	@set -e; for LINE in ${GOPACKAGES}; do golint -set_exit_status=true $${LINE} ; done

vet:
	@set -e; for LINE in ${GOPACKAGES}; do go vet --tags $(GOTAGS) $${LINE} ; done

helm.package:
	-rm $$(pwd)/portieris-$(VERSION).tgz
	helm package helm/portieris

helm.install.local: helm.package
	-kubectl create ns portieris
	-kubectl get secret $(PULLSECRET) -o yaml | sed 's/namespace: default/namespace: portieris/' | kubectl create -f - 
	helm install -n portieris portieris $$(pwd)/portieris-$(VERSION).tgz --set image.host=$(HUB) --set image.tag=$(TAG) --set image.pullSecret=$(PULLSECRET)

helm.install: helm.package
	helm install -n portieris portieris $$(pwd)/portieris-$(VERSION).tgz

helm.clean:
	-helm/cleanup.sh portieris

e2e:
	-helm package helm/portieris
	@go test -v ./test/e2e

e2e.local: helm.install.local e2e.quick

e2e.local.ics: helm.install.local e2e.quick.ics

e2e.quick: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.wildcards e2e.quick.generic e2e.quick.simple.imagepolicy e2e.quick.vulnerability
e2e.quick.ics: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.armada e2e.quick.wildcards e2e.quick.generic e2e.quick.simple.imagepolicy
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.imagepolicy:
	@go test -v ./test/e2e --no-install --trust-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.clusterimagepolicy:
	@go test -v ./test/e2e --no-install --trust-cluster-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.wildcards:
	@go test -v ./test/e2e --no-install --wildcards-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.armada:
	@go test -v ./test/e2e --no-install --armada
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.generic:
	go test -v ./test/e2e --no-install --generic
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.simple.imagepolicy:
	@go test -v ./test/e2e --no-install --simple-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.vulnerability:
	@go test -v ./test/e2e --no-install --vulnerability
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | grep -v portieris | awk '{ print $$1 }' | grep -v NAME)

e2e.clean: helm.clean

.PHONY: code-generator regenerate

code-generator:
	git clone https://github.com/kubernetes/code-generator.git --branch v0.17.16 $(GOPATH)/src/k8s.io/code-generator

regenerate:
	$(GOPATH)/src/k8s.io/code-generator/generate-groups.sh all github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/client github.com/IBM/portieris/pkg/apis portieris.cloud.ibm.com:v1


