// Copyright 2018, 2020 Portieris Authors.
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

func TestWrapper_GetSecretKey(t *testing.T) {
	tests := []struct {
		name       string
		secret     *corev1.Secret
		namespace  string
		secretName string
		wantKey    []byte
		wantErr    bool
	}{
		{
			name:       "should return key",
			secret:     createSecret("name", "namespace", "key", []byte(`testkey`)),
			secretName: "name",
			namespace:  "namespace",
			wantKey:    []byte(`testkey`),
		},
		{
			name:       "error if secret not found",
			wantErr:    true,
			secret:     createSecret("wrong-name", "namespace", ".dockerconfig", []byte("{}")),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if no key or .dockerconfigjson key",
			wantErr:    true,
			secret:     createSecret("name", "namespace", "motKey", []byte("{}")),
			secretName: "name",
			namespace:  "namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeClientset := k8sfake.NewSimpleClientset(tt.secret)
			w := NewKubeClientsetWrapper(kubeClientset)
			key, err := w.GetSecretKey(tt.namespace, tt.secretName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantKey, key)
			}
		})
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
			name: "should return token, does not complain about missing auth field",
			secret: createSecret("name", "namespace", ".dockerconfigjson",
				[]byte(`{ "auths": { "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com" } } }`)),
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
			name: "should override username and password with values from auth field if set",
			secret: createSecret("name", "namespace", ".dockercfg",
				[]byte(`{ "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "cGluZzpwb25n" } }`)),
			secretName: "name",
			namespace:  "namespace",
			registry:   "us.icr.io",
			wantUser:   "ping",
			wantPass:   "pong",
		},
		{
			name: "should override username and password with values from auth field if set and padded",
			secret: createSecret("name", "namespace", ".dockercfg",
				[]byte(`{ "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "d2liYmxlOndvYmJsZQ==" } }`)),
			secretName: "name",
			namespace:  "namespace",
			registry:   "us.icr.io",
			wantUser:   "wibble",
			wantPass:   "wobble",
		},
		{
			name: "should override username and password with values from auth field and accept colon in the password",
			secret: createSecret("name", "namespace", ".dockerconfigjson",
				[]byte(`{ "auths": { "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "d2liYmxlOndvYmJsZTp3aXRoY29sb24=" } } }`)),
			secretName: "name",
			namespace:  "namespace",
			registry:   "us.icr.io",
			wantUser:   "wibble",
			wantPass:   "wobble:withcolon",
		},
		{
			name: "should revert to username/password fields if auth field is set but badly formed",
			secret: createSecret("name", "namespace", ".dockerconfigjson",
				[]byte(`{ "auths": { "us.icr.io": { "username": "token", "password": "registry-token", "email": "email@email.com", "auth": "d2liYmxl" } } }`)),
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
			username, password, err := w.GetSecretToken(tt.namespace, tt.secretName, tt.registry)
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

func createSecretBasic(name, namespace, data1 string, data2 string) *corev1.Secret {
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
			data1: []byte(data1),
			data2: []byte(data2),
		},
	}
}

func TestWrapper_GetBasicCredentials(t *testing.T) {
	tests := []struct {
		name       string
		secret     *corev1.Secret
		namespace  string
		secretName string
		wantUser   string
		wantPass   string
		wantErr    bool
	}{
		{
			name:       "should return credentials",
			secret:     createSecretBasic("name", "namespace", "username", "password"),
			secretName: "name",
			namespace:  "namespace",
			wantUser:   "username",
			wantPass:   "password",
		},
		{
			name:       "should return empty no error if no name",
			secret:     createSecretBasic("name", "namespace", "username", "password"),
			secretName: "",
			namespace:  "namespace",
			wantUser:   "",
			wantPass:   "",
		},
		{
			name:       "error if secret not found",
			wantErr:    true,
			secret:     createSecretBasic("wrong-name", "namespace", "username", "password"),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if no username",
			wantErr:    true,
			secret:     createSecretBasic("name", "namespace", "userfoo", "password"),
			secretName: "name",
			namespace:  "namespace",
		},
		{
			name:       "error if no password",
			wantErr:    true,
			secret:     createSecretBasic("name", "namespace", "username", "passwrong"),
			secretName: "name",
			namespace:  "namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeClientset := k8sfake.NewSimpleClientset(tt.secret)
			w := NewKubeClientsetWrapper(kubeClientset)
			username, password, err := w.GetBasicCredentials(tt.namespace, tt.secretName)
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
