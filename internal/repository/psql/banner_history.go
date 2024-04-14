package psql

import (
	"avito_intern/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type BannerHistory interface {
	CreateHistoryBanner(ctx context.Context, banner models.Banner) error
	GetBannerWithHistory(ctx context.Context, bannerID, limit int) ([]models.Banner, error)
	CleanBannerHistory(ctx context.Context, limit int) error
}

// insert into history_banner (banner_id,content,is_active) values ($1,$2,$3) returning id.
func (db Pg) CreateHistoryBanner(ctx context.Context, banner models.Banner) error {
	client := db.getDB(ctx)
	q := `insert into banner_history (banner_id,content,is_active,created_at,updated_at) values ($1,$2,$3,$4,$5) returning id`
	var bannerHistoryID int
	if err := client.QueryRow(ctx, q,
		banner.ID,
		banner.Content,
		banner.IsActive,
		banner.CreatedAt,
		banner.UpdatedAt).
		Scan(&bannerHistoryID); err != nil {
		return err
	}
	fmt.Println(banner)
	if len(banner.Tags) > 0 {
		if err := db.AddTagsToHistoryBanner(ctx, bannerHistoryID, banner.Feature, banner.Tags...); err != nil {
			return err
		}
	} else {
		q := `insert into banner_tag_history (banner_history_id,feature_id,tag_id) values ($1,$2,null) returning banner_history_id`
		if err := client.QueryRow(ctx, q, bannerHistoryID, banner.Feature).Scan(&banner.ID); err != nil {
			return err
		}
	}

	return nil
}

// insert into history_banner_tag (banner_id,feature_id,tag_id) values (args...) returning banner_id.
func (db Pg) AddTagsToHistoryBanner(ctx context.Context, bannerHistoryID int, featureID int32, tags ...int32) error {
	client := db.getDB(ctx)
	q := `insert into banner_tag_history (banner_history_id,feature_id,tag_id) values `
	input := []any{bannerHistoryID, featureID}
	for i, tag := range tags {
		q += fmt.Sprintf(" ($1,$2,$%d),", i+3)
		input = append(input, tag)
	}
	q = q[:len(q)-1] + " returning banner_history_id"

	if err := client.QueryRow(ctx, q, input...).Scan(&bannerHistoryID); err != nil {
		return err
	}
	return nil
}

// insert into history_banner_tag (banner_id,feature_id,tag_id) values (args...) returning banner_id.
func (db Pg) GetBannerWithHistory(ctx context.Context, bannerID, limit int) ([]models.Banner, error) {
	client := db.getDB(ctx)

	q := `
	
	 select 
		 bh.banner_id as id, bh.id as banner_history_id,
		 bh.content ,
		 bh.is_active ,
		 bh.created_at ,
		 bh.updated_at ,
		 bth.tag_id,
		 bth.feature_id   from 
	(select * from banner_history where banner_history.banner_id = $1  order by updated_at desc limit  $2) bh 
	join banner_tag_history as bth on  bth.banner_history_id = bh.id order  by updated_at desc`
	rows, err := client.Query(ctx, q, bannerID, limit)
	if err != nil {
		return nil, err
	}
	rowReses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BannerHistoryTags])
	if err != nil {
		return nil, err
	}
	if len(rowReses) == 0 {
		return nil, pgx.ErrNoRows
	}

	mapping := make(map[int]*models.Banner)

	for _, rowRes := range rowReses {
		id := rowRes.BannerHistoryID
		curr, ok := mapping[id]
		if !ok {
			curr = &models.Banner{
				ID:        bannerID,
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

func (db Pg) CleanBannerHistory(ctx context.Context, limit int) error {
	client := db.getDB(ctx)
	q := `
		delete from banner_history where id in (
		    select id from  
              (select 
                   id, 
                   row_number() over (partition by banner_id  order by updated_at desc) banner_count 
               from banner_history) bh 
          where bh.banner_count > $1
		) returning id
	`
	if err := client.QueryRow(ctx, q, limit).Scan(&limit); err != nil {
		return err
	}
	return nil
}
