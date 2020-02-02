module github.com/IBM/portieris

go 1.13

require (
	github.com/Sirupsen/logrus v0.0.0-00010101000000-000000000000 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/containers/image/v5 v5.1.0
	github.com/docker/distribution v2.6.2+incompatible
	github.com/docker/go v1.5.1-1 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/btree v1.0.0 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/kubernetes/apiextensions-apiserver v0.0.0-20181121072900-e8a638592964
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	github.com/theupdateframework/notary v0.6.1
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20181026184759-d1dc89ebaebe
	k8s.io/apiextensions-apiserver v0.0.0-20181026191334-ba848ee89ca3
	k8s.io/apimachinery v0.0.0-20181022183627-f71dbbc36e12
	k8s.io/client-go v0.0.0-20181026185218-bf181536cb4d
	k8s.io/kube-openapi v0.0.0-20181114233023-0317810137be // indirect
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.2.0
