/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KratosSpec defines the desired state of Kratos.
type KratosSpec struct {
	// ProjectID specifies the Google Cloud project in which to create the cluster
	ProjectID string `json:"projectID"`

	// ClusterName is the name of the GKE cluster to be created
	ClusterName string `json:"clusterName"`

	// Region is the geographical region where the cluster will be deployed
	Region string `json:"region"`

	// NodePools defines the node pool configurations
	NodePools []NodePoolSpec `json:"nodePools"`

	// Networking defines network-related settings
	Networking NetworkingSpec `json:"networking"`
}

// NodePoolSpec defines the node pool configuration
type NodePoolSpec struct {
	Name        string `json:"name"`
	MachineType string `json:"machineType"`
	NodeCount   int    `json:"nodeCount"`
	AutoScaling bool   `json:"autoScaling"`
	MinNodes    int    `json:"minNodes,omitempty"`
	MaxNodes    int    `json:"maxNodes,omitempty"`
	DiskSizeGB  int    `json:"diskSizeGB,omitempty"`
	Preemptible bool   `json:"preemptible,omitempty"`
}

// NetworkingSpec defines the network configurations for the cluster
type NetworkingSpec struct {
	VPCName       string `json:"vpcName"`
	SubnetName    string `json:"subnetName"`
	EnableIPAlias bool   `json:"enableIPAlias"`
	PodCIDR       string `json:"podCIDR,omitempty"`
	ServiceCIDR   string `json:"serviceCIDR,omitempty"`
}

// KratosStatus defines the observed state of Kratos.
type KratosStatus struct {
	// Phase represents the current phase of cluster provisioning
	Phase ClusterPhase `json:"phase"`

	// Conditions provide detailed information about the cluster status
	Conditions []Condition `json:"conditions,omitempty"`

	// ClusterEndpoint is the API server endpoint of the created cluster
	ClusterEndpoint string `json:"clusterEndpoint,omitempty"`

	// NodePoolsStatus gives status information of each node pool
	NodePoolsStatus []NodePoolStatus `json:"nodePoolsStatus,omitempty"`

	// ErrorMessage captures any errors encountered during provisioning
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Condition represents a specific condition of the cluster
type Condition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// NodePoolStatus represents the observed state of a node pool
type NodePoolStatus struct {
	Name   string `json:"name"`
	Ready  int    `json:"ready"`
	Total  int    `json:"total"`
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Kratos is the Schema for the kratos API.
type Kratos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KratosSpec   `json:"spec,omitempty"`
	Status KratosStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KratosList contains a list of Kratos.
type KratosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kratos `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kratos{}, &KratosList{})
}
