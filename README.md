# DefaultEnv CF CLI Plugin

This is a Cloud Foundry CLI plugin designed to aid local development of [multi-target applications (MTAs)](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) in Cloud Foundry. The default-env command creates a local default-env.json file with the environment variables of the specified Cloud Foundry app - typically connection details for bound Cloud Foundry services such as SAP HANA HDI Containers, XSUAA (User Account and Authentication) and intra-MTA destinations defined in mta.yaml. Environment variables written to default-env.json include `VCAP_APPLICATION`, `VCAP_SERVICES` and `destinations`. The default-env.json file is used by [@sap/approuter](https://www.npmjs.com/package/@sap/approuter) and [@sap/hdi-deploy](https://www.npmjs.com/package/@sap/hdi-deploy) when running locally and it's also possible to use default-env.json from your own Node.js applications via [@sap/xsenv](https://www.npmjs.com/package/@sap/xsenv) as follows:

`const xsenv = require('@sap/xsenv');`

`xsenv.loadEnv();`

# Requirements

Installed CloudFoundry CLI - ensure that CloudFoundry CLI is installed and working. For more information about installation of CloudFoundry CLI, please visit the official CloudFoundry [documentation](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).

# Download and Installation

Check whether you have a previous version installed, using the command: `cf plugins`. If the DefaultEnv plugin is already installed, you need to uninstall it first and then to install the new version. You can uninstall the plugin using command `cf uninstall-plugin DefaultEnv`.

## CF Community Plugin Repository

The DefaultEnv CF CLI Plugin is available on the CF Community Repository. To install the latest available version of the DefaultEnv CLI Plugin execute the following:

`cf install-plugin DefaultEnv`

If you do not have the community repository in your CF CLI you can add it first by executing:

`cf add-plugin-repo CF-Community https://plugins.cloudfoundry.org`

## Manual Installation

Alternatively you can install any version of the plugin by manually downloading it from the releases page and installing the binaries for your specific operating system.

### Download

The latest version of the plugin can also be downloaded from the project's [releases](https://github.com/saphanaacademy/DefaultEnv/releases/latest). Download the plugin for your platform (Darwin, Linux, Windows). The name for the correct plugin for each platform can be found in the table below.

| Mac OS X 64 bit | Windows 32 bit   | Windows 64 bit   | Linux 32 bit       | Linux 64 bit       |
| --------------- | ---------------- | ---------------- | ------------------ | ------------------ |
| DefaultEnv.osx  | DefaultEnv.win32 | DefaultEnv.win64 | DefaultEnv.linux32 | DefaultEnv.linux64 |

### Installation

Install the plugin, using the following command:

```
cf install-plugin <path-to-the-plugin> -f
```

Note: if you are running on a Unix-based system, you need to make the plugin executable before installing it. In order to achieve this, execute the following commad `chmod +x <path-to-the-plugin>`

## Usage

The DefaultEnv CF plugin supports the following commands:

| Command Name  | Command Description                                                                                 |
| ------------- | --------------------------------------------------------------------------------------------------- |
| `default-env` | Create default-env.json file with environment variables of an app. Usage `cf default-env myapp-srv` |

For more information, see the command help output available via `cf [command] --help` or `cf help [command]`.

# Further Information

Tutorials:

- [SAP Business Technology Platform Onboarding](https://www.youtube.com/playlist?list=PLkzo92owKnVw3l4fqcLoQalyFi9K4-UdY)
- [SAP HANA Academy](https://www.youtube.com/saphanaacademy)

# License

This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the [LICENSE](https://github.com/saphanaacademy/DefaultEnv/blob/master/LICENSE) file.
