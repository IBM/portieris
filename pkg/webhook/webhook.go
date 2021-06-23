// Copyright 2018, 2021 Portieris Authors.
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
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"

	"github.com/IBM/portieris/pkg/controller"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

// Server is the admission webhook server
type Server struct {
	name string
	// Request multiplexer
	mux *http.ServeMux
	// Controller for this webhook
	controller controller.Interface

	serverCert, serverKey []byte
}

// NewServer creates a new admission webhook server with the passed controller handling the admissions
func NewServer(name string, ctrl controller.Interface, cert, key []byte) *Server {

	return &Server{
		name:       name,
		mux:        http.NewServeMux(),
		controller: ctrl,
		serverCert: cert,
		serverKey:  key,
	}
}

// HandleAdmissionRequest handles an incoming request and calls the controllers admit function
// It writes the response from the Admit to the response writer
func (s *Server) HandleAdmissionRequest(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var admissionReview admissionv1.AdmissionReview
	responder := &AdmissionResponder{}
	deserializer := codec.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &admissionReview); err != nil {
		responder.ToAdmissionResponse(fmt.Errorf("bad request, unable to decode admission review body: %v", err))
		responder.Write(w, admissionReview)
		return
	}
	if admissionReview.Request == nil {
		responder.ToAdmissionResponse(errors.New("bad request, no admission request present"))
		responder.Write(w, admissionReview)
		return
	}
	admissionResponse := s.controller.Admit(admissionReview.Request)
	w.Write(reviewResponseToByte(admissionResponse, admissionReview))
}

// HandleLiveness responds to a Kubernetes Liveness probe
// Fail this request if Kubernetes should restart this instance
func (s *Server) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
}

// HandleReadiness responds to a Kubernetes Readiness probe
// Fail this request if this instance can't accept traffic, but Kubernetes shouldn't restart it
func (s *Server) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
}

// Run starts the server
func (s *Server) Run() {
	flag.Parse()
	// TODO: Use mutual tls after we agree on what cert the apiserver should use.
	// glog.Info("Creating TLS config...")
	// cert := getAPIServerCert(clientset)
	// apiserverCA := x509.NewCertPool()
	// apiserverCA.AppendCertsFromPEM(cert)
	certs, err := tls.X509KeyPair(s.serverCert, s.serverKey)
	if err != nil {
		panic(fmt.Sprintf("unable to load certs: %v", err))
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certs},
		MinVersion:   tls.VersionTLS12,
		// TODO: Use mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth determines the server's policy for TLS Client Authentication. The default is NoClientCert.
		// ClientCAs:    apiserverCA,
		// ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientAuth: tls.NoClientCert,
	}
	s.mux.HandleFunc("/admit", s.HandleAdmissionRequest)
	s.mux.HandleFunc("/health/liveness", s.HandleLiveness)
	s.mux.HandleFunc("/health/readiness", s.HandleReadiness)
	port := "8000"
	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		Handler:   s.mux,
		TLSConfig: tlsConfig,
	}
	glog.Infof("Starting %v Webhook on port %s...", s.name, port)
	server.ListenAndServeTLS("", "")
}

func reviewResponseToByte(admissionResponse *admissionv1.AdmissionResponse, admissionReview admissionv1.AdmissionReview) []byte {
	response := admissionv1.AdmissionReview{}
	if admissionResponse != nil {
		response.TypeMeta = admissionReview.TypeMeta
		response.Response = admissionResponse
		if admissionReview.Request != nil {
			response.Response.UID = admissionReview.Request.UID
			// reset the Object and OldObject, they are not needed in a response.
			admissionReview.Request.Object = runtime.RawExtension{}
			admissionReview.Request.OldObject = runtime.RawExtension{}
		}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
		responder := &AdmissionResponder{}
		responder.ToAdmissionResponse(fmt.Errorf("bad request, unable to decode admission review body: %v", err))
		resp = reviewResponseToByte(responder.Flush(), response)
	}
	return resp
}
