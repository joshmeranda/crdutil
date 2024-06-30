package filler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	EditorEnvVar = "EDITOR"
)

type Editor struct {
	EditorPath string
}

func DefaultEditor() (Filler, error) {
	editor := os.Getenv(EditorEnvVar)
	if editor == "" {
		return nil, fmt.Errorf("could not determine default editor")
	}

	return &Editor{
		EditorPath: editor,
	}, nil
}

func (e Editor) writeDefaultData(gvk schema.GroupVersionKind, crdVersion *apiextensionsv1.CustomResourceDefinitionVersion) (string, error) {
	filled, err := Default{}.Fill(gvk, crdVersion)
	if err != nil {
		return "", fmt.Errorf("could not marshal data: %w", err)
	}

	data, err := yaml.Marshal(filled)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default crd: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "crdutil-*.yaml")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		return "", fmt.Errorf("could not write data to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}

func (e Editor) Fill(gvk schema.GroupVersionKind, crdVersion *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]any, error) {
	tmpPath, err := e.writeDefaultData(gvk, crdVersion)
	if err != nil {
		return nil, fmt.Errorf("could not write default data: %w", err)
	}

	cmd := exec.Command(e.EditorPath, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); errors.Is(err, &exec.ExitError{}) {
		return nil, fmt.Errorf("editor exited with error: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("could not run editor: %w", err)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("could not read temp file: %w", err)
	}

	dataMap := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &dataMap); err != nil {
		return nil, fmt.Errorf("could not unmarshal data: %w", err)
	}

	return dataMap, nil
}
