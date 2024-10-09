package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func IsPasswordCorrect(hashedPassword, givenPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword)) == nil
}
