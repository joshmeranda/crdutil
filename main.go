package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joshmeranda/crdutil/filler"
	"github.com/joshmeranda/crdutil/getter"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	CrdFile    string
	CrdVersion string
)

func init() {
	flag.StringVar(&CrdFile, "f", "", "the file to read the base crd from")
	flag.StringVar(&CrdVersion, "v", "v1", "the crd version to inspect")
}

func run() error {
	flag.Parse()

	var crdGetter getter.Getter
	if CrdFile != "" {
		crdGetter = getter.FromFile(CrdFile)
	} else {
		flag.Usage()
		return fmt.Errorf("expected one of '-f' to be specified")
	}

	f, err := filler.DefaultEditor()
	if err != nil {
		return fmt.Errorf("could not create default filler: %w", err)
	}

	crd, err := crdGetter(context.TODO())
	if err != nil {
		return fmt.Errorf("")
	}

	var crdVersion *apiextensionsv1.CustomResourceDefinitionVersion
	for _, v := range crd.Spec.Versions {
		if v.Name == CrdVersion {
			crdVersion = &v
		}
	}

	if crdVersion == nil {
		return fmt.Errorf("crd had no matching crd versions '%s'", CrdVersion)
	}

	data, err := f.Fill(crdVersion)
	if err != nil {
		return fmt.Errorf("could not fill data: %w", err)
	}

	fmt.Printf("%+v", data)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
