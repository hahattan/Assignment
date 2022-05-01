package interfaces

type DBClient interface {
	AddToken(token string, expiry string) error
	UpdateToken(token string, expiry string) error
	TokenExists(token string) (bool, error)
	ExpiryByToken(token string) (string, error)
	AllTokens() (map[string]string, error)
}
