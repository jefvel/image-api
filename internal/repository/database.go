package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image-api/internal/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ImageDetails struct {
	ID       int        `json:"id"`
	Metadata Metadata   `json:"metadata"`
	Data     *ImageData `json:"-"`
}

type Metadata struct {
	Size      int       `json:"size"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`
}

type ImageData string

var DB *sql.DB

var (
	ErrNotFound = errors.New("Image not found")
)

func init() {
	c := config.GetConfig()

	psqlConf := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DbAddress, c.DbPort, c.DbUser, c.DbPass, c.DbName,
	)

	var err error
	DB, err = sql.Open("postgres", psqlConf)

	if err != nil {
		log.Fatalln("Could not connect to database", err)
	}

	if err = migrate(); err != nil {
		log.Fatalln("Could not create table", err)
	}
}

func migrate() error {
	_, err := DB.Query(`CREATE TABLE IF NOT EXISTS images(
		id serial primary key,
		size INTEGER,
		width INTEGER,
		height INTEGER,
		format text,
		created_at TIMESTAMP DEFAULT now(),
		data text
	);`)

	return err
}

func ListImageMetadata(ctx context.Context) ([]ImageDetails, error) {
	rows, error := DB.QueryContext(ctx, `SELECT id, size, width, height, format, created_at FROM images`)
	if error != nil {
		return nil, error
	}

	defer rows.Close()

	result := make([]ImageDetails, 0)
	for rows.Next() {
		d := ImageDetails{}
		m := &d.Metadata

		err := rows.Scan(&d.ID, &m.Size, &m.Width, &m.Height, &m.Format, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetImageAndMetadata(ctx context.Context, ID int) (*ImageDetails, error) {
	d := &ImageDetails{}
	m := &d.Metadata

	row := DB.QueryRowContext(ctx, `SELECT id, size, width, height, format, created_at, data FROM images WHERE id = $1`, ID)

	err := row.Scan(&d.ID, &m.Size, &m.Width, &m.Height, &m.Format, &m.CreatedAt, &d.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return d, nil
}

func SaveImage(ctx context.Context, m Metadata, data ImageData) (*ImageDetails, error) {
	stmt, err := DB.PrepareContext(ctx, `INSERT INTO images (size, width, height, format, data) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, m.Size, m.Width, m.Height, m.Format, data)

	var id int
	var created_at time.Time
	err = row.Scan(&id, &created_at)
	if err != nil {
		return nil, err
	}

	res := &ImageDetails{
		Metadata: m,
	}

	res.ID = id
	res.Metadata.CreatedAt = created_at

	return res, nil
}

func UpdateImage(ctx context.Context, id int, m Metadata, data ImageData) (*ImageDetails, error) {
	stmt, err := DB.PrepareContext(ctx, `UPDATE images SET size = $1, width = $2, height = $3, format = $4, data = $5 WHERE id = $6`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, m.Size, m.Width, m.Height, m.Format, data, id)
	if err != nil {
		return nil, err
	}

	updatedRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if updatedRows == 0 {
		return nil, ErrNotFound
	}

	details := &ImageDetails{
		ID:       id,
		Metadata: m,
	}

	return details, nil
}
