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
	SchemeBuilder.Register(&DatabaseClass{}, &DatabaseClassList{})
	SchemeBuilder.Register(&DatabaseAccessClass{}, &DatabaseAccessClassList{})
}

type AuthenticationType string

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
type DatabaseClass struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// DriverName is the name of driver associated with this database
	DriverName string `json:"driverName"`

	// Parameters is an opaque map for passing in configuration to a driver
	// for creating the bucket
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseClass `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
type DatabaseAccessClass struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// DriverName is the name of driver associated with
	// this DatabaseAccess
	DriverName string `json:"driverName"`

	// AuthenticationType denotes the style of authentication
	AuthenticationType AuthenticationType `json:"authenticationType"`

	// Parameters is an opaque map for passing in configuration to a driver
	// for granting access to a bucket
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseAccessClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseAccessClass `json:"items"`
}
