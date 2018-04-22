GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOPACKAGES=$(shell go list ./... | grep -v /vendor/)

copyright:
	@${GOPATH}/src/github.com/IBM/portieris/scripts/copyright.sh

test:
	$(GOPATH)/bin/gotestcover -v -coverprofile=cover.out ${GOPACKAGES}

fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

lint:
	@set -e; for LINE in ${GOPACKAGES}; do golint -set_exit_status=true $${LINE} ; done

vet:
	@set -e; for LINE in ${GOPACKAGES}; do go vet $${LINE} ; done

helm.install:
	-rm $$(pwd)/portieris-0.2.0.tgz
	helm package helm/portieris
	helm install -n portieris $$(pwd)/portieris-0.2.0.tgz

helm.clean:
	-helm/cleanup.sh portieris

e2e:
	-helm package install/helm/portieris
	@go test -v ./test/e2e --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz

e2e.quick: e2e.quick.trust.imagepolicy e2e.quick.trust.clusterimagepolicy e2e.quick.armada e2e.quick.wildcards e2e.quick.generic
	- kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.imagepolicy:
	@go test -v ./test/e2e --no-install --trust-image-policy --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.trust.clusterimagepolicy:
	@go test -v ./test/e2e --no-install --trust-cluster-image-policy --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.wildcards:
	@go test -v ./test/e2e --no-install --wildcards-image-policy --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.armada:
	@go test -v ./test/e2e --no-install --armada --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.quick.generic:
	@go test -v ./test/e2e --no-install --generic --kubeconfig $$HOME/.kube/config --helmChart $$(pwd)/portieris-0.2.0.tgz
	-kubectl delete namespace $$(kubectl get namespaces | grep -v ibm | grep -v kube | grep -v default | awk '{ print $$1 }' | grep -v NAME)

e2e.helm:
	kubectl apply -f test/helm/tiller-rbac.yaml
	helm init --service-account tiller

e2e.clean: helm.clean


