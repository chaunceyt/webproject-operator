package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebProjectSpec defines the desired state of WebProject
type WebProjectSpec struct {
	ProjectDomainName       string `json:"projectdomainname"`
	ReleaseName             string `json:"releasename"`
	ProjectName             string `json:"projectname"`
	WebImage                string `json:"webimage"`
	CLIImage                string `json:"cliimage"`
	CacheImage              string `json:"cacheimage"`
	CachePort               int32  `json:"cacheport"`
	DatabaseImage           string `json:"databaseimage"`
	FileStorageSize         string `json:"filestoragesize"`
	FileStorageMountPath    string `json:"filestoragemountpath"`
	DatabaseName            string `json:"databasename"`
	DatabaseUser            string `json:"databaseuser"`
	DatabaseUserPassword    string `json:"databaseuserpassword"`
	DatabaseStorageSize     string `json:"databasestoragesize"`
	DatabaseStoreMountPath  string `json:"databasestoragemountpath"`
	DatabaseRootPassword    string `json:"databaserootpassword"`
	DockerConfigUsername    string `json:"dockerconfigusername,omitempty"`
	DockerConfigPassword    string `json:"dockerconfiguserpassword,omitempty"`
	DockerConfigEmail       string `json:"dockerconfiguseremail,omitempty"`
	DockerConfigRegistryURL string `json:"dockerconfigregistryurl,omitempty"`
	InitContainerScript     string `json:"initcontainerscript,omitempty"`
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
