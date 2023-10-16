package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (hashedPassword string, err error) {
	bytePassword := []byte(password)
	hashedPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	hashedPassword = string(hashedPasswordByte)
	return
}

func ComparePassword(password string, hashedPassword string) bool {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(hashedPassword)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword) == nil
}
