package psql

import (
	"avito_intern/internal/models"
	"context"
	"fmt"
)

type Banner interface {
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	AddTagsToBanner(ctx context.Context, banner_id int, tags ...int32) error
}

func (db Pg) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	tx, err := db.getDb(ctx)
	fmt.Println(tx)
	if err != nil {
		return -1, err
	}
	q := `insert into banner (feature_id,content,is_active) values ($1,$2,$3) returning id`

	var bannerID int
	if err := tx.QueryRow(ctx, q, banner.Feature, banner.Content, banner.IsActive).
		Scan(&bannerID); err != nil {
		return -1, err
	}
	if err := db.AddTagsToBanner(ctx, bannerID, banner.Tags...); err != nil {
		return -1, err
	}
	return bannerID, nil
}

func (db Pg) AddTagsToBanner(ctx context.Context, banner_id int, tags ...int32) error {
	tx, err := db.getDb(ctx)
	if err != nil {
		return err
	}
	q := `insert into banner_tag (banner_id,tag_id) values`
	input := []any{banner_id}
	for i, tag := range tags {
		q += fmt.Sprintf(" ($1,$%d),", i+2)
		input = append(input, tag)
	}
	q = q[:len(q)-1] + "returning banner_id"
	if err := tx.QueryRow(ctx, q, input...).Scan(&banner_id); err != nil {
		return err
	}
	return nil
}
