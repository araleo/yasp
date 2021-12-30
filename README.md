# YASP - Yet Another Snitch Program

Yasp is a simplified snitch cli tool to help with code maitenance, heavily inspired by (copied from) [/tsoding/snitch](https://github.com/tsoding/snitch). 

## Features

* Pretty print a directory file structure.
* Print all prints and todos statements found within the code.
* List open issues found in GitLab and submit new ones.
* Supports ignore like files to ignore files and directories for all commands.
* Customization of print, todos and issues commands.

## Usage

### Quickstart

The default directory is the cwd and the default ignore file is .yaspignore.

```
go build .

./snitch -c <command>
```

Available commands:

* ls: pretty prints the directory file structure.
* print: prints all the print statements found within the directory files.
* todo: prints all the todo statements.
* diag: prints both the print and todo statements.
* issues: prints project's GitLab issues.
* snitch: finds every instance of TODO! in the code and submits it as a new GitLab issue.


### Optional flags

- -d /path/to/dir
  - Absolute path to the root dir in which the program will run.
- -i /path/to/file
  - Absolute path to the ignore like file.


### Customization

The default ignore file, and the supported print and todo commands can be changed in the yasp.yml file.

### Env

To use GitLab's issues functionalities some environment variables must be set in a .env file according to the following:

```
GITLAB_TOKEN=<token>
GITLAB_API_URL=https://gitlab.com/api/v4/projects
GITLAB_PROJECT_ID=<project id>
```