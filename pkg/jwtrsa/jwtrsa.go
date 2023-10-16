package jwtrsa

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type GenerateJWTInput struct {
	PrivateKey   string
	Claims       map[string]interface{}
	TimeToExpire time.Duration
}

func updatePrivateKey(signKeyStr string) error {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signKeyStr))
	if err != nil {
		return errors.New("error in pkg jwtrsa, updatePrivateKey, when jwt.ParseRSAPRivateKeyFromPEM")
	}
	privateKey = signKey
	return nil
}

func updatePublicKey(signKeyStr string) error {
	signKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(signKeyStr))
	if err != nil {
		return errors.New("error in pkg jwtrsa, updatePublicKey when jwt.ParseRSAPublicKeyFromPEM")
	}

	publicKey = signKey
	return nil
}

// GetPrivateKey returns *rsa.PrivateKey is path contains valid rsa private key, or error if path or key is invalid
func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	err := updatePrivateKey(path)
	if err != nil {
		return nil, errors.New("error in pkg jwtrsa, GetPrivateKey, when updatePrivateKey")
	}
	return privateKey, nil
}

// GetPublicKey returns *rsa.PublicKey is path contains valid rsa public key, or error if path or key is invalid
func GetPublicKey(path string) (*rsa.PublicKey, error) {
	err := updatePublicKey(path)
	if err != nil {
		return nil, errors.New("error in pkg jwtrsa, GetPublicKey, when updatePublicKey")
	}
	return publicKey, nil
}

// GenerateJWT will generate JWT and its expired time based on input that defined by GenerateJWTInput
func GenerateJWT(input GenerateJWTInput) (tokenStr string, expiresIn time.Time, err error) {
	err = updatePrivateKey(input.PrivateKey)
	if err != nil {
		return tokenStr, expiresIn, errors.New("error in pkg jwtrsa, GenerateJWT, when updatePrivateKey")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	// Set claims
	tclaims := token.Claims.(jwt.MapClaims)
	for k, v := range input.Claims {
		tclaims[k] = v
	}
	exp := time.Now().Add(input.TimeToExpire)
	tclaims["exp"] = exp.Unix()
	// Generate encoded token and send it as response.
	t, err := token.SignedString(privateKey)
	return t, exp, err
}

// ValidateJWT will validate token string based on public key.
// It will return *jwt.Token if public key and token is valid, otherwise return error.
func ValidateJWT(publicKeyPath string, tokenStr string) (token *jwt.Token, err error) {
	err = updatePublicKey(publicKeyPath)
	if err != nil {
		return token, errors.New("error in pkg jwtrsa, ValidateJWT, when updatePublicKey")
	}
	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("signing method invalid")
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.New("error in pkg jwtrsa ValidateJWT when parsing JWT")
	}
	return token, nil
}

// DecodeJWT will decode token jwt to claims.
func DecodeJWT(token *jwt.Token) (claims jwt.MapClaims) {
	return token.Claims.(jwt.MapClaims)
}
