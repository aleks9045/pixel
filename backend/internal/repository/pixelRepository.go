package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/aleks9045/pixels/internal/model"
)

type PixelRepository struct {
	conn           *sql.DB
	fieldWidth     int
	fieldHeight    int
	pixelArraySize int
}

func NewPixelRepository(conn *sql.DB, cfg config.PixelFieldConfig) *PixelRepository {
	return &PixelRepository{
		conn:           conn,
		fieldWidth:     cfg.FieldWidth,
		fieldHeight:    cfg.FieldHeight,
		pixelArraySize: cfg.FieldWidth * cfg.FieldHeight,
	}
}

func (r *PixelRepository) GetField() (model.PixelField, error) {

	rows, err := r.conn.Query("SELECT color FROM pixel ORDER BY y, x")

	defer rows.Close()
	if err != nil {
		return model.PixelField{}, err
	}
	var pixels = make([]byte, r.pixelArraySize)
	i := 0
	for rows.Next() {
		if err := rows.Scan(&pixels[i]); err != nil {
			return model.PixelField{}, err
		}
		i++
	}
	if rows.Err() != nil {
		return model.PixelField{}, err
	}

	return model.PixelField{Width: int32(r.fieldWidth), Height: int32(r.fieldHeight), PixelColors: pixels}, nil
}

func (r *PixelRepository) PutPixel(wsMessage *model.WsPixel, userId string) error {

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(`
	UPDATE users
	SET last_placed_pixel = NOW()
	WHERE user_id = $1
	  AND (
		  last_placed_pixel IS NULL
		  OR NOW() >= last_placed_pixel + time_limit * INTERVAL '1 millisecond'
	  )
	RETURNING id
`, userId).Scan(&userID)

	if err == sql.ErrNoRows {
		return errors.New("pixel time limit")
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE pixel
		 SET color = $1,
		     placed_by = $2
		 WHERE x = $3 AND y = $4`,
		wsMessage.Color,
		userID,
		wsMessage.X,
		wsMessage.Y,
	)
	if err != nil {
		return errors.New("Invalid coordinates")
	}

	_, err = tx.Exec(
		`INSERT INTO pixel_history (x, y, placed_by)
		 VALUES ($1, $2, $3)`,
		wsMessage.X,
		wsMessage.Y,
		userID,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PixelRepository) PutPixelRegion(wsMessage *model.WsPixelRegion, userId string) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int

	err = tx.QueryRow(`
	UPDATE users
	SET last_placed_pixel = NOW()
	WHERE user_id = $1
	  AND (
		  last_placed_pixel IS NULL
		  OR NOW() >= last_placed_pixel + time_limit * INTERVAL '1 millisecond'
	  )
	RETURNING id
`, userId).Scan(&userID)

	if err == sql.ErrNoRows {
		return errors.New("pixel time limit")
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE pixel
		 SET color = $1,
		     placed_by = $2
			 WHERE x BETWEEN $3 AND $4
  				AND y BETWEEN $5 AND $6`,
		wsMessage.Color,
		userID,
		wsMessage.X1,
		wsMessage.X2,
		wsMessage.Y1,
		wsMessage.Y2,
	)
	if err != nil {
		return errors.New("Invalid coordinates")
	}

	width := wsMessage.X2 - wsMessage.X1 + 1
	height := wsMessage.Y2 - wsMessage.Y1 + 1

	size := width * height

	xArr := make([]int32, size)
	yArr := make([]int32, size)
	placedByArr := make([]int, size)

	index := 0

	log.Println(userID)

	for x := wsMessage.X1; x <= wsMessage.X2; x++ {
		for y := wsMessage.Y1; y <= wsMessage.Y2; y++ {
			xArr[index] = x
			yArr[index] = y
			placedByArr[index] = userID
			index++
		}
	}
	log.Println(index, len(placedByArr))
	_, err = tx.Exec(
		`INSERT INTO pixel_history (x, y, placed_by)
		 SELECT * FROM unnest(
		 $1::smallint[],
		 $2::smallint[],
		 $3::int[]
		 )`,
		xArr,
		yArr,
		placedByArr,
	)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *PixelRepository) GetPixelAuthor(x, y int) (*model.User, error) {

	var (
		userName     string
		name         string
		timeLimit    int32
		avatarURL    string
		studentGroup sql.NullString
	)

	err := r.conn.QueryRow(
		`SELECT username, name, avatar_url, student_group, time_limit
		 FROM users u
		 JOIN pixel p ON p.placed_by = u.id
		 WHERE p.x = $1 and p.y = $2`, x, y,
	).Scan(
		&userName,
		&name,
		&avatarURL,
		&studentGroup,
		&timeLimit,
	)

	if err != nil {
		return nil, err
	}

	return &model.User{
		Username:     userName,
		Name:         name,
		AvatarUrl:    avatarURL,
		StudentGroup: studentGroup.String,
		TimeLimit:    timeLimit,
	}, nil
}
