package controller

import (
	"github.com/saurabhbali/mongodbatlas-operator/pkg/controller/mongodbatlasdatabase"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mongodbatlasdatabase.Add)
}
