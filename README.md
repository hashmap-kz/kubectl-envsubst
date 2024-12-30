# kubectl-envsubst - a plugin for kubectl, used for expand env-vars in manifests

### Features:

- --only: cmd flag, that consumes a list of names that allowed for expansion
- --prefix-only: cmd flag, that consumes a prefix (APP_), and variables that not match the prefix will be ignored 

### Client usage:
```
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
```