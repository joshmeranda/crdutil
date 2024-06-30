package filler

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Filler interface {
	Fill(schema.GroupVersionKind, *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]any, error)
}
