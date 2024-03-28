package hashing

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 32)
	return string(hashed), err
}

func ComparePlainToHashed(plain string, hashed string) (bool, error) {
	toCompare, err := HashPassword(plain)
	if err != nil {
		return false, err
	}
	return toCompare == hashed, nil
}
