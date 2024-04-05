package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type DefaultEnvPlugin struct{}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func (c *DefaultEnvPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "default-env" {
		if len(args) != 2 {
			fmt.Println("Please specify an app")
			return
		}
		app, err := cliConnection.GetApp(args[1])
		handleError(err)
		url := fmt.Sprintf("/v2/apps/%s/env", app.Guid)
		env, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)
		handleError(err)
		var envJSON map[string]interface{}
		json.Unmarshal([]byte(strings.Join(env, "")), &envJSON)
		f, err := os.Create("default-env.json")
		handleError(err)
		_, err = f.Write([]byte("{"))
		handleError(err)
		env1, err := json.Marshal(envJSON["system_env_json"])
		handleError(err)
		str1 := strings.Trim(string(env1), "{}")
		_, err = f.Write([]byte(str1))
		handleError(err)
		_, err = f.Write([]byte("},"))
		handleError(err)
		env2, err := json.Marshal(envJSON["application_env_json"])
		handleError(err)
		str2 := strings.Trim(string(env2), "{}")
		_, err = f.Write([]byte(str2))
		handleError(err)
		_, err = f.Write([]byte("},"))
		handleError(err)
		env3, err := json.Marshal(envJSON["environment_json"])
		handleError(err)
		str3 := strings.Trim(string(env3), "{}")
		_, err = f.Write([]byte(str3))
		handleError(err)
		_, err = f.Write([]byte("}"))
		handleError(err)
		err = f.Close()
		handleError(err)
		fmt.Println("Environment variables for " + args[1] + " written to default-env.json")
	}
}

func (c *DefaultEnvPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DefaultEnv",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
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
