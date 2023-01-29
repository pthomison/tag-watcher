/*
Copyright 2023.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TagReflectorSpec defines the desired state of TagReflector
type TagReflectorSpec struct {
	Registry string `json:"registry,omitempty"`
	Regex    string `json:"regex,omitempty"`
}

// TagReflectorStatus defines the observed state of TagReflector
type TagReflectorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TagReflector is the Schema for the tagreflectors API
type TagReflector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TagReflectorSpec   `json:"spec,omitempty"`
	Status TagReflectorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TagReflectorList contains a list of TagReflector
type TagReflectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TagReflector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TagReflector{}, &TagReflectorList{})
}
