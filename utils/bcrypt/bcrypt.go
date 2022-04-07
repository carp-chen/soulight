package bcrypt

import "golang.org/x/crypto/bcrypt"

const (
	PassWordCost = 12 //密码加密难度
)

//SetPassword 设置密码
func SetPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return "", err
	}
	PasswordDigest := string(bytes)
	return PasswordDigest, nil
}

//CheckPassword 校验密码
func CheckPassword(password string, PasswordDigest string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(PasswordDigest), []byte(password))
	return err == nil
}
