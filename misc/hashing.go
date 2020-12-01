package misc

import "golang.org/x/crypto/bcrypt"

// EncryptString encrypts the string
func EncryptString(str string) (string, error) {
	hBytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hBytes), nil
}

// CompareHash asserts the correctness of the encrypted argument
func CompareHash(hashed string, raw string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashed), 
		[]byte(raw),
		); err != nil {
			return false, err
		}
	
	return true, nil
}
