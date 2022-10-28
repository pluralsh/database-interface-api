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

import (
	crhelperTypes "github.com/pluralsh/controller-reconcile-helper/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}

const (
	// DatabaseReadyCondition used when database is ready.
	DatabaseReadyCondition crhelperTypes.ConditionType = "DatabaseReady"

	// FailedToCreateDatabaseReason used when grpc method for database creation failed.
	FailedToCreateDatabaseReason = "FailedToCreateDatabase"
)

type DatabaseSpec struct {
	// DriverName is the name of driver associated with this database
	DriverName string `json:"driverName"`

	// Name of the DatabaseClass specified in the DatabaseRequest
	DatabaseClassName string `json:"databaseClassName"`

	// Name of the DatabaseRequest that resulted in the creation of this Database
	// In case the Database object was created manually, then this should refer
	// to the DatabaseRequest with which this Database should be bound
	DatabaseRequest *corev1.ObjectReference `json:"databaseRequest"`

	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`

	// ExistingDatabaseID is the unique id of the database.
	// This field will be empty when the Database is dynamically provisioned by operator.
	// +optional
	ExistingDatabaseID string `json:"existingBucketID,omitempty"`

	// DeletionPolicy is used to specify how to handle deletion. There are 2 possible values:
	//  - Retain: Indicates that the database should not be deleted (default)
	//  - Delete: Indicates that the database should be deleted
	//
	// +optional
	// +kubebuilder:default:=Retain
	DeletionPolicy DeletionPolicy `json:"deletionPolicy"`
}

type DatabaseStatus struct {
	// Ready is a boolean condition to reflect the successful creation
	// of a database.
	Ready bool `json:"ready,omitempty"`

	// DatabaseID is the unique id of the database
	// +optional
	DatabaseID string `json:"databaseID,omitempty"`

	// Conditions defines current state.
	// +optional
	Conditions crhelperTypes.Conditions `json:"conditions,omitempty"`
}

// GetConditions returns the list of conditions for a WireGuardServer API object.
func (db *Database) GetConditions() crhelperTypes.Conditions {
	return db.Status.Conditions
}

// SetConditions will set the given conditions on a WireGuardServer object.
func (db *Database) SetConditions(conditions crhelperTypes.Conditions) {
	db.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Database ready status"

type Database struct {
	metav1.TypeMeta `json:",inline"`
	// +optional

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DatabaseSpec `json:"spec,omitempty"`

	// +optional
	Status DatabaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}
