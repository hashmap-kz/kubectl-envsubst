package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/cmd/apply"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/yaml"
)

type EnvsubstApplyPlugin struct {
	configFlags     *genericclioptions.ConfigFlags
	ioStreams       genericiooptions.IOStreams
	filenameOptions resource.FilenameOptions
	allowedVars     []string
	allowedPrefixes []string
	strict          bool
	verbose         bool
	configFile      string
}

type Config struct {
	AllowedVars     []string `yaml:"allowedVars" json:"allowedVars"`
	AllowedPrefixes []string `yaml:"allowedPrefixes" json:"allowedPrefixes"`
	Strict          bool     `yaml:"strict" json:"strict"`
	Verbose         bool     `yaml:"verbose" json:"verbose"`
}

func NewEnvsubstApplyPlugin(streams genericiooptions.IOStreams) *EnvsubstApplyPlugin {
	return &EnvsubstApplyPlugin{
		configFlags: genericclioptions.NewConfigFlags(true),
		ioStreams:   streams,
	}
}

func (p *EnvsubstApplyPlugin) Complete(args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--filename", "-f":
			if i+1 < len(args) {
				p.filenameOptions.Filenames = append(p.filenameOptions.Filenames, args[i+1])
				i++
			} else {
				return fmt.Errorf("flag --filename requires a value")
			}
		case "--envsubst-allowed":
			if i+1 < len(args) {
				p.allowedVars = strings.Split(args[i+1], ",")
				i++
			} else {
				return fmt.Errorf("flag --envsubst-allowed requires a comma-separated list of environment variables")
			}
		case "--envsubst-allowed-with-prefixes":
			if i+1 < len(args) {
				p.allowedPrefixes = strings.Split(args[i+1], ",")
				i++
			} else {
				return fmt.Errorf("flag --envsubst-allowed-with-prefixes requires a comma-separated list of prefixes")
			}
		case "--strict":
			p.strict = true
		case "--verbose", "-v":
			p.verbose = true
		case "--config":
			if i+1 < len(args) {
				p.configFile = args[i+1]
				i++
			} else {
				return fmt.Errorf("flag --config requires a file path")
			}
		}
	}
	if len(p.filenameOptions.Filenames) == 0 {
		return fmt.Errorf("at least one file must be specified using --filename (-f)")
	}
	return nil
}

func (p *EnvsubstApplyPlugin) LoadConfig() error {
	if p.configFile == "" {
		return nil
	}

	// Read the config file
	data, err := os.ReadFile(p.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the config file
	var config Config
	if strings.HasSuffix(p.configFile, ".yaml") || strings.HasSuffix(p.configFile, ".yml") {
		err = yaml.Unmarshal(data, &config)
	} else if strings.HasSuffix(p.configFile, ".json") {
		err = json.Unmarshal(data, &config)
	} else {
		return fmt.Errorf("unsupported config file format: %s", p.configFile)
	}
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Merge config values
	if len(p.allowedVars) == 0 {
		p.allowedVars = config.AllowedVars
	}
	if len(p.allowedPrefixes) == 0 {
		p.allowedPrefixes = config.AllowedPrefixes
	}
	if !p.strict {
		p.strict = config.Strict
	}
	if !p.verbose {
		p.verbose = config.Verbose
	}

	return nil
}

func (p *EnvsubstApplyPlugin) Validate() error {
	if len(p.allowedVars) > 0 && len(p.allowedPrefixes) > 0 {
		return fmt.Errorf("--envsubst-allowed and --envsubst-allowed-with-prefixes are mutually exclusive")
	}
	if len(p.filenameOptions.Filenames) == 0 {
		return fmt.Errorf("no files specified")
	}
	return nil
}

func (p *EnvsubstApplyPlugin) substituteEnvVariables(input []byte) ([]byte, error) {
	text := string(input)

	// Collect allowed environment variables and prefixes
	envMap := make(map[string]string)
	for _, env := range p.allowedVars {
		if value, exists := os.LookupEnv(env); exists {
			envMap[env] = value
		}
	}

	// Add variables with allowed prefixes
	for _, prefix := range p.allowedPrefixes {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && strings.HasPrefix(parts[0], prefix) {
				envMap[parts[0]] = parts[1]
			}
		}
	}

	// Perform substitution using regex
	re := regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)
	substituted := re.ReplaceAllStringFunc(text, func(match string) string {
		varName := re.FindStringSubmatch(match)[1]
		if value, ok := envMap[varName]; ok {
			return value
		}
		if p.strict {
			return fmt.Sprintf("UNDEFINED(%s)", varName) // Marker for strict mode
		}
		return match
	})

	// Handle strict mode by detecting unresolved variables
	if p.strict {
		unresolved := re.FindAllString(substituted, -1)
		if len(unresolved) > 0 {
			return nil, fmt.Errorf("undefined variables: %v", unresolved)
		}
	}

	return []byte(substituted), nil
}

