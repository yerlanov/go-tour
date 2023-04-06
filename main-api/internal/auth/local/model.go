package local

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID         string `json:"sub"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type UserCreateDto struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
