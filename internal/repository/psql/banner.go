package psql

import (
	"avito_intern/internal/models"
	"context"
	"fmt"
)

type Banner interface {
}

func (db Pg) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	tx, err := db.getTx(ctx)
	if err != nil {
		return -1, err
	}
	q := `insert into banner (feature_id,content,is_active) values ($1,$2,$3)`

	var id int
	if err := tx.QueryRow(ctx, q, banner.Feature, banner.Content, banner.IsActive).
		Scan(&id); err != nil {
		return -1, nil
	}
	if err := db.AddTagsToBanner(ctx, banner.Id, banner.Tags...); err != nil {
		return -1, nil
	}
	return id, nil
}

func (db Pg) AddTagsToBanner(ctx context.Context, banner_id int, tags ...int32) error {
	tx, err := db.getTx(ctx)
	if err != nil {
		return err
	}
	q := `insert into banner_tags (banner_id,tag_id) values`
	input := []any{}
	for i, tag := range tags {
		q += fmt.Sprintf(" ($1,$%d),", i+2)
		input = append(input, tag)
	}
	q = q[:len(q)-1]
	if err := tx.QueryRow(ctx, q, input).Scan(&banner_id); err != nil {
		return err
	}
	return nil
}
