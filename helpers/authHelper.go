package helpers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func VerifyPassword(password string, hash string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintln("email or password is incorrect")
		check = false
	}

	return check, msg
}

func MatchUserTypeToId(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("user_id")
	err = nil

	if uid != userId {
		err = errors.New("invalid user type")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("invalid user type")
		return err
	}

	return err
}
