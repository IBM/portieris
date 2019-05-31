// Copyright 2018 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func createSecret(name, namespace, dataKey string, dataValue []byte) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			dataKey: dataValue,
		},
	}
}

func TestWrapper_GetSecretToken(t *testing.T) {
	tests := []struct {
		name       string
		secret     *corev1.Secret
		namespace  string
		secretName string
		registry   string
		wantUser   string
		wantPass   string
		wantErr    bool
	}{
		{
			name: "should return token",
			secret: createSecret("name", "namespace", ".dockerconfigjson",
				[]byte(`{ "auths": { "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "auth-token" } } }`)),
			secretName: "name",
			namespace:  "namespace",
			registry:   "us.icr.io",
			wantUser:   "token",
			wantPass:   "registry-token",
		},
		{
			name: "should return token for old imagePullSecret format",
			secret: createSecret("name", "namespace", ".dockercfg",
				[]byte(`{ "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "auth-token" } }`)),
			secretName: "name",
			namespace:  "namespace",
			registry:   "us.icr.io",
			wantUser:   "token",
			wantPass:   "registry-token",
		},
		{
			name:       "error if secret not found",
			wantErr:    true,
			secret:     createSecret("wrong-name", "namespace", ".dockerconfig", []byte("{}")),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if it fails parsing .dockerconfigjson",
			wantErr:    true,
			secret:     createSecret("name", "namespace", ".dockerconfigjson", []byte("{")),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if it fails parsing .dockercfg",
			wantErr:    true,
			secret:     createSecret("name", "namespace", ".dockercfg", []byte("{")),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if no .dockercfg or .dockerconfigjson key",
			wantErr:    true,
			secret:     createSecret("name", "namespace", ".notdockercfg", []byte("{}")),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if secret does not match the registry",
			wantErr:    true,
			secret:     createSecret("name", "namespace", ".dockercfg", []byte("{}")),
			secretName: "name",
			namespace:  "namespace",
			registry:   "test.registry.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeClientset := k8sfake.NewSimpleClientset(tt.secret)
			w := NewKubeClientsetWrapper(kubeClientset)
			username, password, err := w.GetSecretToken(tt.namespace, tt.secretName, tt.registry, false)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, username)
				assert.Equal(t, tt.wantPass, password)
			}
		})
	}
}
