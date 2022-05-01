package application

import (
	"errors"
	"math/rand"
	"time"

	"github.com/xyproto/randomstring"

	"github.com/hahattan/assignment/interfaces"
)

const timeFormat = time.RFC3339

func CreateToken(dbClient interfaces.DBClient) (string, error) {
	rand.Seed(time.Now().UnixNano())
	l := 6 + rand.Intn(7)
	token := randomstring.CookieFriendlyString(l)
	err := dbClient.AddToken(token, time.Now().Add(5*time.Minute).Format(timeFormat))
	if err != nil {
		return "", err
	}

	return token, nil
}

func DisableToken(token string, dbClient interfaces.DBClient) error {
	expiry := time.Time{}.Format(timeFormat)
	return dbClient.UpdateToken(token, expiry)
}

func GetAllTokens(dbClient interfaces.DBClient) (map[string]string, error) {
	res, err := dbClient.AllTokens()
	if err != nil {
		return nil, err
	}

	var status string
	var expiry time.Time
	for k, v := range res {
		expiry, err = time.Parse(timeFormat, v)
		if err != nil {
			return nil, err
		}
		if expiry.IsZero() {
			status = "DISABLED"
		} else if time.Now().After(expiry) {
			status = "EXPIRED"
		} else {
			status = "VALID"
		}
		res[k] = status
	}

	return res, nil
}

func ValidateToken(token string, dbClient interfaces.DBClient) error {
	res, err := dbClient.ExpiryByToken(token)
	if err != nil {
		return err
	}

	expiry, err := time.Parse(timeFormat, res)
	if err != nil {
		return err
	}

	if expiry.IsZero() {
		return errors.New("token disabled")
	}
	if time.Now().After(expiry) {
		return errors.New("token expired")
	}

	return nil
}
