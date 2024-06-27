package getter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type Getter func(context.Context) (*apiextensionsv1.CustomResourceDefinition, error)

func FromFile(path string) Getter {
	return func(_ context.Context) (*apiextensionsv1.CustomResourceDefinition, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not reading crd at '%s': %w", path, err)
		}

		crd := &apiextensionsv1.CustomResourceDefinition{}

		if err := yaml.Unmarshal(data, crd); err != nil {
			return nil, fmt.Errorf("could not unmarshal crd: %w", err)
		}

		return crd, nil
	}
}

func FromUrl(url string) Getter {
	return func(ctx context.Context) (*apiextensionsv1.CustomResourceDefinition, error) {
		data := make([]byte, 0)
		out := bytes.NewBuffer(data)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("could not fetch crd: %w", err)
		}
		defer resp.Body.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return nil, fmt.Errorf("could not prepare response body: %w", err)
		}

		crd := &apiextensionsv1.CustomResourceDefinition{}

		if err := yaml.Unmarshal(data, crd); err != nil {
			return nil, fmt.Errorf("could not unmarshal crd: %w", err)
		}

		return nil, err
	}
}

func FromCluster(client clientset.Clientset, name string) Getter {
	return func(ctx context.Context) (*apiextensionsv1.CustomResourceDefinition, error) {
		crd, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get crd from cluster: %w", err)
		}

		return crd, nil
	}
}
