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
	AWSSecretName         string                    `json:"awsSecretName,omitempty"`
	Backup                WebProjectBackup          `json:"backup,omitempty"`
	CLISidecar            WebProjectCLISidecar      `json:"cliSidecar,omitempty"`
	CacheSidecar          WebProjectCacheSidecar    `json:"cacheSidecar,omitempty"`
	DatabaseSidecar       WebProjectDatabaseSidecar `json:"databaseSidecar,omitempty"`
	ImagePullSecrets      WebProjectImagePullSecret `json:"imagePullSecrets,omitempty"`
	FileStorageSize       string                    `json:"fileStorageSize"`
	FileStorageMountPath  string                    `json:"fileStorageMountPath"`
	InitContainerScript   string                    `json:"initContainerScript,omitempty"`
	IngressHosts          []string                  `json:"ingressHosts,omitempty"`
	DeploymentAnnotations map[string]string         `json:"deploymentAnnotations,omitempty"`
	IngressAnnotations    map[string]string         `json:"ingressAnnotations,omitempty"`
	ProjectDomainName     string                    `json:"projectDomainName"`
	ProjectName           string                    `json:"projectName"`
	ReleaseName           string                    `json:"releaseName"`
	SearchSidecar         WebProjectSearchSidecar   `json:"searchSidecar,omitempty"`
	WebContainer          WebProjectWebContainer    `json:"webContainer"`
	CommonConfig          map[string]string         `json:"commonConfig"`
}

// WebProjectWebContainer defines to spec for webcontainer
type WebProjectWebContainer struct {
	Image   string            `json:"image"`
	CronJob WebProjectCronJob `json:"cronJob,omitempty"`
}

// WebProjectDockerConfig defines setup ImagePullSecret for private registry.
type WebProjectImagePullSecret struct {
	Enabled    bool   `json:"enabled"`
	Secretname string `json:"secretName,omitempty"`
}

// WebProjectCronJob defines setup for cronjobs.
type WebProjectCronJob struct {
	Enabled  bool   `json:"enabled"`
	Schedule string `json:"schedule"`
	Script   string `json:"script"`
}

// WebProjectCacheSidecar defines spec for cache sidecar
type WebProjectCacheSidecar struct {
	CronJob       WebProjectCronJob `json:"cronJob,omitempty"`
	Enabled       bool              `json:"enabled"`
	Engine        string            `json:"engine,omitempty"`
	RedisPassword string            `json:"redisPassword,omitempty"`
}

// WebProjectCLISidecar defines spec for cache sidecar
type WebProjectCLISidecar struct {
	Enabled bool   `json:"enabled"`
	Image   string `json:"image,omitempty"`
	Port    int32  `json:"port,omitempty"`
}

// WebProjectSearchSidecar defines spec for cache sidecar
type WebProjectSearchSidecar struct {
	CronJob WebProjectCronJob `json:"cronJob,omitempty"`
	Enabled bool              `json:"enabled"`
	Engine  string            `json:"engine,omitempty"`
	Image   string            `json:"image,omitempty"`
}

// WebProjectDatabaseSidecar defines the desired state for database sidecar
type WebProjectDatabaseSidecar struct {
	Backup                 WebProjectBackup  `json:"backup,omitempty"`
	CronJob                WebProjectCronJob `json:"cronJob,omitempty"`
	DatabaseName           string            `json:"databaseName,omitempty"`
	DatabaseImage          string            `json:"databaseImage,omitempty"`
	DatabaseUser           string            `json:"databaseUser,omitempty"`
	DatabaseStorageSize    string            `json:"databaseStorageSize,omitempty"`
	DatabaseRootPassword   string            `json:"databaseRootPassword,omitempty"`
	DatabaseStoreMountPath string            `json:"databaseStorageMountPath,omitempty"`
	DatabaseUserPassword   string            `json:"databaseUserPassword,omitempty"`
	Enabled                bool              `json:"enabled,omitempty"`
}

// WebProjectBackup defined the spec for backups.
type WebProjectBackup struct {
	BackupSchedule                  string `json:"backupSchedule"`
	BackupScheduledJobsHistorylimit int    `json:"backupScheduledJobsHistorylimit,omitempty"`
	Enabled                         bool   `json:"enabled,omitempty"`
	StorageProvider                 string `json:"storageProvider,omitempty"`
}

// WebProjectStatus defines the observed state of WebProject
type WebProjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// WebProject WebProjectSpec `json:"webproject,omitempty"`
	PodNames        []string `json:"podNames"`
	ConfigMapNames  []string `json:"configMapNames"`
	SecretNames     []string `json:"secretNames"`
	DeploymentNames []string `json:"deploymentNames"`
	CronJobNames    []string `json:"cronJobNames"`
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
