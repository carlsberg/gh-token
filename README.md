# GitHub Token

> [GitHub CLI](https://cli.github.com) extension to create and revoke installation tokens.

## Installation

```bash
$ gh extensions install carlsberg/gh-token
```

## Usage

```
$ gh token --help
Create and revoke GitHub Installation tokens

Usage:
  gh-token [command]

Examples:
# create a token with an installation ID
$ gh token create --app-id 123 app-private-key-path path/to/pem --installation-id 123

# create a token without installation ID
$ gh token create --app-id 123 app-private-key-path path/to/pem --org org-name

# revoke a token
$ gh token revoke --token ghs_123


Available Commands:
  completion  generate the autocompletion script for the specified shell
  create      
  help        Help about any command
  revoke      

Flags:
  -h, --help   help for gh-token

Use "gh-token [command] --help" for more information about a command.
```

## Contributing

Please, see the [Contributing Guide](.github/contributing.md) to learn how you can contribute
to this repository. Every contribution is welcome!

## License

This project is released under the [MIT License](LICENSE).
