package simple

import (
	"bytes"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/containers/image/v5/signature"
)

// Verifier is for verifying simple signing
type Verifier interface {
	TransformPolicies(kWrapper kubernetes.WrapperInterface, namespace string, inPolicies []v1beta1.SimpleRequirement) (*signature.Policy, error)
	CreateRegistryDir(storeURL, storeUser, storePassword string) (string, error)
	VerifyByPolicy(imageToVerify string, credentials credential.Credentials, registriesConfigDir string, simplePolicy *signature.Policy) (*bytes.Buffer, error, error)
	RemoveRegistryDir(dirName string) error
}

type verifier struct{}

// NewVerifier creates a new Verfier
func NewVerifier() Verifier {
	return &verifier{}
}
