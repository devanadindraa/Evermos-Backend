package user

import "golang.org/x/crypto/bcrypt"

func comparePassword(storedHash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain))
	return err == nil
}

func hashPassword(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
