package filler

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Default struct {
}

func (d Default) propValue(prop apiextensionsv1.JSONSchemaProps) (any, error) {
	if prop.Default != nil {
		return prop.Default, nil
	}

	switch prop.Type {
	case "object":
		value, err := d.parseProperties(prop.Properties)
		if err != nil {
			return nil, err
		}

		return value, nil
	case "string":
		return "", nil
	case "integer":
		return 0, nil
	case "boolean":
		return false, nil
	case "array":
		return make([]any, 0), nil
	default:
		return nil, fmt.Errorf("unrecognized type '%s", prop.Type)
	}
}

func (d Default) parseProperties(props map[string]apiextensionsv1.JSONSchemaProps) (map[string]any, error) {
	data := make(map[string]any, len(props))

	for key, prop := range props {
		value, err := d.propValue(prop)
		if err != nil {
			return nil, fmt.Errorf("failed to determine default value: %w", err)
		}

		fmt.Printf("=== [Default parseProperties] 000 '%s' : '%+v'\n", key, value)
		data[key] = value
	}

	return data, nil
}

func (d Default) Fill(crdVersion *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]any, error) {
	data, err := d.parseProperties(crdVersion.Schema.OpenAPIV3Schema.Properties)
	if err != nil {
		return nil, fmt.Errorf("could not parse properties: %w", err)
	}

	return data, nil
}
