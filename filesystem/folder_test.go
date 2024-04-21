package filesystem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	username    = "testUser"
	folderName  = "testFolder"
	filename    = "testFile"
	description = "testDescription"
)

func TestFS_CreateFolder(t *testing.T) {
	f := NewFS()
	ts := time.Now().Unix()

	_, err := f.UserRegister(username)
	require.NoError(t, err)

	// invalid username
	_, err = f.CreateFolder("not-found", folderName, "", ts)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid chars
	_, err = f.CreateFolder(username, "folder/name", "", ts)
	assert.ErrorIs(t, err, ErrInvalidCharacter)

	// create folder
	folder, err := f.CreateFolder(username, folderName, description, ts)
	require.NoError(t, err)
	assert.Equal(t, folderName, folder.Name)
	assert.Equal(t, description, folder.Description)
	assert.Equal(t, ts, folder.UpdatedAt)

	// create folder again
	_, err = f.CreateFolder(username, folderName, description, ts)
	assert.ErrorIs(t, err, ErrDataAlreadyExists)
}

func TestFS_DeleteFolder(t *testing.T) {
	f := NewFS()
	ts := time.Now().Unix()

	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folderName, description, ts)
	require.NoError(t, err)

	// invalid username
	err = f.DeleteFolder("notFound", folderName)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// invalid folder name
	err = f.DeleteFolder(username, "notFound")
	assert.ErrorIs(t, err, ErrDataNotFound)

	// delete folder
	err = f.DeleteFolder(username, folderName)
	require.NoError(t, err)

	// delete folder again
	err = f.DeleteFolder(username, folderName)
	assert.ErrorIs(t, err, ErrDataNotFound)
}

func TestFS_RenameFolder(t *testing.T) {
	const (
		newFolderName = "newTestFolder"
	)
	f := NewFS()
	ts := time.Now().Unix()

	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folderName, description, ts-1)
	require.NoError(t, err)

	// invalid username
	_, err = f.RenameFolder("notFound", folderName, newFolderName, ts)
	assert.ErrorIs(t, err, ErrDataNotFound)
	// invalid folder name
	_, err = f.RenameFolder(username, "notFound", newFolderName, ts)
	assert.ErrorIs(t, err, ErrDataNotFound)

	// rename folder
	folder, err := f.RenameFolder(username, folderName, newFolderName, ts)
	require.NoError(t, err)
	assert.Equal(t, newFolderName, folder.Name)
	assert.Equal(t, ts, folder.UpdatedAt)

	// rename folder again
	_, err = f.RenameFolder(username, folderName, newFolderName, ts)
	assert.ErrorIs(t, err, ErrDataNotFound)
}

func TestFS_ListFolder(t *testing.T) {
	const (
		folder1 = "folder1"
		folder2 = "folder2"
		folder3 = "folder3"
	)
	f := NewFS()
	ts := time.Now().Unix()
	_, err := f.UserRegister(username)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folder1, description, ts-2)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folder2, description, ts-1)
	require.NoError(t, err)
	_, err = f.CreateFolder(username, folder3, description, ts)
	require.NoError(t, err)

	defaultOpt := NewDefaultSortOption()
	folders, err := f.ListFolder(username, defaultOpt)
	require.NoError(t, err)
	require.Len(t, folders, 3)
	assert.Equal(t, folder1, folders[0].Name)
	assert.Equal(t, folder2, folders[1].Name)
	assert.Equal(t, folder3, folders[2].Name)

	nameDesc := defaultOpt
	nameDesc.Order = SortingOrderDesc
	folders, err = f.ListFolder(username, nameDesc)
	require.NoError(t, err)
	require.Len(t, folders, 3)
	assert.Equal(t, folder3, folders[0].Name)
	assert.Equal(t, folder2, folders[1].Name)
	assert.Equal(t, folder1, folders[2].Name)

	createdAsc := defaultOpt
	createdAsc.Field = SortingFieldCreatedTime
	folders, err = f.ListFolder(username, nameDesc)
	require.NoError(t, err)
	require.Len(t, folders, 3)
	assert.Equal(t, folder3, folders[0].Name)
	assert.Equal(t, folder2, folders[1].Name)
	assert.Equal(t, folder1, folders[2].Name)

	createdDesc := defaultOpt
	createdDesc.Field = SortingFieldCreatedTime
	createdDesc.Order = SortingOrderDesc
	folders, err = f.ListFolder(username, nameDesc)
	require.NoError(t, err)
	require.Len(t, folders, 3)
	assert.Equal(t, folder3, folders[0].Name)
	assert.Equal(t, folder2, folders[1].Name)
	assert.Equal(t, folder1, folders[2].Name)
}
