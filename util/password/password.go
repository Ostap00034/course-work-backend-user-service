package password

import (
    "golang.org/x/crypto/bcrypt"
)

// Hash принимает plain-text и возвращает bcrypt-хеш.
func Hash(pass string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    return string(b), err
}

// Compare сверяет plain-text и хеш, возвращает nil, если совпадают.
func Compare(hash, pass string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
