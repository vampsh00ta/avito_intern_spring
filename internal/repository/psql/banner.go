package psql

import (
	"avito_intern/internal/models"
	"avito_intern/internal/repository/psql/complex_queries"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Banner interface {
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	GetBannerForUser(ctx context.Context, userTag int32, featureID int32) (models.Banner, error)
	GetBannerByID(ctx context.Context, ID int) (models.Banner, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error)
	ChangeBanner(ctx context.Context, ID int, banner models.BannerChange) error
	AddTagsToBanner(ctx context.Context, bannerID int, featureID int32, tags ...int32) error
	DeleteBannerByID(ctx context.Context, ID int) error
	DeleteBannerByTagAndFeature(ctx context.Context, featureID, tagID int32) (int, error)
}

// insert into banner (content,is_active) values ($1,$2) returning id.
func (db Pg) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	client := db.getDB(ctx)
	q := `insert into banner (content,is_active) values ($1,$2) returning id`
	var bannerID int
	if err := client.QueryRow(ctx, q, banner.Content, banner.IsActive).
		Scan(&bannerID); err != nil {
		return -1, err
	}
	if err := db.AddTagsToBanner(ctx, bannerID, banner.Feature, banner.Tags...); err != nil {
		return -1, err
	}
	return bannerID, nil
}

// insert into banner_tag (banner_id,feature_id,tag_id) values (args...) returning banner_id.
func (db Pg) AddTagsToBanner(ctx context.Context, bannerID int, featureID int32, tags ...int32) error {
	client := db.getDB(ctx)
	q := `insert into banner_tag (banner_id,feature_id,tag_id) values`
	input := []any{bannerID, featureID}
	for i, tag := range tags {
		q += fmt.Sprintf(" ($1,$2,$%d),", i+3)
		input = append(input, tag)
	}
	q = q[:len(q)-1] + " returning banner_id"
	if err := client.QueryRow(ctx, q, input...).Scan(&bannerID); err != nil {
		return err
	}
	return nil
}

// select banner.*  from
// select * from banner_tag where args.. ) banner_tag
// join banner on banner_tag.banner_id = banner.id limit/offset args...
func (db Pg) GetBannerForUser(ctx context.Context, userTag int32, featureID int32) (models.Banner, error) {
	client := db.getDB(ctx)
	q := `select banner.*  from 
        (select * from banner_tag where tag_id = $1 and  feature_id = $2) banner_tag
		join banner on banner_tag.banner_id = banner.id `
	row, err := client.Query(ctx, q, userTag, featureID)
	if err != nil {
		return models.Banner{}, err
	}
	res, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Banner])
	if err != nil {
		return models.Banner{}, err
	}

	return res, nil
}

// select banner.* , banner_tag.feature_id,tag_id from
// (select *  from  banner where id = $1) banner
// join banner_tag on banner.id = banner_tag.banner_id.
func (db Pg) GetBannerByID(ctx context.Context, ID int) (models.Banner, error) {
	client := db.getDB(ctx)
	q := `select banner.* , banner_tag.feature_id,tag_id from  
		   (select *  from  banner where id = $1) banner
		   join banner_tag on banner.id = banner_tag.banner_id`
	row, err := client.Query(ctx, q, ID)
	if err != nil {
		return models.Banner{}, err
	}
	bannerTags, err := pgx.CollectRows(row, pgx.RowToStructByName[models.BannerTags])
	if err != nil {
		return models.Banner{}, err
	}
	if len(bannerTags) == 0 {
		return models.Banner{}, pgx.ErrNoRows
	}
	var res models.Banner
	for _, bannerTag := range bannerTags {
		res.ID = bannerTag.ID
		res.IsActive = bannerTag.IsActive
		res.Content = bannerTag.Content
		res.Feature = bannerTag.Feature

		res.UpdatedAt = bannerTag.UpdatedAt
		res.CreatedAt = bannerTag.CreatedAt
		if bannerTag.Tag != nil {
			res.Tags = append(res.Tags, *bannerTag.Tag)

		}

	}
	return res, nil
}