func (p *EnvsubstApplyPlugin) applyResources(manifest []byte) error {
	// Create a command factory
	f := cmdutil.NewFactory(p.configFlags)

	// Create the apply command
	cmd := apply.NewCmdApply("kubectl", f, p.ioStreams)

	// Simulate applying resources using the apply command
	cmd.SetArgs([]string{"-f", "-"})
	cmd.SetIn(bytes.NewReader(manifest))

	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("failed to apply resources: %v", err)
	}
	return nil
}

//func (p *EnvsubstApplyPlugin) applyResources(manifest []byte) error {
//
//	result := &cobra.Command{
//		Use:    "kubectl",
//		Short:  "kubectl controls the Kubernetes cluster manager",
//		Hidden: true,
//	}
//
//	genericFlags := genericclioptions.NewConfigFlags(true)
//	genericFlags.AddFlags(result.PersistentFlags())
//	matchVersionFlags := cmdutil.NewMatchVersionFlags(genericFlags)
//	matchVersionFlags.AddFlags(result.PersistentFlags())
//	f := cmdutil.NewFactory(matchVersionFlags)
//
//	cmd := apply.NewCmdApply("kubectl", f, p.ioStreams)
//
//	// Simulate applying resources using the apply command
//	cmd.SetArgs([]string{"-f", "-"})
//	cmd.SetIn(bytes.NewReader(manifest))
//
//	if err := cmd.Execute(); err != nil {
//		return fmt.Errorf("failed to apply resources: %v", err)
//	}
//	return nil
//}

func (p *EnvsubstApplyPlugin) Run() error {
	// Build resource builder
	builder := resource.NewBuilder(p.configFlags).
		Unstructured().
		ContinueOnError().
		FilenameParam(false, &p.filenameOptions)

	infos, err := builder.Do().Infos()
	if err != nil {
		return fmt.Errorf("failed to process files: %v", err)
	}

	var substitutedContent bytes.Buffer
	for _, info := range infos {
		fmt.Println(info)
		manifest := make([]byte, 0)
		if err != nil {
			return fmt.Errorf("failed to marshal object: %v", err)
		}

		substituted, err := p.substituteEnvVariables(manifest)
		if err != nil {
			return fmt.Errorf("failed to substitute environment variables: %v", err)
		}

		substitutedContent.Write(substituted)
		substitutedContent.WriteString("---\n")

		if p.verbose {
			fmt.Fprintln(p.ioStreams.Out, "Substituted Manifest:")
			fmt.Fprintln(p.ioStreams.Out, string(substituted))
		}
	}

	// Apply substituted resources
	err = p.applyResources(substitutedContent.Bytes())
	if err != nil {
		return fmt.Errorf("failed to apply resources: %v", err)
	}

	fmt.Fprintln(p.ioStreams.Out, "Resources successfully applied.")
	return nil
}

func main() {
	streams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	plugin := NewEnvsubstApplyPlugin(streams)

	// Parse arguments
	args := os.Args[1:]
	if err := plugin.Complete(args); err != nil {
		fmt.Fprintln(plugin.ioStreams.ErrOut, err)
		os.Exit(1)
	}

	// Load config file if provided
	if err := plugin.LoadConfig(); err != nil {
		fmt.Fprintln(plugin.ioStreams.ErrOut, err)
		os.Exit(1)
	}

	// Validate inputs
	if err := plugin.Validate(); err != nil {
		fmt.Fprintln(plugin.ioStreams.ErrOut, err)
		os.Exit(1)
	}

	// Run the plugin
	if err := plugin.Run(); err != nil {
		fmt.Fprintln(plugin.ioStreams.ErrOut, err)
		os.Exit(1)
	}
}
