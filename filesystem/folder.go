package filesystem

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

func (f *FS) CreateFolder(username string, folderName string, description string, ts int64) (*Folder, error) {
	user, ok := f.Users[username]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	_, ok = user.Folders[folderName]
	if ok {
		return nil, fmt.Errorf("%w: %s", ErrDataAlreadyExists, folderName)
	}

	if !isValidChar(folderName) {
		return nil, ErrInvalidCharacter
	}

	folder := &Folder{
		Name:        folderName,
		CreatedAt:   ts,
		UpdatedAt:   ts,
		Description: description,
		Files:       make(map[string]*File),
	}

	user.Folders[folderName] = folder
	return folder, nil
}

func (f *FS) DeleteFolder(username string, folderName string) error {
	user, ok := f.Users[username]
	if !ok {
		return fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	_, ok = user.Folders[folderName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrDataNotFound, folderName)
	}

	delete(user.Folders, folderName)
	return nil
}

func (f *FS) RenameFolder(username string, folderName string, newFolderName string, ts int64) (*Folder, error) {
	user, ok := f.Users[username]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	folder, ok := user.Folders[folderName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, folderName)
	}

	if !isValidChar(newFolderName) {
		return nil, ErrInvalidCharacter
	}

	updatedFolder := folder
	updatedFolder.Name = newFolderName
	updatedFolder.UpdatedAt = ts

	delete(user.Folders, folderName)
	user.Folders[newFolderName] = updatedFolder

	return updatedFolder, nil
}

func (f *FS) ListFolder(username string, opt SortOption) ([]*Folder, error) {
	user, ok := f.Users[username]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	folders := lo.Values(user.Folders)
	var compareFunc func(i, j int) bool

	if opt.Field == SortingFieldName {
		compareFunc = func(i, j int) bool {
			if opt.Order == SortingOrderAsc {
				return folders[i].Name < folders[j].Name
			}
			return folders[i].Name > folders[j].Name
		}
	} else if opt.Field == SortingFieldCreatedTime {
		compareFunc = func(i, j int) bool {
			if opt.Order == SortingOrderAsc {
				return folders[i].CreatedAt < folders[j].CreatedAt
			}
			return folders[i].CreatedAt > folders[j].CreatedAt
		}
	}

	sortFolders(folders, compareFunc)
	return folders, nil
}

func sortFolders(folders []*Folder, compareFunc func(i, j int) bool) {
	sort.Slice(folders, compareFunc)
}
