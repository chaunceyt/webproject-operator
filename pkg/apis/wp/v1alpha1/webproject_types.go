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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebProjectSpec defines the desired state of WebProject
type WebProjectSpec struct {
	AWSSecretName         string                    `json:"awssecretname,omitempty"`
	CLISidecar            WebProjectCLISidecar      `json:"clisidecar,omitempty"`
	CacheSidecar          WebProjectCacheSidecar    `json:"cachesidecar,omitempty"`
	DatabaseSidecar       WebProjectDatabaseSidecar `json:"databasesidecar,omitempty"`
	DockerConfig          WebProjectDockerConfig    `json:"dockerconfig,omitempty"`
	FileStorageSize       string                    `json:"filestoragesize"`
	FileStorageMountPath  string                    `json:"filestoragemountpath"`
	InitContainerScript   string                    `json:"initcontainerscript,omitempty"`
	IngressHosts          []string                  `json:"ingresshosts,omitempty"`
	DeploymentAnnotations map[string]string         `json:"deploymentannotations,omitempty"`
	IngressAnnotations    map[string]string         `json:"ingressannotations,omitempty"`
	ProjectDomainName     string                    `json:"projectdomainname"`
	ProjectName           string                    `json:"projectname"`
	ReleaseName           string                    `json:"releasename"`
	SearchSidecar         WebProjectSearchSidecar   `json:"searchsidecar,omitempty"`
	WebContainer          WebProjectWebContainer    `json:"webcontainer"`
	CommonConfig          map[string]string         `json:"commonconfig"`
}

// WebProjectWebContainer defines to spec for webcontainer
type WebProjectWebContainer struct {
	Image string `json:"image"`
}

// WebProjectDockerConfig defines setup ImagePullSecret for private registry.
type WebProjectDockerConfig struct {
	Enabled    bool   `json:"enabled"`
	Secretname string `json:"secretname,omitempty"`
}

// WebProjectCacheSidecar defines spec for cache sidecar
type WebProjectCacheSidecar struct {
	Enabled bool   `json:"enabled"`
	Image   string `json:"image,omitempty"`
	Port    int32  `json:"port,omitempty"`
}

// WebProjectCLISidecar defines spec for cache sidecar
type WebProjectCLISidecar struct {
	Enabled bool   `json:"enabled"`
	Image   string `json:"image,omitempty"`
	Port    int32  `json:"port,omitempty"`
}

// WebProjectCacheSidecar defines spec for cache sidecar
type WebProjectSearchSidecar struct {
	Enabled bool   `json:"enabled"`
	Engine  string `json:"engine,omitempty"`
	Image   string `json:"image,omitempty"`
}

// WebProjectDatabaseSidecar defines the desired state for database sidecar
type WebProjectDatabaseSidecar struct {
	Enabled                bool   `json:"enabled,omitempty"`
	DatabaseName           string `json:"databasename,omitempty"`
	DatabaseImage          string `json:"databaseimage,omitempty"`
	DatabaseUser           string `json:"databaseuser,omitempty"`
	DatabaseStorageSize    string `json:"databasestoragesize,omitempty"`
	DatabaseRootPassword   string `json:"databaserootpassword,omitempty"`
	DatabaseStoreMountPath string `json:"databasestoragemountpath,omitempty"`
	DatabaseUserPassword   string `json:"databaseuserpassword,omitempty"`
}

type WebProjectBackup struct {
	Enabled         bool   `json:"enabled,omitempty"`
	StorageProvider string `json:"storageprovider,omitempty"`
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
