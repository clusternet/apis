/*
Copyright 2021 The Clusternet Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Important: Run "make generated" to regenerate code after modifying this file

type ClusterType string

// These are the valid values for ClusterType
const (
	EdgeCluster ClusterType = "EdgeCluster"

	StandardCluster ClusterType = "StandardCluster"

	// todo: add more types
)

type ClusterSyncMode string

// These are the valid values for ClusterSyncMode
const (
	// Push means that all the resource changes in the parent cluster will be synchronized, pushed and applied to child clusters.
	Push ClusterSyncMode = "Push"

	// Pull means that the agent, known as 'clusternet-agent', running in the child cluster will watch, synchronize
	// and apply all the resource changes from the parent cluster to child cluster.
	Pull ClusterSyncMode = "Pull"

	// Dual combines both Push and Pull mode.
	Dual ClusterSyncMode = "Dual"
)

const (
	// ClusterReady means cluster is ready.
	ClusterReady = "Ready"
)

// ClusterRegistrationRequestSpec defines the desired state of ClusterRegistrationRequest
type ClusterRegistrationRequestSpec struct {
	// ClusterID, a Random (Version 4) UUID, is a unique value in time and space value representing for child cluster.
	// It is typically generated by the clusternet agent on the successful creation of a "self-cluster" Lease
	// in the child cluster.
	// Also it is not allowed to change on PUT operations.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	ClusterID types.UID `json:"clusterId"`

	// ClusterType denotes the type of the child cluster.
	//
	// +optional
	// +kubebuilder:validation:Type=string
	ClusterType ClusterType `json:"clusterType,omitempty"`

	// ClusterName is the cluster name.
	// a lower case alphanumeric characters or '-', and must start and end with an alphanumeric character
	//
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:Pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9]))*"
	ClusterName string `json:"clusterName,omitempty"`

	// ClusterNamespace is the dedicated namespace of the cluster.
	//
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	ClusterNamespace string `json:"clusterNamespace,omitempty"`

	// ClusterLabels is the labels of the child cluster.
	//
	// +optional
	// +kubebuilder:validation:Type=object
	ClusterLabels map[string]string `json:"clusterLabels,omitempty"`

	// SyncMode decides how to sync resources from parent cluster to child cluster.
	//
	// +optional
	// +kubebuilder:default=Pull
	// +kubebuilder:validation:Enum=Push;Pull;Dual
	SyncMode ClusterSyncMode `json:"syncMode,omitempty"`
}

// ClusterRegistrationRequestStatus defines the observed state of ClusterRegistrationRequest
type ClusterRegistrationRequestStatus struct {
	// DedicatedNamespace is a dedicated namespace for the child cluster, which is created in the parent cluster.
	//
	// +optional
	DedicatedNamespace string `json:"dedicatedNamespace,omitempty"`

	// DedicatedToken is populated by clusternet-hub when Result is RequestApproved.
	// With this token, the client could have full access on the resources created in DedicatedNamespace.
	//
	// +optional
	DedicatedToken []byte `json:"token,omitempty"`

	// CACertificate is the public certificate that is the root of trust for parent cluster
	// The certificate is encoded in PEM format.
	//
	// +optional
	CACertificate []byte `json:"caCertificate,omitempty"`

	// Result indicates whether this request has been approved.
	// When all necessary objects have been created and ready for child cluster registration,
	// this field will be set to "Approved". If any illegal updates on this object, "Illegal" will be set to this filed.
	//
	// +optional
	Result *ApprovedResult `json:"result,omitempty"`

	// ErrorMessage tells the reason why the request is not approved successfully.
	//
	// +optional
	ErrorMessage string `json:"errorMessage,omitempty"`

	// ManagedClusterName is the name of ManagedCluster object in the parent cluster corresponding to the child cluster
	//
	// +optional
	ManagedClusterName string `json:"managedClusterName,omitempty"`
}

type ApprovedResult string

// These are the possible results for a cluster registration request.
const (
	RequestDenied   ApprovedResult = "Denied"
	RequestApproved ApprovedResult = "Approved"
	RequestFailed   ApprovedResult = "Failed"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster",shortName=clsrr,categories=clusternet
// +kubebuilder:printcolumn:name="CLUSTER ID",type=string,JSONPath=`.spec.clusterId`,description="The unique id for the cluster"
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=`.status.result`,description="The status of current cluster registration request"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterRegistrationRequest is the Schema for the clusterregistrationrequests API
type ClusterRegistrationRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterRegistrationRequestSpec   `json:"spec,omitempty"`
	Status ClusterRegistrationRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRegistrationRequestList contains a list of ClusterRegistrationRequest
type ClusterRegistrationRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRegistrationRequest `json:"items"`
}

// ManagedClusterSpec defines the desired state of ManagedCluster
type ManagedClusterSpec struct {
	// ClusterID, a Random (Version 4) UUID, is a unique value in time and space value representing for child cluster.
	// It is typically generated by the clusternet agent on the successful creation of a "self-cluster" Lease
	// in the child cluster.
	// Also it is not allowed to change on PUT operations.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	ClusterID types.UID `json:"clusterId"`

	// ClusterType denotes the type of the child cluster.
	//
	// +optional
	// +kubebuilder:validation:Type=string
	ClusterType ClusterType `json:"clusterType,omitempty"`

	// SyncMode decides how to sync resources from parent cluster to child cluster.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Enum=Push;Pull;Dual
	SyncMode ClusterSyncMode `json:"syncMode"`

	// Taints has the "effect" on any resource that does not tolerate the Taint.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
}

// ManagedClusterStatus defines the observed state of ManagedCluster
type ManagedClusterStatus struct {
	// lastObservedTime is the time when last status from the series was seen before last heartbeat.
	// RFC 3339 date and time at which the object was acknowledged by the Clusternet Agent.
	// +optional
	LastObservedTime metav1.Time `json:"lastObservedTime,omitempty"`

	// k8sVersion is the Kubernetes version of the cluster
	// +optional
	KubernetesVersion string `json:"k8sVersion,omitempty"`

	// platform indicates the running platform of the cluster
	// +optional
	Platform string `json:"platform,omitempty"`

	// APIServerURL indicates the advertising url/address of managed Kubernetes cluster
	// +optional
	APIServerURL string `json:"apiserverURL,omitempty"`

	// Healthz indicates the healthz status of the cluster
	// which is deprecated since Kubernetes v1.16. Please use Livez and Readyz instead.
	// Leave it here only for compatibility.
	// +optional
	Healthz bool `json:"healthz"`

	// Livez indicates the livez status of the cluster
	// +optional
	Livez bool `json:"livez"`

	// Readyz indicates the readyz status of the cluster
	// +optional
	Readyz bool `json:"readyz"`

	// AppPusher indicates whether to allow parent cluster deploying applications in Push or Dual Mode.
	// Mainly for security concerns.
	// +optional
	AppPusher bool `json:"appPusher,omitempty"`

	// UseSocket indicates whether to use socket proxy when connecting to child cluster.
	//
	// +optional
	UseSocket bool `json:"useSocket,omitempty"`

	// Allocatable is the sum of allocatable resources for nodes in the cluster
	// +optional
	Allocatable corev1.ResourceList `json:"allocatable,omitempty"`

	// Capacity is the sum of capacity resources for nodes in the cluster
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// ClusterCIDR is the CIDR range of the cluster
	// +optional
	ClusterCIDR string `json:"clusterCIDR,omitempty"`

	// ServcieCIDR is the CIDR range of the services
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// NodeStatistics is the info summary of nodes in the cluster
	// +optional
	NodeStatistics NodeStatistics `json:"nodeStatistics,omitempty"`

	// PodStatistics is the info summary of pods in the cluster
	// +optional
	PodStatistics *PodStatistics `json:"podStatistics,omitempty"`

	// ResourceUsage is the cpu(m) and memory(Mi) already used in the cluster
	// +optional
	ResourceUsage *ResourceUsage `json:"resourceUsage,omitempty"`

	// Conditions is an array of current cluster conditions.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// heartbeatFrequencySeconds is the frequency at which the agent reports current cluster status
	// +optional
	HeartbeatFrequencySeconds *int64 `json:"heartbeatFrequencySeconds,omitempty"`

	// PredictorEnabled indicates whether predictor is enabled.
	// +optional
	PredictorEnabled bool `json:"predictorEnabled,omitempty"`

	// PredictorAddress shows the predictor address
	// +optional
	PredictorAddress string `json:"predictorAddress,omitempty"`

	// PredictorDirectAccess indicates whether the predictor can be accessed directly by clusternet-scheduler
	PredictorDirectAccess bool `json:"predictorDirectAccess,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=mcls,categories=clusternet
// +kubebuilder:printcolumn:name="CLUSTER ID",type=string,JSONPath=`.spec.clusterId`,description="The unique id for the cluster"
// +kubebuilder:printcolumn:name="CLUSTER TYPE",type=string,JSONPath=`.spec.clusterType`,description="The type of the cluster",priority=100
// +kubebuilder:printcolumn:name="SYNC MODE",type=string,JSONPath=`.spec.syncMode`,description="The cluster sync mode"
// +kubebuilder:printcolumn:name="KUBERNETES",type=string,JSONPath=".status.k8sVersion"
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ManagedCluster is the Schema for the managedclusters API
type ManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterSpec   `json:"spec,omitempty"`
	Status ManagedClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterList contains a list of ManagedCluster
type ManagedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedCluster `json:"items"`
}

type NodeStatistics struct {
	// ReadyNodes is the number of ready nodes in the cluster
	// +optional
	ReadyNodes int32 `json:"readyNodes,omitempty"`

	// NotReadyNodes is the number of not ready nodes in the cluster
	// +optional
	NotReadyNodes int32 `json:"notReadyNodes,omitempty"`

	// UnknownNodes is the number of unknown nodes in the cluster
	// +optional
	UnknownNodes int32 `json:"unknownNodes,omitempty"`

	// LostNodes is the number of states lost nodes in the cluster
	// +optional
	LostNodes int32 `json:"lostNodes,omitempty"`
}

type PodStatistics struct {
	// RunningPods is the number of running pods in the cluster
	// +optional
	RunningPods int32 `json:"runningPods,omitempty"`

	// TotalPods is the number of all pods in the cluster
	// +optional
	TotalPods int32 `json:"totalPods,omitempty"`
}

type ResourceUsage struct {
	// CpuUsage is the total cpu(m) already used in the whole cluster, k8s reserved not include
	// +optional
	CpuUsage resource.Quantity `json:"cpuUsage,omitempty"`

	// MemoryUsage is the total memory(Mi) already used in the whole cluster, k8s reserved not include
	// +optional
	MemoryUsage resource.Quantity `json:"memoryUsage,omitempty"`
}
