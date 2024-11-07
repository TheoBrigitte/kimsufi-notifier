# Run

This documentation describe how to run the check script using different CI services.

Examples are provided, they will run the check script every hour and send a notifications to OpsGenie. You can customize them to fit your needs using environment variables described in [config.env.example](config.env.example) file.

## GitHub Actions

Create a new workflow file in your repository under `.github/workflows/` directory. You can use the provided example [check.yaml.example](.github/workflows/check.yaml.example) file.

You can set environement variable using Github's `Repository secrets` at `https://github.com/<owner>/kimsufi-notifier/settings/secrets/actions`.

More info on [Github Actions](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions).

## CircleCI

Create a new config file in your repository under `.circleci/config.yml`. You can use the provided example [config.yml.example](.circleci/config.yml.example).

You can set environement variable using CircleCI [contexts](https://circleci.com/docs/contexts/#create-and-use-a-context. Alternatively you can use CircleCI [Project's environement variables](https://circleci.com/docs/set-environment-variable/#set-an-environment-variable-in-a-project).

