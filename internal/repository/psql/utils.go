package psql

import "fmt"

func buildGetBannersQuery(tagID, featureID, limit, offset int32) (string, []any) {
	q := `select banner.*,banner_tag.tag_id,banner_tag.feature_id   from banner 
		join  (select * from banner_tag
		`
	//(select * from banner_tag )
	//banner_tag on banner_tag.banner_id = banner.id
	argCount := 1
	args := []any{}
	if tagID != 0 {
		if argCount == 1 {
			q += " where "
		}
		q += fmt.Sprintf("tag_id = $%d", argCount)
		argCount += 1
		args = append(args, tagID)

	}
	if featureID != 0 {
		if argCount > 1 {
			q += " and "
		}
		q += fmt.Sprintf("feature_id = $%d", argCount)
		argCount += 1
		args = append(args, featureID)

	}
	q += ") banner_tag on banner_tag.banner_id = banner.id"
	if limit > 0 {
		q += fmt.Sprintf(" limit  $%d", argCount)
		argCount += 1
		args = append(args, limit)

	}
	if offset > 0 {
		q += fmt.Sprintf(" offset  $%d", argCount)
		argCount += 1
		args = append(args, offset)

	}
	return q, args

}
