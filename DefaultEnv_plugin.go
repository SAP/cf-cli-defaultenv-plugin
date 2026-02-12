package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	cfclient "github.com/cloudfoundry/go-cfclient/v3/client"
	cfconfig "github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

// DefaultEnvPlugin allows users to export environment variables of an app into a JSON file
type DefaultEnvPlugin struct{}

func (*DefaultEnvPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "default-env" {
		return
	}
	if err := runDefaultEnv(cliConnection, args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func (c *DefaultEnvPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DefaultEnv",
		Version: plugin.VersionType{
			Major: 2,
			Minor: 0,
			Build: 0,
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
					Options: map[string]string{
						"f":      "output file name (default: default-env.json)",
						"guid":   "specify the app GUID directly instead of app name",
						"stdout": "write output to stdout instead of a file",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DefaultEnvPlugin))
}

// createClient creates a new CF client using the access token and API endpoint from the CLI connection
func createClient(connection plugin.CliConnection) (*cfclient.Client, error) {
	accessToken, err := connection.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}

	endpoint, err := connection.ApiEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve API endpoint: %w", err)
	}

	cfg, err := cfconfig.New(endpoint, cfconfig.Token(accessToken, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create CF client config: %w", err)
	}

	client, err := cfclient.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create CF client: %w", err)
	}

	return client, nil
}

// findAppByName retrieves the app with the specified name in the given space
func findAppByName(ctx context.Context, client *cfclient.Client, appName, spaceGUID string) (*resource.App, error) {
	app, err := client.Applications.Single(ctx, &cfclient.AppListOptions{
		Names:      cfclient.Filter{Values: []string{appName}},
		SpaceGUIDs: cfclient.Filter{Values: []string{spaceGUID}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app: %w", err)
	}
	return app, nil
}

// findEnvironmentByAppGUID retrieves the environment variables of the app with the specified GUID
func findEnvironmentByAppGUID(ctx context.Context, client *cfclient.Client, appGUID string) (*AppEnvironment, error) {
	seg := url.PathEscape(appGUID)
	p, err := url.JoinPath("/v3/apps", seg, "env")
	if err != nil {
		return nil, fmt.Errorf("failed to construct API path: %w", err)
	}

	var environment AppEnvironment
	req, err := http.NewRequestWithContext(ctx, "GET", client.ApiURL(p), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := client.ExecuteAuthRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := json.NewDecoder(resp.Body).Decode(&environment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &environment, nil
}

// runDefaultEnv fetches and merges the environment variables of a specified CF app into a JSON file
func runDefaultEnv(cliConnection plugin.CliConnection, args []string) error {
	fs := flag.NewFlagSet("cf-defaultenv-plugin", flag.ExitOnError)
	outFileFlag := fs.String("f", "default-env.json", "output file name")
	appGUIDFlag := fs.String("guid", "", "specify the app GUID directly instead of app name")
	writeStdoutFlag := fs.Bool("stdout", false, "write output to stdout instead of a file")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	ctx := context.Background()

	client, err := createClient(cliConnection)
	if err != nil {
		return err
	}

	currentSpace, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return fmt.Errorf("failed to retrieve current space: %w", err)
	}

	var appGUID string
	if *appGUIDFlag != "" {
		// if the user specified the app GUID directly, use it instead of looking up the app by name
		appGUID = *appGUIDFlag

		_, _ = fmt.Fprintln(os.Stderr, "info: using provided app GUID", appGUID)
	} else {
		appName := fs.Arg(0)
		if appName == "" {
			_, _ = fmt.Fprintln(os.Stderr, "error: app name or -guid flag must be specified")
			fs.Usage()
			return nil
		}

		app, err := findAppByName(ctx, client, appName, currentSpace.Guid)
		if err != nil {
			return err
		}
		appGUID = app.GUID

		_, _ = fmt.Fprintln(os.Stderr, "info: retrieved app", appName, "with GUID", appGUID)
	}

	env, err := findEnvironmentByAppGUID(ctx, client, appGUID)
	if err != nil {
		return fmt.Errorf("failed to retrieve app environment: %w", err)
	}

	result := Merge(env.SystemEnvVars, env.AppEnvVars, env.EnvVars)

	if *writeStdoutFlag {
		if err := marshalAndWriteStdout(result); err != nil {
			return err
		}
		_, _ = fmt.Fprintln(os.Stderr, "success: environment variables written to stdout")
	} else {
		if err := marshalAndWrite(result, *outFileFlag); err != nil {
			return err
		}
		_, _ = fmt.Fprintln(os.Stderr, "success: environment variables written to", *outFileFlag)
	}
	return nil
}

// AppEnvironment is a stripped down version of [cfclient.Environment]
// that only contains the fields we care about for this plugin with a slightly different structure to make it easier
// to merge into a single map.
type AppEnvironment struct {
	EnvVars       map[string]any `json:"environment_variables,omitempty"`
	SystemEnvVars map[string]any `json:"system_env_json,omitempty"`      // VCAP_SERVICES
	AppEnvVars    map[string]any `json:"application_env_json,omitempty"` // VCAP_APPLICATION
}

// Merge merges multiple maps into a single map
func Merge[Map ~map[K]V, K comparable, V any](maps ...Map) Map {
	result := make(Map)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// marshalAndWrite marshals [v] into JSON and writes it to a file
func marshalAndWrite(v any, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err := json.NewEncoder(f).Encode(v); err != nil {
		return err
	}
	return nil
}

// marshalAndWriteStdout [v] into JSON and writes it to stdout
func marshalAndWriteStdout(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
