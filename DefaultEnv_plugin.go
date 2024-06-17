package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// DefaultEnvPlugin allows users to export environment variables of an app into a JSON file
type DefaultEnvPlugin struct{}

var (
	// ErrAppNotSpecified is returned when the app name is not provided by the user
	ErrAppNotSpecified = fmt.Errorf("please specify an app")
)

func (*DefaultEnvPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "default-env" {
		return
	}
	if err := runDefaultEnv(cliConnection, args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func (c *DefaultEnvPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DefaultEnv",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 7,
			Minor: 2,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "default-env",
				Alias:    "de",
				HelpText: "Create default-env.json file with environment variables of an app.",
				UsageDetails: plugin.Usage{
					Usage: "cf default-env APP",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DefaultEnvPlugin))
}

// environmentResponse from /v3/apps/:guid/env
type environmentResponse struct {
	SystemEnvJson        map[string]interface{} `json:"system_env_json"`
	ApplicationEnvJson   map[string]interface{} `json:"application_env_json"`
	EnvironmentVariables map[string]interface{} `json:"environment_variables"`
}

// Merge all environment variables into one map
func (e environmentResponse) Merge() map[string]interface{} {
	content := make(map[string]interface{})
	for k, v := range e.SystemEnvJson {
		content[k] = v
	}
	for k, v := range e.ApplicationEnvJson {
		content[k] = v
	}
	for k, v := range e.EnvironmentVariables {
		content[k] = v
	}
	return content
}

// marshalAndWrite marshals v (any) into JSON and writes it to a file
func marshalAndWrite(v interface{}, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

// runDefaultEnv fetches and merges the environment variables of a specified CF app into a JSON file
func runDefaultEnv(cliConnection plugin.CliConnection, args []string) error {
	if len(args) != 2 {
		return ErrAppNotSpecified
	}

	app, err := cliConnection.GetApp(args[1])
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v3/apps/%s/env", app.Guid)
	env, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return err
	}

	var data environmentResponse
	if err = json.Unmarshal([]byte(strings.Join(env, "")), &data); err != nil {
		return err
	}

	if err = marshalAndWrite(data.Merge(), "default-env.json"); err != nil {
		return err
	}

	fmt.Println("Environment variables for " + args[1] + " written to default-env.json")
	return nil
}
