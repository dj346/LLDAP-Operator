package internal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group   = "lldap-operator.github.com"
	Version = "v1alpha1"
)

// GroupVersion is group version used to register these objects.
var GroupVersion = schema.GroupVersion{Group: Group, Version: Version}

// SchemeBuilder is used to add go types to the GroupVersionKind scheme.
var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

// AddToScheme adds the types in this group-version to the given scheme.
var AddToScheme = SchemeBuilder.AddToScheme

// LLDAPUserSpec defines your spec fields.
type LLDAPUserSpec struct {
	Username         string   `json:"username"`
	DisplayName      string   `json:"displayName,omitempty"`
	PrimaryEmail     string   `json:"primaryEmail,omitempty"`
	AdditionalEmails []string `json:"additionalEmails,omitempty"`
	Groups           []string `json:"groups,omitempty"`
}

// LLDAPUserStatus is optional; keep minimal for now.
type LLDAPUserStatus struct {
	Synced bool `json:"synced,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type LLDAPUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LLDAPUserSpec   `json:"spec,omitempty"`
	Status LLDAPUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type LLDAPUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LLDAPUser `json:"items"`
}

// ---- Group ----

type LLDAPGroupSpec struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	GIDNumber   *int64   `json:"gidNumber,omitempty"`
	Members     []string `json:"members,omitempty"`
}

type LLDAPGroupStatus struct {
	Synced      bool  `json:"synced,omitempty"`
	MemberCount int32 `json:"memberCount,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type LLDAPGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LLDAPGroupSpec   `json:"spec,omitempty"`
	Status LLDAPGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type LLDAPGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LLDAPGroup `json:"items"`
}

// Register the types with the scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&LLDAPUser{},
		&LLDAPUserList{},
		&LLDAPGroup{},
		&LLDAPGroupList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
