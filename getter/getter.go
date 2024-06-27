package getter

import (
	"context"
	"fmt"
	"os"

	// "gopkg.in/yaml.v3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type Getter func(context.Context) (*apiextensionsv1.CustomResourceDefinition, error)

func FromFile(path string) Getter {
	return func(_ context.Context) (*apiextensionsv1.CustomResourceDefinition, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed reading crd at '%s': %w", path, err)
		}

		crd := &apiextensionsv1.CustomResourceDefinition{}

		if err := yaml.Unmarshal(data, crd); err != nil {
			return nil, fmt.Errorf("failde to unmarshal crd: %w", err)
		}

		return crd, nil
	}
}
