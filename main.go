package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
	"os"
)

func main() {
	file := os.Args[1]
	client := util.GetClient()

	readFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(readFile, 4096)
	for {
		obj := &unstructured.Unstructured{}
		if err := decoder.Decode(obj); err != nil {
			break
		}
		if err := client.Apply(obj); err != nil {
			log.Fatal(fmt.Errorf("failed to apply resource: %w", err))
		}
	}

}
