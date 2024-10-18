package database

import "work-space-backend/utils"

func InsertBoard(brd utils.BoardDTO) (utils.BoardDTO, error) {
	var b utils.BoardDTO
	stmt, err := Db.Prepare(`
	INSERT INTO board (name,type,slug,owner,space_id,data)
	SELECT $1, $2, $3, $4, $5, $6
	RETURNING *`)

	if err != nil {
		return b, err
	}
	defer stmt.Close()

	row, qerr := stmt.Query(brd.Name, brd.Type, brd.Slug, brd.Owner, brd.SpaceId, brd.Data)

	if qerr != nil {
		return b, qerr
	}
	row.Next()
	err = row.Scan(&b.Id, &b.Name, &b.Type, &b.Slug, &b.Data, &b.Owner, &b.SpaceId)
	if err != nil {
		return b, err
	}
	return b, nil

}

func SelectBordData(slug string) (utils.BoardDTO, error) {
	var b utils.BoardDTO
	stmt, err := Db.Prepare(`SELECT * FROM board WHERE slug = $1`)
	if err != nil {
		return b, err
	}
	defer stmt.Close()

	row, qerr := stmt.Query(slug)
	if qerr != nil {
		return b, qerr
	}

	row.Next()
	err = row.Scan(&b.Id, &b.Name, &b.Type, &b.Slug, &b.Data, &b.Owner, &b.SpaceId)
	if err != nil {
		return b, err
	}
	return b, nil
}

func UpdateBoardData(slug string, data string) error {
	stmt, err := Db.Prepare(`
		UPDATE board
		SET data = $1
		WHERE slug = $2`)

	if err != nil {
		return err
	}

	_, qerr := stmt.Query(data, slug)
	if qerr != nil {
		return qerr
	}

	return nil
}
