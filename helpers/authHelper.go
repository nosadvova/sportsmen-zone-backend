package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeToId(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("user_id")
	err = nil

	if userType != "Sportsman" || userType != "Trainer" && uid != userId {
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
