package filesystem

type FileSystem interface {
	UserRegister(name string) (*User, error)
	CreateFolder(username string, folderName string, description string, ts int64) (*Folder, error)
	DeleteFolder(username string, folderName string) error
	RenameFolder(username string, folderName string, newFolderName string, ts int64) (*Folder, error)
	ListFolder(username string, opt SortOption) ([]*Folder, error)
	CreateFile(username string, folderName string, fileName string, description string, ts int64) (*File, error)
	DeleteFile(username string, folderName string, fileName string) error
	ListFile(username string, folderName string, opt SortOption) ([]*File, error)
}

type FS struct {
	Users map[string]*User
}

func NewFS() *FS {
	return &FS{
		Users: make(map[string]*User),
	}
}
