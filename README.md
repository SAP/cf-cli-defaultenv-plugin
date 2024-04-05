# DefaultEnv CF CLI Plugin

This is a Cloud Foundry CLI plugin designed to aid local development of [multi-target applications (MTAs)](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) in Cloud Foundry. The `default-env` command creates a local `default-env.json` file with the environment variables of the specified Cloud Foundry app - typically connection details for bound Cloud Foundry services such as SAP HANA HDI Containers, XSUAA (User Account and Authentication) and intra-MTA destinations defined in `mta.yaml`. Environment variables written to `default-env.json` include `VCAP_APPLICATION`, `VCAP_SERVICES` and `destinations`. The `default-env.json` file is used by [@sap/approuter](https://www.npmjs.com/package/@sap/approuter) and [@sap/hdi-deploy](https://www.npmjs.com/package/@sap/hdi-deploy) when running locally and it's also possible to use `default-env.json` from your own Node.js applications via [@sap/xsenv](https://www.npmjs.com/package/@sap/xsenv) as follows:

```javascript
const xsenv = require("@sap/xsenv");
xsenv.loadEnv();
```

# Requirements

Installed CloudFoundry CLI - ensure that CloudFoundry CLI is installed and working. For more information about installation of CloudFoundry CLI, please visit the official CloudFoundry [documentation](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).

# Download and Installation

Check whether you have a previous version installed, using the command: `cf plugins`. If the DefaultEnv plugin is already installed, you need to uninstall it first and then to install the new version. You can uninstall the plugin using command `cf uninstall-plugin DefaultEnv`.

## CF Community Plugin Repository

The DefaultEnv CF CLI Plugin is available on the CF Community Repository. To install the latest available version of the DefaultEnv CLI Plugin execute the following:

```console
cf install-plugin DefaultEnv
```

If you do not have the community repository in your CF CLI you can add it first by executing:

```console
cf add-plugin-repo CF-Community https://plugins.cloudfoundry.org
```

## Manual Installation

Alternatively you can install any version of the plugin by manually downloading it from the releases page and installing the binaries for your specific operating system.

The latest version of the plugin can be downloaded from the project's [releases](https://github.com/sap/cf-cli-defaultenv-plugin/releases/latest). Download the plugin for your platform (Darwin, Linux, Windows) and install the plugin, using the following command:

```console
cf install-plugin <path-to-the-binary> -f
```

> [!NOTE]
> If you are running on a Unix-based system, you need to make the plugin executable before installing it. In order to achieve this, execute the following commad
>
> ```console
> chmod +x <path-to-the-plugin>
> ```

## Usage

The DefaultEnv CF plugin supports the following commands:

| Command Name  | Command Description                                                                               |
| ------------- | ------------------------------------------------------------------------------------------------- |
| `default-env` | Create `default-env.json` file with environment variables of an app. Usage `cf default-env myapp` |

For more information, see the command help output available via `cf [command] --help` or `cf help [command]`.

# Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/cf-cli-defaultenv-plugin/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

# Security / Disclosure

If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/cf-cli-defaultenv-plugin/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

# Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

# Licensing

Copyright 2024 SAP SE or an SAP affiliate company and cf-cli-defaultenv-plugin contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/cf-cli-defaultenv-plugin).
