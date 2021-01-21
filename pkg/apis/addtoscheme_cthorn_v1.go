package apis

import (
	v "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v.SchemeBuilder.AddToScheme)
}
