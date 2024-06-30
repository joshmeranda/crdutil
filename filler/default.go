package filler

import (
	"fmt"
	"reflect"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func setIfAbsentZeroedOrEmpty(data map[string]any, key string, newValue any) {
	value, found := data[key]
	if !found {
		data[key] = newValue
	}

	oldType := reflect.TypeOf(value)

	if !reflect.ValueOf(value).IsZero() {
		if oldType.Kind() == reflect.Map && len(value.(map[string]any)) > 0 {
			fmt.Printf("=== [setIfAbsentZeroedOrEmpty] 000 ===\n")
			return
		}
	}

	data[key] = newValue
}

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
			return nil, fmt.Errorf("failed to parse object properties: %w", err)
		}

		return value, nil
	case "string":
		return "", nil
	case "integer":
		return 0, nil
	case "number":
		return 0.0, nil
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

		data[key] = value
	}

	return data, nil
}

func (d Default) Fill(gvk schema.GroupVersionKind, crdVersion *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]any, error) {
	data, err := d.parseProperties(crdVersion.Schema.OpenAPIV3Schema.Properties)
	if err != nil {
		return nil, fmt.Errorf("could not parse properties: %w", err)
	}

	setIfAbsentZeroedOrEmpty(data, "apiVersion", gvk.GroupVersion().String())
	setIfAbsentZeroedOrEmpty(data, "kind", gvk.Kind)
	setIfAbsentZeroedOrEmpty(data, "metadata", map[string]any{"name": strings.ToLower(gvk.Kind) + "-example"})

	delete(data, "status")

	return data, nil
}
