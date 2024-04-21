package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFS_UserRegister(t *testing.T) {
	f := NewFS()

	// invalid chars
	user, err := f.UserRegister("user@name")
	assert.ErrorIs(t, err, ErrInvalidCharacter)

	// register user
	user, err = f.UserRegister(username)
	require.NoError(t, err)
	assert.Equal(t, user.Name, username)

	// register user again
	user, err = f.UserRegister(username)
	assert.ErrorIs(t, err, ErrDataAlreadyExists)
}
