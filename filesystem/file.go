package filesystem

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

func (f *FS) CreateFile(username string, folderName string, fileName string, description string, ts int64) (*File, error) {
	user, ok := f.Users[username]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	folder, ok := user.Folders[folderName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, folderName)
	}

	_, ok = folder.Files[fileName]
	if ok {
		return nil, fmt.Errorf("%w: %s", ErrDataAlreadyExists, fileName)
	}

	if !isValidChar(fileName) {
		return nil, ErrInvalidCharacter
	}

	file := &File{
		Name:        fileName,
		CreatedAt:   ts,
		UpdatedAt:   ts,
		Description: description,
	}

	folder.Files[fileName] = file
	return file, nil
}

func (f *FS) DeleteFile(username string, folderName string, fileName string) error {
	user, ok := f.Users[username]
	if !ok {
		return fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	folder, ok := user.Folders[folderName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrDataNotFound, folderName)
	}

	_, ok = folder.Files[fileName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrDataNotFound, fileName)
	}

	delete(folder.Files, fileName)
	return nil
}

func (f *FS) ListFile(username string, folderName string, opt SortOption) ([]*File, error) {
	user, ok := f.Users[username]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, username)
	}

	folder, ok := user.Folders[folderName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDataNotFound, folderName)
	}

	files := lo.Values(folder.Files)
	var compareFunc func(i, j int) bool

	if opt.Field == SortingFieldName {
		compareFunc = func(i, j int) bool {
			if opt.Order == SortingOrderAsc {
				return files[i].Name < files[j].Name
			}
			return files[i].Name > files[j].Name
		}
	} else if opt.Field == SortingFieldCreatedTime {
		compareFunc = func(i, j int) bool {
			if opt.Order == SortingOrderAsc {
				return files[i].CreatedAt < files[j].CreatedAt
			}
			return files[i].CreatedAt > files[j].CreatedAt
		}
	}

	sortFiles(files, compareFunc)
	return files, nil
}

func sortFiles(files []*File, compareFunc func(i, j int) bool) {
	sort.Slice(files, compareFunc)
}
