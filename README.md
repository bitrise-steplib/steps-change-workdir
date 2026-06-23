# Change Working Directory for subsequent Steps

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-change-workdir?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-change-workdir/releases)

Changes the Working Directory for all subsequent Steps to a specified path.

<details>
<summary>Description</summary>

This Step changes the base Working Directory for every subsequent Step in the Workflow to the specified path.

Useful for projects which expect a fixed project/workspace hierarchy, like `Go` which requires the project inside a "workspace" with the structure: `$GOPATH/src/domain/user/project`.

### Configuring the Step

1. Set the **Directory path** to the path you want subsequent Steps to use as their Working Directory.
2. Use the **Create path if doesn't exist?** input to control whether the path is created when it is missing.

### Troubleshooting
If the Step fails with a "does not exist" error, either set **Create path if doesn't exist?** to `true`, or make sure an earlier Step creates the directory.
The provided path is resolved to an absolute path before it is applied, so relative paths are resolved against the current Working Directory.

### Related Steps
- [Change value in file](https://www.bitrise.io/integrations/steps/change-value)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `path` | The path to set as the base Working Directory for every subsequent Step. It is resolved to an absolute path before being applied. | required |  |
| `is_create_path` | When set to `true`, the Step creates the directory (including parent directories) if it does not exist yet. When set to `false`, the Step fails if the directory is missing. | required | `true` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-change-workdir/pulls) and [issues](https://github.com/bitrise-steplib/steps-change-workdir/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
