GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v test/ | grep -v pkg/apis/)

VERSION=simple-alpha
TAG=$(VERSION)

.PHONY: test

image: 
	docker build -t portieris:$(TAG) .

push: image
	docker tag portieris:$(TAG) $(HUB)/portieris:$(TAG)
	docker push $(HUB)/portieris:$(TAG)

test-deps:
	@go get golang.org/x/lint/golint

alltests: test-deps fmt lint vet copyright-check test

test: 
	echo 'mode: atomic' > cover.out
	go list ./... | xargs -n1 -I{} sh -c 'go test --tags containers_image_openpgp -covermode=atomic -coverprofile=cover.tmp {} && tail -n +2 cover.tmp >> cover.out'
	rm cover.tmp

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
	-rm $$(pwd)/portieris-$(VERSION).tgz
	helm package helm/portieris
	helm install -n portieris $$(pwd)/portieris-$(VERSION).tgz --set image.host=$(HUB) --set image.tag=$(TAG)

helm.install:
	-rm $$(pwd)/portieris-$(VERSION).tgz
	helm package helm/portieris
	helm install -n portieris $$(pwd)/portieris-$(VERSION).tgz

helm.clean:
	-helm/cleanup.sh portieris

e2e:
	-helm package helm/portieris
	@go test -v ./test/e2e --helmChart $$(pwd)/portieris-$(VERSION).tgz

e2e.local: helm.install.local e2e.quick

e2e.local.ics: helm.install.local e2e.quick.ics

e2e.quick: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.wildcards e2e.quick.generic
	- kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.ics: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.armada e2e.quick.wildcards e2e.quick.generic
	- kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.imagepolicy:
	@go test -v ./test/e2e --no-install --trust-image-policy --helmChart $$(pwd)/portieris-$(VERSION).tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.clusterimagepolicy:
	@go test -v ./test/e2e --no-install --trust-cluster-image-policy --helmChart $$(pwd)/portieris-$(VERSION).tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.wildcards:
	@go test -v ./test/e2e --no-install --wildcards-image-policy --helmChart $$(pwd)/portieris-$(VERSION).tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.armada:
	@go test -v ./test/e2e --no-install --armada --helmChart $$(pwd)/portieris-$(VERSION).tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.generic:
	go test -v ./test/e2e --no-install --generic --helmChart $$(pwd)/portieris-$(VERSION).tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.helm:
	kubectl apply -f test/helm/tiller-rbac.yaml
	helm init --service-account tiller

e2e.clean: helm.clean

.PHONY: code-generator regenerate

code-generator:
	git clone https://github.com/kubernetes/code-generator.git --branch irelease-1.15 $(GOPATH)/src/k8s.io/code-generator

regenerate:
	$(GOPATH)/src/k8s.io/code-generator/generate-groups.sh all github.com/IBM/portieris/pkg/apis/securityenforcement/client github.com/IBM/portieris/pkg/apis securityenforcement:v1beta1


