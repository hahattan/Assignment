package filesystem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFS_CreateFile(t *testing.T) {
	f := NewFS()
	ts := time.Now().Unix()

	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folderName, "", ts)
	require.NoError(t, err)

	// invalid user
	_, err = f.CreateFile("notFound", folderName, filename, "", ts)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid folder
	_, err = f.CreateFile(username, "notFound", filename, "", ts)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid chars
	_, err = f.CreateFile(username, folderName, "file/name", "", ts)
	assert.ErrorIs(t, err, ErrInvalidCharacter)

	// create file
	file, err := f.CreateFile(username, folderName, filename, description, ts)
	require.NoError(t, err)
	assert.Equal(t, filename, file.Name)
	assert.Equal(t, description, file.Description)
	assert.Equal(t, ts, file.UpdatedAt)

	// create file
	_, err = f.CreateFile(username, folderName, filename, description, ts)
	assert.ErrorIs(t, err, ErrDataAlreadyExists)
}

func TestFS_DeleteFile(t *testing.T) {
	f := NewFS()
	ts := time.Now().Unix()

	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folderName, "", ts)
	require.NoError(t, err)
	_, err = f.CreateFile(username, folderName, filename, description, ts)
	require.NoError(t, err)

	// invalid user
	err = f.DeleteFile("notFound", folderName, filename)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid folder
	err = f.DeleteFile(username, "notFound", filename)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid file
	err = f.DeleteFile(username, folderName, "notFound")
	assert.ErrorIs(t, err, ErrDataNotFound)

	// delete file
	err = f.DeleteFile(username, folderName, filename)
	require.NoError(t, err)

	// delete file again
	err = f.DeleteFile(username, folderName, filename)
	assert.ErrorIs(t, err, ErrDataNotFound)
}

func TestFS_ListFile(t *testing.T) {
	const (
		file1 = "file1"
		file2 = "file2"
		file3 = "file3"
	)
	f := NewFS()
	ts := time.Now().Unix()
	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folderName, "", ts-3)
	require.NoError(t, err)
	_, err = f.CreateFile(username, folderName, file1, "", ts-2)
	require.NoError(t, err)
	_, err = f.CreateFile(username, folderName, file2, "", ts-1)
	require.NoError(t, err)
	_, err = f.CreateFile(username, folderName, file3, "", ts)
	require.NoError(t, err)

	defaultOpt := NewDefaultSortOption()
	files, err := f.ListFile(username, folderName, defaultOpt)
	require.NoError(t, err)
	require.Len(t, files, 3)
	assert.Equal(t, file1, files[0].Name)
	assert.Equal(t, file2, files[1].Name)
	assert.Equal(t, file3, files[2].Name)

	nameDesc := defaultOpt
	nameDesc.Order = SortingOrderDesc
	files, err = f.ListFile(username, folderName, nameDesc)
	require.NoError(t, err)
	require.Len(t, files, 3)
	assert.Equal(t, file3, files[0].Name)
	assert.Equal(t, file2, files[1].Name)
	assert.Equal(t, file1, files[2].Name)

	createdAsc := defaultOpt
	createdAsc.Field = SortingFieldCreatedTime
	files, err = f.ListFile(username, folderName, createdAsc)
	require.NoError(t, err)
	require.Len(t, files, 3)
	assert.Equal(t, file1, files[0].Name)
	assert.Equal(t, file2, files[1].Name)
	assert.Equal(t, file3, files[2].Name)

	createdDesc := defaultOpt
	createdDesc.Field = SortingFieldCreatedTime
	createdDesc.Order = SortingOrderDesc
	files, err = f.ListFile(username, folderName, createdDesc)
	require.NoError(t, err)
	require.Len(t, files, 3)
	assert.Equal(t, file3, files[0].Name)
	assert.Equal(t, file2, files[1].Name)
	assert.Equal(t, file1, files[2].Name)
}
