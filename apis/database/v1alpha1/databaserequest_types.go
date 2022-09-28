/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func init() {
	SchemeBuilder.Register(&DatabaseRequest{}, &DatabaseRequestList{})
}

type DatabaseRequestSpec struct {

	// DatabaseClassName name of the DatabaseClass
	DatabaseClassName string `json:"databaseClassName,omitempty"`

	// Engine name
	Engine string `json:"engine"`

	// Name of a database object.
	// If unspecified, then a new Database will be dynamically provisioned
	// +optional
	ExistingDatabaseName string `json:"existingBucketName,omitempty"`
}

type DatabaseRequestStatus struct {

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// DatabaseName is the name of the provisioned Database in response
	// to this DatabaseRequest.
	// +optional
	DatabaseName string `json:"databaseName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
type DatabaseRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseRequestSpec   `json:"spec,omitempty"`
	Status DatabaseRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseRequest `json:"items"`
}
