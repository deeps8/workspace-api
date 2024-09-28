package auth

import (
	"encoding/json"
	"net/http"
	"work-space-backend/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func AuthGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		/*
			1. Fetch the session data from store
			2. check for validity of token
			3. if ( valid ) -> set the user_details to ctx
			4. else ->
				a) clear the session values
				b) return the error
		*/
		token_details, _, sess := GetSessionData(c)

		if !token_details.Valid() {
			sess.Options.MaxAge = -1
			return c.JSON(http.StatusUnauthorized, utils.Res{Message: "User is not authorized", Ok: false})
		}

		return next(c)
	}
}

func UnAuthGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		/*
			1. Fetch the session data from store
			2. check for validity of token
			3. if ( valid ) ->
				a) set the user_details to ctx
		*/

		token_details, _, _ := GetSessionData(c)

		if token_details.Valid() {
			return c.JSON(http.StatusBadRequest, utils.Res{Message: "User already authenticated", Ok: false})
		}

		return next(c)
	}
}

func GetSessionData(c echo.Context) (utils.TokenDetail, utils.UserDetails, *sessions.Session) {
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
