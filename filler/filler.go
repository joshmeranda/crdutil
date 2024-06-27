package filler

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Filler interface {
	Fill(*apiextensionsv1.CustomResourceDefinitionVersion) (map[string]any, error)
}
