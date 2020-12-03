/*
Copyright 2019 Bloomberg Finance LP.

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
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// StatefulSetOptions defines custom options for StatefulSets
type StatefulSetOptions struct {
	// Annotations to be added for the StatefulSet.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for the StatefulSet.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// PodManagementPolicy defines the policy for creating pods under a stateful set.
	// Override the default value of Parallel.
	// This cannot be updated on an existing StatefulSet, the StatefulSet must be deleted and recreated for a change in this field to take effect.
	//
	// +optional
	// +kubebuilder:validation:Enum=OrderedReady;Parallel
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
}

// DeploymentOptions defines custom options for Deployments
type DeploymentOptions struct {
	// Annotations to be added for the Deployment.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for the Deployment.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// PodOptions defines the common pod configuration for Pods, including when used
// in deployments, stateful-sets, etc.
type PodOptions struct {
	// The scheduling constraints on pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Resources is the resource requirements for the container.
	// This field cannot be updated once the cluster is created.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Additional non-data volumes to load into the default container.
	// +optional
	Volumes []AdditionalVolume `json:"volumes,omitempty"`

	// PodSecurityContext is the security context for the pod.
	// +optional
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// Additional environment variables to pass to the default container.
	// +optional
	EnvVariables []corev1.EnvVar `json:"envVars,omitempty"`

	// Annotations to be added for pods.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for pods.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Tolerations to be added for the StatefulSet.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Node Selector to be added for the StatefulSet.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Liveness probe parameters
	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// Readiness probe parameters
	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`

	// Startup probe parameters
	// +optional
	StartupProbe *corev1.Probe `json:"startupProbe,omitempty"`

	// PriorityClassName for the pod
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Sidecar containers to run in the pod. These are in addition to the Solr Container
	// +optional
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`

	// Additional init containers to run in the pod.
	// These will run along with the init container that sets up the "solr.xml".
	// +optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
}

// ServiceOptions defines custom options for services
type ServiceOptions struct {
	// Annotations to be added for the Service.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for the Service.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// IngressOptions defines custom options for ingresses
type IngressOptions struct {
	// Annotations to be added for the Ingress.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for the Ingress.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// ConfigMapOptions defines custom options for configMaps
type ConfigMapOptions struct {
	// Annotations to be added for the ConfigMap.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added for the ConfigMap.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Name of a user provided ConfigMap in the same namespace containing a custom solr.xml
	// +optional
	ProvidedConfigMap string `json:"providedConfigMap,omitempty"`
}

// AdditionalVolume provides information on additional volumes that should be loaded into pods
type AdditionalVolume struct {
	// Name of the volume
	Name string `json:"name"`

	// Source is the source of the Volume to be loaded into the solrCloud Pod
	Source corev1.VolumeSource `json:"source"`

	// DefaultContainerMount defines how to mount this volume into the default container.
	// If this volume is to be used only with sidecar or non-default init containers,
	// then this option is not necessary.
	// +optional
	DefaultContainerMount *corev1.VolumeMount `json:"defaultContainerMount,omitempty"`
}

// ContainerImage defines the fields needed for a Docker repository image. The
// format here matches the predominant format used in Helm charts.
type ContainerImage struct {
	// +optional
	Repository string `json:"repository,omitempty"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
	// +optional
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
}

func (c *ContainerImage) withDefaults(repo string, version string, policy corev1.PullPolicy) (changed bool) {
	if c.Repository == "" && repo != "" {
		changed = true
		c.Repository = repo
	}
	if c.Tag == "" && version != "" {
		changed = true
		c.Tag = version
	}
	if c.PullPolicy == "" && policy != "" {
		changed = true
		c.PullPolicy = policy
	}
	return changed
}

func (c *ContainerImage) ToImageName() (name string) {
	return c.Repository + ":" + c.Tag
}

func ImageVersion(image string) (version string) {
	split := strings.Split(image, ":")
	if len(split) < 2 {
		return ""
	} else {
		return split[1]
	}
}

// SolrNodeStatus is the status of a solrNode in the cloud, with readiness status
// and internal and external addresses
type ZookeeperConnectionInfo struct {
	// The connection string to connect to the ensemble from within the Kubernetes cluster
	// +optional
	InternalConnectionString string `json:"internalConnectionString,omitempty"`

	// The connection string to connect to the ensemble from outside of the Kubernetes cluster
	// If external and no internal connection string is provided, the external cnx string will be used as the internal cnx string
	// +optional
	ExternalConnectionString *string `json:"externalConnectionString,omitempty"`

	// The ChRoot to connect solr at
	// +optional
	ChRoot string `json:"chroot,omitempty"`

	// ZooKeeper ACL to use when connecting with ZK.
	// This ACL should have ALL permission in the given chRoot.
	// +optional
	AllACL *ZookeeperACL `json:"acl,omitempty"`

	// ZooKeeper ACL to use when connecting with ZK for reading operations.
	// This ACL should have READ permission in the given chRoot.
	// +optional
	ReadOnlyACL *ZookeeperACL `json:"readOnlyAcl,omitempty"`
}

func (ci *ZookeeperConnectionInfo) withDefaults() (changed bool) {
	if ci.InternalConnectionString == "" {
		if ci.ExternalConnectionString != nil {
			changed = true
			ci.InternalConnectionString = *ci.ExternalConnectionString
		}
	}
	if ci.ChRoot == "" {
		changed = true
		ci.ChRoot = "/"
	} else if !strings.HasPrefix(ci.ChRoot, "/") {
		changed = true
		ci.ChRoot = "/" + ci.ChRoot
	}
	return changed
}

// ZookeeperSpec defines the internal zookeeper ensemble to run for solr
type ZookeeperACL struct {
	// The name of the Kubernetes Secret that stores the username and password for the ACL.
	// This secret must be in the same namespace as the solrCloud or prometheusExporter is running in.
	SecretRef string `json:"secret"`

	// The name of the key in the given secret that contains the ACL username
	UsernameKey string `json:"usernameKey"`

	// The name of the key in the given secret that contains the ACL password
	PasswordKey string `json:"passwordKey"`
}
