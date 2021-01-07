/*
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

// WebProjectSpec defines the desired state of WebProject
type WebProjectSpec struct {
	AWSSecretName           string                  `json:"awssecretname,omitempty"`
	CacheImage              string                  `json:"cacheimage,omitempty"`
	CachePort               int32                   `json:"cacheport,omitempty"`
	CLIImage                string                  `json:"cliimage,omitempty"`
	DatabaseName            string                  `json:"databasename,omitempty"`
	DatabaseImage           string                  `json:"databaseimage,omitempty"`
	DatabaseUser            string                  `json:"databaseuser,omitempty"`
	DatabaseStorageSize     string                  `json:"databasestoragesize,omitempty"`
	DatabaseRootPassword    string                  `json:"databaserootpassword,omitempty"`
	DatabaseStoreMountPath  string                  `json:"databasestoragemountpath,omitempty"`
	DatabaseUserPassword    string                  `json:"databaseuserpassword,omitempty"`
	DockerConfigEmail       string                  `json:"dockerconfiguseremail,omitempty"`
	DockerConfigPassword    string                  `json:"dockerconfiguserpassword,omitempty"`
	DockerConfigRegistryURL string                  `json:"dockerconfigregistryurl,omitempty"`
	DockerConfigUsername    string                  `json:"dockerconfigusername,omitempty"`
	FileStorageSize         string                  `json:"filestoragesize"`
	FileStorageMountPath    string                  `json:"filestoragemountpath"`
	InitContainerScript     string                  `json:"initcontainerscript,omitempty"`
	IngressHost             []WebProjectIngressHost `json:"ingresshosts,omitempty"`
	ProjectDomainName       string                  `json:"projectdomainname"`
	ProjectName             string                  `json:"projectname"`
	ReleaseName             string                  `json:"releasename"`
	WebImage                string                  `json:"webimage"`
}

type WebProjectIngressHost struct {
	Hostname string `json:"ingresshost,omitempty"`
}

// WebProjectStatus defines the observed state of WebProject
type WebProjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WebProject is the Schema for the webprojects API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=webprojects,scope=Namespaced
type WebProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebProjectSpec   `json:"spec,omitempty"`
	Status WebProjectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WebProjectList contains a list of WebProject
type WebProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebProject{}, &WebProjectList{})
}
