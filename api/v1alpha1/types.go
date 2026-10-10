// +kubebuilder:object:generate=true
// +groupName=example.org
package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// let these kubebuilder validation lives in the crossplane controlplane

type XBucketsSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=us-east-1;us-east-2
	Region string `json:"region"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Names []string `json:"names"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,path=xbuckets
// +kubebuilder:storageversion
type XBuckets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec XBucketsSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type XBucketsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XBuckets `json:"items"`
}
