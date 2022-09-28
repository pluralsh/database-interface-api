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
	SchemeBuilder.Register(&DatabaseAccess{}, &DatabaseAccessList{})
}

type DatabaseAccessSpec struct {
	// DatabaseRequestName is the name of the DatabaseRequest.
	DatabaseRequestName string `json:"databaseRequestName"`

	// DatabaseAccessClassName is the name of the DatabaseAccessClass
	DatabaseAccessClassName string `json:"bucketAccessClassName"`

	// CredentialsSecretName is the name of the secret that operator should populate
	// with the credentials. If a secret by this name already exists, then it is
	// assumed that credentials have already been generated. It is not overridden.
	// This secret is deleted when the DatabaseAccess is delted.
	CredentialsSecretName string `json:"credentialsSecretName"`
}

type DatabaseAccessStatus struct {
	// AccountID is the unique ID for the account in the OSP. It will be populated
	// by the COSI sidecar once access has been successfully granted.
	// +optional
	AccountID string `json:"accountID,omitempty"`

	// AccessGranted indicates the successful grant of privileges to access the bucket
	// +optional
	AccessGranted bool `json:"accessGranted"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
type DatabaseAccess struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DatabaseAccessSpec `json:"spec,omitempty"`

	// +optional
	Status DatabaseAccessStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseAccessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseAccess `json:"items"`
}
