// Copyright 2017 Istio Authors
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

// This file describes the abstract model of services (and their instances) as
// represented in Istio. This model is independent of the underlying platform
// (Kubernetes, Mesos, etc.). Platform specific adapters found populate the
// model object with various fields, from the metadata found in the platform.
// The platform independent proxy code uses the representation in the model to
// generate the configuration files for the Layer 7 proxy sidecar. The proxy
// code is specific to individual proxy implementations

package gateway

import (
	"istio.io/api/networking/v1alpha3"

	"istio.io/istio/pkg/config/protocol"
)

// IsTLSServer returns true if this server is non HTTP, with some TLS settings for termination/passthrough
func IsTLSServer(server *v1alpha3.Server) bool {
	if server.Tls != nil && !protocol.Parse(server.Port.Protocol).IsHTTP() {
		return true
	}
	return false
}

// IsHTTPServer returns true if this server is using HTTP or HTTPS with termination
func IsHTTPServer(server *v1alpha3.Server) bool {
	p := protocol.Parse(server.Port.Protocol)
	if p.IsHTTP() {
		return true
	}

	if p == protocol.HTTPS && server.Tls != nil && !IsPassThroughServer(server) {
		return true
	}

	return false
}

// IsPassThroughServer returns true if this server does TLS passthrough (auto or manual)
func IsPassThroughServer(server *v1alpha3.Server) bool {
	if server.Tls == nil {
		return false
	}

	if server.Tls.Mode == v1alpha3.Server_TLSOptions_PASSTHROUGH ||
		server.Tls.Mode == v1alpha3.Server_TLSOptions_AUTO_PASSTHROUGH {
		return true
	}

	return false
}
