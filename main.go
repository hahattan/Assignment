package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/hahattan/assignment/filesystem"
)

func main() {
	fs := filesystem.NewFS()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("# ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // Remove newline character

		commandList := strings.Split(input, " ")
		command := commandList[0]
		switch command {
		case "register":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: register [username]\n")
				continue
			}
			username := strings.ToLower(commandList[1])
			user, err := fs.UserRegister(username)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Add %s successfully.\n", user.Name)
		case "create-folder":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: create-folder [username] [foldername] [description]?\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: create-folder [username] [foldername] [description]?\n")
				continue
			}

			ts := time.Now().Unix()
			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			var description string
			if len(args) > 2 {
				description = args[2]
			}
			folder, err := fs.CreateFolder(username, folderName, description, ts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Create %s successfully.\n", folder.Name)
		case "delete-folder":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: delete-folder [username] [foldername]\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: delete-folder [username] [foldername]\n")
				continue
			}

			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			err := fs.DeleteFolder(username, folderName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Delete %s successfully.\n", folderName)
		case "rename-folder":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: rename-folder [username] [foldername] [new-folder-name]\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 3 {
				fmt.Fprintf(os.Stderr, "Usage: rename-folder [username] [foldername] [new-folder-name]\n")
				continue
			}

			ts := time.Now().Unix()
			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			newFolderName := strings.ToLower(args[2])
			folder, err := fs.RenameFolder(username, folderName, newFolderName, ts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Rename %s to %s successfully.\n", folderName, folder.Name)
		case "list-folders":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: list-folders [username] [--sort-name|--sort-created]? [asc|desc]?\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 1 {
				fmt.Fprintf(os.Stderr, "Usage: list-folders [username] [--sort-name|--sort-created]? [asc|desc]?\n")
				continue
			}

			username := strings.ToLower(args[0])
			opt := filesystem.NewDefaultSortOption()
			if len(args) > 1 {
				sortBy := strings.ToLower(args[1])
				if sortBy != "--sort-name" && sortBy != "--sort-created" {
					fmt.Fprintf(os.Stderr, "Usage: list-folders [username] [--sort-name|--sort-created]? [asc|desc]?\n")
					continue
				}
				if sortBy == "--sort-created" {
					opt.Field = filesystem.SortingFieldCreatedTime
				}
			}
			if len(args) > 2 {
				sortOrder := strings.ToLower(args[2])
				if sortOrder != "asc" && sortOrder != "desc" {
					fmt.Fprintf(os.Stderr, "Usage: list-folders [username] [--sort-name|--sort-created]? [asc|desc]?\n")
					continue
				}
				if sortOrder == "desc" {
					opt.Order = filesystem.SortingOrderDesc
				}
			}

			folders, err := fs.ListFolder(username, opt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}

			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			for _, folder := range folders {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", folder.Name, folder.Description, time.Unix(folder.CreatedAt, 0).Format(time.RFC822))
			}
			w.Flush()
		case "create-file":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: create-file [username] [foldername] [filename] [description]?\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 3 {
				fmt.Fprintf(os.Stderr, "Usage: create-file [username] [foldername] [filename] [description]?\n")
				continue
			}

			ts := time.Now().Unix()
			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			filename := strings.ToLower(args[2])
			var description string
			if len(args) > 3 {
				description = args[3]
			}
			file, err := fs.CreateFile(username, folderName, filename, description, ts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Create %s successfully.\n", file.Name)
		case "delete-file":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: delete-file [username] [foldername] [filename]\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 3 {
				fmt.Fprintf(os.Stderr, "Usage: delete-file [username] [foldername] [filename]\n")
				continue
			}

			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			filename := strings.ToLower(args[2])
			err := fs.DeleteFile(username, folderName, filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "Delete %s successfully.\n", filename)
		case "list-files":
			if len(commandList) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created]? [asc|desc]?\n")
				continue
			}
			args := commandList[1:]
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created]? [asc|desc]?\n")
				continue
			}

			username := strings.ToLower(args[0])
			folderName := strings.ToLower(args[1])
			opt := filesystem.NewDefaultSortOption()
			if len(args) > 2 {
				sortBy := strings.ToLower(args[2])
				if sortBy != "--sort-name" && sortBy != "--sort-created" {
					fmt.Fprintf(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created]? [asc|desc]?\n")
					continue
				}
				if sortBy == "--sort-created" {
					opt.Field = filesystem.SortingFieldCreatedTime
				}
			}
			if len(args) > 3 {
				sortOrder := strings.ToLower(args[3])
				if sortOrder != "asc" && sortOrder != "desc" {
					fmt.Fprintf(os.Stderr, "Usage: list-files [username] [foldername] [--sort-name|--sort-created]? [asc|desc]?\n")
					continue
				}
				if sortOrder == "desc" {
					opt.Order = filesystem.SortingOrderDesc
				}
			}

			files, err := fs.ListFile(username, folderName, opt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}

			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			for _, file := range files {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", file.Name, file.Description, time.Unix(file.CreatedAt, 0).Format(time.RFC822))
			}
			w.Flush()
		case "exit":
			return
		default:
			fmt.Fprint(os.Stderr, "Error: Unknown command\n")
		}
	}
}
