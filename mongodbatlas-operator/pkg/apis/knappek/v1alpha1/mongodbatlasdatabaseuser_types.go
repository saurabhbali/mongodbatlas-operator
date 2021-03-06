package v1alpha1

import (
	//ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Role allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type Role struct {
	DatabaseName   string `json:"databaseName,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
	RoleName       string `json:"roleName,omitempty"`
}

// Role allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type Scope struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// MongoDBAtlasDatabaseUserRequestBody defines the Request Body Parameters when creating/updating a database user
type MongoDBAtlasDatabaseUserRequestBody struct {
	Password        string    `json:"password,omitempty"`
	DeleteAfterDate string    `json:"deleteAfterDate,omitempty"`
	DatabaseName    string    `json:"databaseName,omitempty"`
	Roles           []Role    `json:"roles,omitempty"`
	Scopes          []Scope   `json:"scopes,omitempty"`
}

// MongoDBAtlasDatabaseUserSpec defines the desired state of MongoDBAtlasDatabaseUser
type MongoDBAtlasDatabaseUserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ProjectName                         string `json:"projectName,project"`
	MongoDBAtlasDatabaseUserRequestBody `json:",inline"`
}

// MongoDBAtlasDatabaseUserStatus defines the observed state of MongoDBAtlasDatabaseUser
type MongoDBAtlasDatabaseUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	GroupID         string    `json:"groupID,omitempty"`
	Username        string    `json:"username,omitempty"`
	DeleteAfterDate string    `json:"deleteAfterDate,omitempty"`
	DatabaseName    string    `json:"databaseName,omitempty"`
	Roles           []Role    `json:"roles,omitempty"`
        Scopes          []Scope   `json:"scopes,omitempty"`
}

// IsNotEqual does not purely compare equality like reflect.DeepEqual does
// It returns false if a and b are equal
// It returns false if a and b are not equal but a is nil
// It returns true if a and b are not equal and a is not nil
func IsNotEqual(a, b interface{}) bool {
	if a != b {
		return !IsZeroValue(a)
	}
	return false
}

// IsZeroValue returns true if input interface is the corresponding zero value
func IsZeroValue(i interface{}) bool {
	if i == nil {
		return true
	} // nil interface
	if i == "" {
		return true
	} // zero value of a string
	if i == 0.0 {
		return true
	} // zero value of a float64
	if i == 0 {
		return true
	} // zero value of an int
	if i == false {
		return true
	} // zero value of a boolean
	return false
}

// IsMongoDBAtlasDatabaseUserToBeUpdated is used to compare spec.MongoDBAtlasDatabaseUserRequestBody with status
func IsMongoDBAtlasDatabaseUserToBeUpdated(m1 MongoDBAtlasDatabaseUserRequestBody, m2 MongoDBAtlasDatabaseUserStatus) bool {
	if ok := IsNotEqual(m1.DeleteAfterDate, m2.DeleteAfterDate); ok {
		return true
	}
	for idx, role := range m1.Roles {
		if ok := IsNotEqual(role.DatabaseName, m2.Roles[idx].DatabaseName); ok {
			return true
		}
		if ok := IsNotEqual(role.CollectionName, m2.Roles[idx].CollectionName); ok {
			return true
		}
		if ok := IsNotEqual(role.RoleName, m2.Roles[idx].RoleName); ok {
			return true
		}
	}
	return false
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBAtlasDatabaseUser is the Schema for the mongodbatlasdatabaseusers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mongodbatlasdatabaseusers,scope=Namespaced
type MongoDBAtlasDatabaseUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBAtlasDatabaseUserSpec   `json:"spec,omitempty"`
	Status MongoDBAtlasDatabaseUserStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBAtlasDatabaseUserList contains a list of MongoDBAtlasDatabaseUser
type MongoDBAtlasDatabaseUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBAtlasDatabaseUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBAtlasDatabaseUser{}, &MongoDBAtlasDatabaseUserList{})
}
