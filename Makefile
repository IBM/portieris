GOFILES=$(shell find . -type f -name '*.go' -not -path "./code-generator/*" -not -path "./pkg/apis/*")
GOPACKAGES=$(shell go list ./... | grep -v test/ | grep -v pkg/apis/)

VERSION=v0.13.7
TAG=$(VERSION)
GOTAGS='containers_image_openpgp'

.PHONY: test nancy push test-deps alltests copyright-check copyright fmt detect-secrets image

portieris:
	CGO_ENABLED=0 go build \
	-ldflags="-X github.com/IBM/portieris/internal/info.Version=$(VERSION)" -a \
	-tags containers_image_openpgp -o portieris ./cmd/portieris

deps.jsonl: portieris
	go version -m -v portieris | (grep dep || true) | awk '{print "{\"Path\": \""$$2 "\", \"Version\": \"" $$3 "\"}"}' > deps.jsonl

nancy: deps.jsonl
	cat deps.jsonl | nancy --skip-update-check --loud sleuth
 
detect-secrets:
	detect-secrets audit .secrets.baseline

image: 
	docker build --build-arg PORTIERIS_VERSION=$(VERSION) -t portieris:$(TAG) .

push:
	docker tag portieris:$(TAG) $(HUB)/portieris:$(TAG)
	docker push $(HUB)/portieris:$(TAG)

test-deps:
	go install golang.org/x/lint/golint@latest

alltests: test-deps fmt lint vet copyright-check test

test:
	./scripts/makeTest.sh "${GOPACKAGES}" ${GOTAGS}

copyright:
	./scripts/copyright.sh

copyright-check:
	./scripts/copyright-check.sh

fmt:
	@gofmt -l ${GOFILES}
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

e2e.quick: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.wildcards e2e.quick.generic e2e.quick.simple.imagepolicy e2e.quick.simple.clusterimagepolicy
e2e.quick.ics: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.armada e2e.quick.wildcards e2e.quick.generic e2e.quick.simple.imagepolicy
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.trust.imagepolicy:
	go test -v ./test/e2e --no-install --trust-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.trust.clusterimagepolicy:
	go test -v ./test/e2e --no-install --trust-cluster-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.wildcards:
	go test -v ./test/e2e --no-install --wildcards-image-policy
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.armada:
	go test -v ./test/e2e --no-install --armada
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.generic:
	go test -v ./test/e2e --no-install --generic
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.simple.clusterimagepolicy:
	go test -v ./test/e2e --no-install --simple-cluster-image-policy
	-kubectl delete namespace secretnamespace
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.quick.simple.imagepolicy:
	-kubectl delete namespace secretnamespace
	go test -v ./test/e2e --no-install --simple-image-policy
	-kubectl delete namespace secretnamespace
	-kubectl delete namespace $$(kubectl get namespaces | grep '[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*-[0-9a-f]*' | awk '{ print $$1 }' )

e2e.clean: helm.clean

.PHONY: code-generator regenerate

code-generator:
	go install k8s.io/code-generator@v0.24.0

regenerate:
	bash $(GOPATH)/pkg/mod/k8s.io/code-generator@v0.24.0/generate-groups.sh all github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/client github.com/IBM/portieris/pkg/apis portieris.cloud.ibm.com:v1