// select banner.*,banner_tag.tag_id,banner_tag.feature_id   from banner
// join  (select distinct p1.banner_id from banner_tag p1
// join banner_tag p2  on p2.banner_id = p1.banner_id and p2.feature_id   = p1.feature_id and p2.tag_id   = $%d) banner_filter
//
//	on banner_tag.banner_id = banner.id
//
// join banner_tag on   banner_filter.banner_id = banner_tag.banner_id.
func (db Pg) GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error) {
	client := db.getDB(ctx)

	q, args := complex_queries.GetBanners(tagID, featureID, limit, offset)
	rows, err := client.Query(ctx, q, args...)
	if err != nil {

		return nil, err
	}
	rowReses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BannerTags])
	if err != nil {
		return nil, err
	}
	if len(rowReses) == 0 {
		return nil, pgx.ErrNoRows
	}

	mapping := make(map[int]*models.Banner)

	for _, rowRes := range rowReses {
		id := rowRes.ID
		curr, ok := mapping[id]
		if !ok {
			curr = &models.Banner{
				ID:        rowRes.ID,
				Feature:   rowRes.Feature,
				Content:   rowRes.Content,
				IsActive:  rowRes.IsActive,
				CreatedAt: rowRes.CreatedAt,
				UpdatedAt: rowRes.UpdatedAt,
				Tags:      make([]int32, 0),
			}
			mapping[id] = curr
		}
		if rowRes.Tag != nil {
			curr.Tags = append(curr.Tags, *rowRes.Tag)

		}
	}

	res := make([]models.Banner, len(mapping))
	i := 0
	for _, value := range mapping {
		res[i] = *value
		i++
	}
	return res, nil
}

// delete from banner where id = $1.
func (db Pg) DeleteBannerByID(ctx context.Context, ID int) error {
	client := db.getDB(ctx)
	q := `delete from banner where id = $1 returning id`
	if err := client.QueryRow(ctx, q, ID).Scan(&ID); err != nil {
		return err
	}
	return nil
}

func (db Pg) DeleteBannerByTagAndFeature(ctx context.Context, featureID, tagID int32) (int, error) {
	client := db.getDB(ctx)
	var bannerID int
	q := `delete from banner 
       	where id = (select banner_id  from banner_tag where feature_id = $1  and tag_id = $2)
        returning id`
	if err := client.QueryRow(ctx, q, featureID, tagID).Scan(&bannerID); err != nil {
		return -1, err
	}
	return bannerID, nil
}

func (db Pg) ChangeBannerTagsOrFeature(ctx context.Context, ID int, featureID *int32, tagIDs *[]int32) error {
	client := db.getDB(ctx)

	if tagIDs != nil && len(*tagIDs) > 0 {
		argsTags := []int32{}
		if tagIDs != nil {
			argsTags = append(argsTags, *tagIDs...)

		}
		var res int32

		q := `delete from banner_tag where banner_id = $1 returning feature_id`

		if err := client.QueryRow(ctx, q, ID).Scan(&res); err != nil {
			return err
		}
		if featureID == nil {
			featureID = &res
		}
		if err := db.AddTagsToBanner(ctx, ID, *featureID, argsTags...); err != nil {
			return err
		}

	} else if featureID != nil && tagIDs == nil {

		q := `update  banner_tag set feature_id = $2 where banner_id = $1 returning banner_id`
		if err := client.QueryRow(ctx, q, ID, *featureID).Scan(&ID); err != nil {
			return err
		}
	} else if tagIDs != nil && len(*tagIDs) == 0 {
		var res int32
		fmt.Println("asd")
		q := `delete from banner_tag where banner_id = $1 returning feature_id`

		if err := client.QueryRow(ctx, q, ID).Scan(&res); err != nil {
			return err
		}
		q = `insert into   banner_tag (banner_id,feature_id,tag_id) values($1,$2,null) returning banner_id`

		if err := client.QueryRow(ctx, q, ID, res).Scan(&ID); err != nil {
			return err
		}
	}

	return nil
}

func (db Pg) ChangeBanner(ctx context.Context, ID int, banner models.BannerChange) error {
	client := db.getDB(ctx)
	q := `update banner set updated_at = now()`
	argCount := 2
	args := []any{ID}
	if banner.Content != nil {
		q += fmt.Sprintf(" , content = $%d", argCount)
		args = append(args, banner.Content)
		argCount++
	}
	if banner.IsActive != nil {
		q += fmt.Sprintf(" , is_active = $%d", argCount)
		args = append(args, banner.IsActive)
		argCount++
	}

	q += " where id = $1 returning id"

	if err := client.QueryRow(ctx, q, args...).Scan(&ID); err != nil {
		return err
	}

	if err := db.ChangeBannerTagsOrFeature(ctx, ID, banner.Feature, banner.Tags); err != nil {
		return err
	}
	return nil
}
