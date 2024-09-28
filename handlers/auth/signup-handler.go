package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

const GoogleOauthURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// this will store the google config globally
var GoogleConfig *oauth2.Config

func GoogleSignup(c echo.Context) error {
	// get the oauthURL from goole-config and redirect user to it
	oAuthURL := GoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, oAuthURL)
}

func GoogleRedirect(c echo.Context) error {
	/*
		1. get the authorization code from the Request URL.
		2. converts an authorization code into a token.
		3. create google-http-client and get the response.
		4. convert the above response into object/struct (user-details).
		5. store that user-details in DB.
		6. create session and store data.
	*/
	// fetch the auth-code
	authCode := c.FormValue("code")

	// convert the auth-code into token
	token, terr := GoogleConfig.Exchange(context.Background(), authCode)
	if terr != nil {
		return c.JSON(http.StatusInternalServerError, utils.Res{Message: "Something went wrong while fetching token"})
	}
	if token == nil {
		return c.JSON(http.StatusInternalServerError, utils.Res{Message: "No Token"})
	}

	// create google-http-client and get the response
	glClient := GoogleConfig.Client(context.Background(), token)
	glRes, glErr := glClient.Get(GoogleOauthURL + token.AccessToken)
	if glErr != nil {
		return c.JSON(http.StatusInternalServerError, utils.Res{Message: "Something went wrong while fetching user details"})
	}
	defer glRes.Body.Close()

	// convert the above response into struct
	var userDetails utils.UserDetails
	derr := json.NewDecoder(glRes.Body).Decode(&userDetails)
	if derr != nil {
		return c.JSON(http.StatusAccepted, utils.Res{Message: "JSON decode error"})
	}

	// store user details in DB
	userid, iserr := database.InsertUser(c, userDetails)
	if iserr != nil {
		return c.JSON(http.StatusInternalServerError, utils.Res{Message: "Error in storing user details"})
	}
	userDetails.Id = userid

	// create/get session
	sess, _ := Store.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(time.Until(token.Expiry)),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   os.Getenv("DOMAIN"),
	}

	tokenDetails := utils.TokenDetail{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
		TokenType:    token.TokenType,
	}

	gob.Register(userDetails)
	gob.Register(tokenDetails)

	sess.Values["tokenDetails"] = tokenDetails
	sess.Values["userDetails"] = userDetails

	sOk := sess.Save(c.Request(), c.Response().Writer)
	if sOk != nil {
		return c.JSON(http.StatusInternalServerError, utils.Res{Message: "Error in session creation"})
	}

	// redirect the user to frontend.
	RedirectURL := os.Getenv("FRONTEND_URL")
	return c.Redirect(http.StatusTemporaryRedirect, RedirectURL)
}

func GetUserFromSession(c echo.Context) error {
	_, user_details, _ := GetSessionData(c)
	if user_details.Id == "" {
		return c.JSON(http.StatusUnauthorized, utils.Res{Message: "User is not authorized"})
	}
	return c.JSON(http.StatusOK, user_details)
}

func LogoutUser(c echo.Context) error {
	sess, _ := Store.Get(c.Request(), "session")
	sess.Values["tokenDetails"] = nil
	sess.Values["userDetails"] = nil
	sess.Options.MaxAge = -1

	sess.Save(c.Request(), c.Response().Writer)

	c.Set("userDetails", nil)
	return c.JSON(http.StatusOK, nil)
}
