package filesystem

type File struct {
	Name        string
	CreatedAt   int64
	UpdatedAt   int64
	Description string
	Content     string
}

type Folder struct {
	Name        string
	CreatedAt   int64
	UpdatedAt   int64
	Description string
	Files       map[string]*File
}

type User struct {
	Name    string
	Folders map[string]*Folder
}

type SortingField string
type SortingOrder string

const (
	SortingFieldName        SortingField = "--sort-name"
	SortingFieldCreatedTime SortingField = "--sort-created"
	SortingOrderAsc         SortingOrder = "asc"
	SortingOrderDesc        SortingOrder = "desc"
)

type SortOption struct {
	Field SortingField
	Order SortingOrder
}

func NewDefaultSortOption() SortOption {
	return SortOption{
		Field: SortingFieldName,
		Order: SortingOrderAsc,
	}
}
