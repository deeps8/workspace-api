package database

import (
	"work-space-backend/utils"

	"github.com/labstack/echo/v4"
)

// Insert the user details into table if user does not exsists
// and return userID in both cases
func InsertUser(c echo.Context, user utils.UserDetails) (string, error) {
	stmt, prerr := Db.Prepare(`
		WITH s AS( SELECT id FROM users WHERE client_id = $1 ), 
		ns AS( INSERT INTO users(client_id ,email ,name ,picture ,created_at) 
			SELECT $1 , $2, $3, $4,NOW()
			WHERE NOT EXISTS (
				SELECT 1 from s
			) 
			RETURNING id
		)
		select id from ns union all select id from s
	`)
	if prerr != nil {
		return "", prerr
	}

	defer stmt.Close()

	row, rerr := stmt.Query(user.Id, user.Email, user.Name, user.Picture)
	if rerr != nil {
		return "", rerr
	}

	var id string
	row.Next()
	row.Scan(&id)

	return id, nil
}
