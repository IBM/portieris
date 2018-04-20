// Copyright 2018 IBM
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

package webhook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	fakeController "github.com/IBM/portieris/pkg/controller/fakecontroller"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getTestWebhookServer() *Server {
	return &Server{
		name:       "test",
		controller: &fakeController.Controller{},
	}
}

func TestServer_handleAdmissionRequest(t *testing.T) {

	tests := []struct {
		name            string
		admissionReview admissionv1beta1.AdmissionReview
		allowed         bool
	}{
		{
			name: "Calls controller admit with admissionRequest from http request",
			admissionReview: admissionv1beta1.AdmissionReview{
				Request: &admissionv1beta1.AdmissionRequest{UID: "requestUID"},
			},
			allowed: true,
		},
		{
			name:    "Returns reject if http request doesn't contain an admissionRequest",
			allowed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := getTestWebhookServer()
			bytesIn, _ := json.Marshal(tt.admissionReview)
			req, _ := http.NewRequest("POST", "", bytes.NewBuffer(bytesIn))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.handleAdmissionRequest)
			handler.ServeHTTP(rr, req)

			var reviewOut admissionv1beta1.AdmissionReview
			bytesOut, _ := ioutil.ReadAll(rr.Body)
			json.Unmarshal(bytesOut, &reviewOut)
			assert.Equal(t, tt.allowed, reviewOut.Response.Allowed)
		})
	}
}

func Test_reviewResponseToByte(t *testing.T) {
	tests := []struct {
		name                string
		admissionResponse   *admissionv1beta1.AdmissionResponse
		admissionReview     admissionv1beta1.AdmissionReview
		wantAdmissionReview admissionv1beta1.AdmissionReview
	}{
		{
			name:              "Review response matches inputted review with matching request/response UID",
			admissionResponse: &admissionv1beta1.AdmissionResponse{UID: "responseUID", Allowed: true},
			admissionReview: admissionv1beta1.AdmissionReview{
				Request: &admissionv1beta1.AdmissionRequest{UID: "requestUID"},
			},
			wantAdmissionReview: admissionv1beta1.AdmissionReview{
				Response: &admissionv1beta1.AdmissionResponse{UID: "requestUID", Allowed: true},
			},
		},
		{
			name:              "Review response object and oldobject are reset",
			admissionResponse: &admissionv1beta1.AdmissionResponse{},
			admissionReview: admissionv1beta1.AdmissionReview{
				Request: &admissionv1beta1.AdmissionRequest{
					UID:       "requestUID",
					Object:    runtime.RawExtension{Raw: []byte("SomeBytes")},
					OldObject: runtime.RawExtension{Raw: []byte("SomeBytes")},
				},
			},
			wantAdmissionReview: admissionv1beta1.AdmissionReview{
				Response: &admissionv1beta1.AdmissionResponse{UID: "requestUID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := reviewResponseToByte(tt.admissionResponse, tt.admissionReview)
			var review admissionv1beta1.AdmissionReview
			json.Unmarshal(bytes, &review)
			assert.Equal(t, tt.wantAdmissionReview, review)
		})
	}
}
