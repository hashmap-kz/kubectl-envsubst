package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
	"os"
)

func complex() {
	file := os.Args[1]
	client, err := util.GetClient()
	if err != nil {
		log.Fatal(err)
	}

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

func simple() {
	flags, err := util.ParseCmdFlags(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	part1 := []string{}
	part1 = append(part1, "--namespace")
	part1 = append(part1, flags.Namespace)
	for _, f := range flags.Filenames {
		part1 = append(part1, "-f")
		part1 = append(part1, f)
	}
	part1 = append(part1, flags.Others...)

	cmd, err := util.ExecCmd("kubectl", part1...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmd.CombinedOutput())
}

func main() {
	simple()
}
