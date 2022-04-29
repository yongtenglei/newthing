package scryptx

import (
	"fmt"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"

	"golang.org/x/crypto/scrypt"
)

func PasswordEncrypt(password string) string {
	dk, err := scrypt.Key([]byte(password),
		[]byte(settings.UserServiceConf.ScryptConf.Salt),
		32768, 8, 1, 32)
	if err != nil {
		zap.S().Errorw("PasswordEncrypt scrypt.Key failed",
			"err", err.Error())
		fmt.Println("===========", err)
	}
	return fmt.Sprintf("%x", string(dk))
}

func PasswordValidate(password, expectPassword string) bool {
	scryptedPassword := PasswordEncrypt(password)

	if scryptedPassword == expectPassword {
		return true
	}

	return false
}
