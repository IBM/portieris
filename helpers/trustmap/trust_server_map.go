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

package trustmap

import (
	"strings"
)

// TrustServerFn A simple type alias to represent a function that takes image and suffix and returns trust url
type TrustServerFn func(string, string) string

// Identity Returns a configured function that just returns a string. For static trust server hosts.
func Identity(value string) TrustServerFn {
	return func(registryHostname string, imageHostname string) string {
		return value
	}
}

// IBMRegional IBM Sponsored Trust server, depends on the regional part of the docker image hostname.
func IBMRegional(registryHostname string, imageHostname string) string {
	trustSuffix := "bluemix.net:4443"
	return "https://" + strings.TrimSuffix(imageHostname, registryHostname) + trustSuffix
}

// TrustServerMap Easy way to link known registries to their sponsored trust servers
var TrustServerMap = map[string]TrustServerFn{
	"docker.io":   Identity("https://notary.docker.io"),
	"quay.io":     Identity("https://quay.io:4443"),
	"bluemix.net": IBMRegional,
}
