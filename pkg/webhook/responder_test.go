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

package webhook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdmissionResponder(t *testing.T) {
	t.Run("should include the error in the response", func(t *testing.T) {
		responder := &AdmissionResponder{}
		err := fmt.Errorf("FAKE_ERROR")
		responder.ToAdmissionResponse(err)
		resp := responder.Flush()
		assert.Equal(t, "\nFAKE_ERROR", resp.Result.Message)
		assert.False(t, resp.Allowed)
	})

	t.Run("should include the patches in the response", func(t *testing.T) {
		responder := &AdmissionResponder{}
		patch := []byte("patch")
		responder.SetPatch(patch)
		responder.SetAllowed()
		resp := responder.Flush()
		assert.Equal(t, string(patch), string(resp.Patch))
		assert.True(t, resp.Allowed)
	})
}
