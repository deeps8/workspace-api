package handlers

import (
	"net/http"
	"os"
	"work-space-backend/handlers/auth"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitHandler(app *echo.Group) {
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Workspace Backend")
	})

	// capturing the google config from env
	auth.GoogleConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
	}

	app.GET("/session", auth.GetUserFromSession)

	authGrp := app.Group("/auth")
	{
		/*	Creating routes for
			- signup : which will redirect user to Google Ouath page
			- google/redirect : which will handle the redirection of user from Oauth page
		*/
		authGrp.GET("/google/signup", auth.GoogleSignup)
		authGrp.GET("/google/redirect", auth.GoogleRedirect)
	}
}
