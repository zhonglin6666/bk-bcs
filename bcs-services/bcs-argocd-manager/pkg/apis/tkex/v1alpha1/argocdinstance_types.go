/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArgocdInstanceSpec defines the desired state of ArgocdInstance
type ArgocdInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ArgocdInstance. Edit argocdinstance_types.go to remove/update
	Foo string `json:"foo,omitempty" protobuf:"bytes,1,opt,name=foo"`
}

// ArgocdInstanceStatus defines the observed state of ArgocdInstance
type ArgocdInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ArgocdInstance is the Schema for the argocdinstances API
type ArgocdInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ArgocdInstanceSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ArgocdInstanceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

//+kubebuilder:object:root=true

// ArgocdInstanceList contains a list of ArgocdInstance
type ArgocdInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ArgocdInstance `json:"items" protobuf:"bytes,2,rep,name=items"`
}
