package auth

import (
	"encoding/json"
	"log"
	"work-space-backend/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func GetSessionData(c echo.Context) (utils.TokenDetail, utils.UserDetails, *sessions.Session) {
	log.Printf("cookies%+v", c.Cookies())
	sess, _ := Store.Get(c.Request(), "session")

	var token_details utils.TokenDetail
	var user_details utils.UserDetails
	sess_token_details, _ := json.Marshal(sess.Values["tokenDetails"])
	sess_user_details, _ := json.Marshal(sess.Values["userDetails"])
	json.Unmarshal(sess_token_details, &token_details)
	json.Unmarshal(sess_user_details, &user_details)

	if token_details.Valid() {
		c.Set("user_details", user_details)
	} else {
		c.Set("user_details", nil)
	}

	return token_details, user_details, sess

}
