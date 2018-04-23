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
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AdmissionResponder is a helper for handling admission response creation
// It supports adding and returning multiple errors to the user
type AdmissionResponder struct {
	allowed bool
	errors  []string
	patches []byte
}

// Flush creates the admission response to return
func (a *AdmissionResponder) Flush() *v1beta1.AdmissionResponse {
	if a.allowed && !a.HasErrors() {
		res := &v1beta1.AdmissionResponse{
			Allowed: true,
		}

		if a.patches != nil {
			res.Patch = a.patches
			pt := v1beta1.PatchTypeJSONPatch
			res.PatchType = &pt
		}
		return res
	}
	return &v1beta1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Message: fmt.Sprintf("\n%s", strings.Join(a.errors, "\n")),
		},
	}
}

// HasErrors returns a true if there are errors false if not
func (a *AdmissionResponder) HasErrors() bool {
	return len(a.errors) != 0
}

// SetAllowed sets the admission response to allow the admission
func (a *AdmissionResponder) SetAllowed() {
	a.allowed = true
}

// IsAllowed returns a true if the admission is allowed false if not
func (a *AdmissionResponder) IsAllowed() bool {
	return a.allowed
}

// SetPatch sets the patches to be applied by the api server
func (a *AdmissionResponder) SetPatch(patch []byte) {
	a.patches = patch
}

// Write writes the output of flush to the passed responsewriter
func (a *AdmissionResponder) Write(w http.ResponseWriter, ar v1beta1.AdmissionReview) {
	resp := reviewResponseToByte(a.Flush(), ar)
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}
}

// ToAdmissionResponse adds an error to the response
func (a *AdmissionResponder) ToAdmissionResponse(err error) {
	glog.Error(err)
	a.errors = append(a.errors, err.Error())
}

// StringToAdmissionResponse adds a string as an error to the response
func (a *AdmissionResponder) StringToAdmissionResponse(msg string) {
	a.errors = append(a.errors, msg)
}
