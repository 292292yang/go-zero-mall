package cryptx

import "golang.org/x/crypto/bcrypt"

// 加密password
func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 检测密码
func CheckPassword(password string, encrypted string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	return err == nil
}
