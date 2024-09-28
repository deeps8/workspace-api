package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"work-space-backend/utils"
)

func InsertWorkspace(ws utils.SpaceCreateDTO) error {

	memArr := make([]string, 0)
	for _, m := range ws.Members {
		memArr = append(memArr, fmt.Sprintf(`('%v',(SELECT id from spaceins))`, m.Id))
	}
	memQry := strings.Join(memArr, ",\n")

	query := fmt.Sprintf(`
	WITH spaceins AS (
		INSERT INTO workspace (name,overview,slug,owner,created_at)
		SELECT $1, $2, $3, $4, NOW()
		RETURNING *
	),
	memins AS (
		INSERT INTO member (user_id,space_id)
		VALUES %v 
	)
	SELECT * FROM spaceins;`, memQry)

	stmt, prerr := Db.Prepare(query)
	if prerr != nil {
		return prerr
	}

	defer stmt.Close()

	_, err := stmt.Exec(ws.Name, ws.Overview, ws.Slug, ws.Owner)
	if err != nil {
		return err
	}

	return nil
}

func GetAllWorkspace() ([]utils.SpaceGetDTO, error) {
	var spaces []utils.SpaceGetDTO
	stmt, err := Db.Prepare(`
	SELECT 
		w.*,
		o.id as user_id, o.email, o.name as user_name, o.picture,
		json_agg(u) AS members
	FROM 
		workspace w
	LEFT JOIN 
		member m ON w.id = m.space_id
	LEFT JOIN 
		users u ON m.user_id = u.id
	JOIN 
		users o ON w.owner = o.id
	GROUP BY 
		w.id, o.id
	`)
	// id	name	overview	slug	owner	created_at	updated_at
	// user_id	email	user_name	picture	members
	if err != nil {
		return spaces, err
	}
	defer stmt.Close()
	rows, qerr := stmt.Query()
	log.Printf("query run")
	if qerr != nil {
		return spaces, qerr
	}
	var space utils.SpaceGetDTO
	var mems []byte
	for rows.Next() {
		err = rows.Scan(&space.Id, &space.Name, &space.Overview, &space.Slug, &space.Owner, &space.CreatedAt, &space.UpdatedAt,
			&space.OwnerDetails.Id, &space.OwnerDetails.Email, &space.OwnerDetails.Name, &space.OwnerDetails.Picture, &mems)
		if err != nil {
			return spaces, err
		}
		json.Unmarshal(mems, &space.Members)
		spaces = append(spaces, space)
	}
	return spaces, nil
}
