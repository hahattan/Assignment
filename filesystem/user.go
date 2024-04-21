package filesystem

import "fmt"

func (f *FS) UserRegister(name string) (*User, error) {
	_, ok := f.Users[name]
	if ok {
		return nil, fmt.Errorf("%w: %s", ErrDataAlreadyExists, name)
	}

	if !isValidChar(name) {
		return nil, ErrInvalidCharacter
	}

	user := &User{
		Name:    name,
		Folders: make(map[string]*Folder),
	}

	f.Users[name] = user
	return user, nil
}
