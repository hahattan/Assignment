# Virtual File System

This project is a simple filesystem simulation written in Go. It allows users to create, delete, and list folders and files.

## Build

### Native binaries

#### Prerequisites

- [Go 1.22](https://go.dev/doc/install)

Use the `Makefile` in the root directory:
```shell
make build
```

## Usage

The application supports the following commands:

- `register [username]`: Register a new user.
- `create-folder [username] [foldername] [description]?`: Create a new folder for a user.
- `delete-folder [username] [foldername]`: Delete a folder for a user.
- `rename-folder [username] [foldername] [new-folder-name]`: Rename a folder for a user.
- `list-folders [username] [--sort-name|--sort-created]? [asc|desc]?`: List all folders for a user.
- `create-file [username] [foldername] [filename] [description]?`: Create a new file in a user's folder.
- `delete-file [username] [foldername] [filename]`: Delete a file in a user's folder.
- `list-files [username] [foldername] [--sort-name|--sort-created]? [asc|desc]?`: List all files in a user's folder.

## Input Validation

For all commands that require user input, the following validation rules apply:

- The length of the input must be less than 20 characters.
- The input can only contain letters (a-z, A-Z), numbers (0-9), and spaces.

These rules apply to `username`, `foldername` and `filename`.
