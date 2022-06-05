# semver-cli
[![build](https://github.com/JFWenisch/semver-cli/actions/workflows/build-go.yml/badge.svg?branch=main)](https://github.com/JFWenisch/semver-cli/actions/workflows/build-go.yml)

semver-cli is an application that eases the creation of tags and releases during CI. semver cli is written in go using cobra.

## Quick start
```
# Get the latest tag found in the repository (local & remote)
$ semver-cli tags list -l
v11.22.33

# Output the tag on minor changes
$ semver-cli tags bump -t minor -d -p v
v11.23.0
```
## Usage
Given a version number MAJOR.MINOR.PATCH, increment the:

- MAJOR version when you make incompatible API changes,
- MINOR version when you add functionality in a backwards compatible manner, and
- PATCH version when you make backwards compatible bug fixes.



```
$ semver-cli tags bump -h
Funcitionality to identify and/or create the next tag in relation to semver conventions

Usage:
  semver-cli tags bump [flags]

Flags:
  -d, --dry-run         Outputs the next determined version without creating it
  -p, --prefix string   The prefix for tagging e.g. 'v'
  -t, --type string     Type of commit, e.g. 'major', 'minor' or 'patch'
```
## Build
```
go get -d -v ./...
go install -v ./...
go build
```

