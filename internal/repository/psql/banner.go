package psql

import (
	"avito_intern/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Banner interface {
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	GetBannerForUser(ctx context.Context, userTag int32, featureID int32) (models.Banner, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error)

	AddTagsToBanner(ctx context.Context, bannerID int, featureID int32, tags ...int32) error
}

// insert into banner (content,is_active) values ($1,$2) returning id
func (db Pg) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	tx, err := db.getDb(ctx)
	if err != nil {
		return -1, err
	}
	q := `insert into banner (content,is_active) values ($1,$2) returning id`
	var bannerID int
	if err := tx.QueryRow(ctx, q, banner.Content, banner.IsActive).
		Scan(&bannerID); err != nil {
		return -1, err
	}
	if err := db.AddTagsToBanner(ctx, bannerID, banner.Feature, banner.Tags...); err != nil {
		return -1, err
	}
	return bannerID, nil
}

// insert into banner_tag (banner_id,feature_id,tag_id) values (args...) returning banner_id
func (db Pg) AddTagsToBanner(ctx context.Context, bannerID int, featureID int32, tags ...int32) error {
	tx, err := db.getDb(ctx)
	if err != nil {
		return err
	}
	q := `insert into banner_tag (banner_id,feature_id,tag_id) values`
	input := []any{bannerID, featureID}
	for i, tag := range tags {
		q += fmt.Sprintf(" ($1,$2,$%d),", i+3)
		input = append(input, tag)
	}
	q = q[:len(q)-1] + "returning banner_id"
	if err := tx.QueryRow(ctx, q, input...).Scan(&bannerID); err != nil {
		return err
	}
	return nil
}

// select banner.*  from
// select * from banner_tag where args.. ) banner_tag
// join banner on banner_tag.banner_id = banner.id limit/offset args...
func (db Pg) GetBannerForUser(ctx context.Context, userTag int32, featureID int32) (models.Banner, error) {

	tx, err := db.getDb(ctx)
	if err != nil {
		return models.Banner{}, err
	}
	q := `select banner.*  from 
        (select * from banner_tag where tag_id = $1 and  feature_id = $2) banner_tag
		join banner on banner_tag.banner_id = banner.id `
	row, err := tx.Query(ctx, q, userTag, featureID)
	if err != nil {
		return models.Banner{}, err
	}
	res, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Banner])
	if err != nil {
		return models.Banner{}, err
	}

	return res, nil
}

// select banner.*,banner_tag.tag_id,banner_tag.feature_id   from banner
// join  (select * from banner_tag ) banner_tag on banner_tag.banner_id = banner.id
func (db Pg) GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error) {

	tx, err := db.getDb(ctx)
	if err != nil {
		return nil, err
	}

	//q := `select banner.*,banner_tag.tag_id,banner_tag.feature_id   from banner
	//	join  (select * from where tagID = @tagID and featureID = @featureID ) banner_tag on banner_tag.banner_id = banner.id
	//	limit @limit offset @offset
	////	`

	q, args := buildGetBannersQuery(tagID, featureID, limit, offset)
	var res []models.Banner
	rows, err := tx.Query(ctx, q, args...)
	rowReses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BannerTags])
	if err != nil {
		return nil, err
	}

	var mapping map[int]*models.Banner
	mapping = make(map[int]*models.Banner)

	for _, rowRes := range rowReses {
		id := rowRes.Id
		curr, ok := mapping[id]
		if !ok {
			curr = &models.Banner{
				Id:        rowRes.Id,
				Feature:   rowRes.Feature,
				Content:   rowRes.Content,
				IsActive:  rowRes.IsActive,
				CreatedAt: rowRes.CreatedAt,
				UpdatedAt: rowRes.UpdatedAt,
				Tags:      make([]int32, 0),
			}
			mapping[id] = curr
		}
		curr.Tags = append(curr.Tags, rowRes.Tag)
	}

	for _, value := range mapping {
		res = append(res, *value)
	}
	return res, nil
}
