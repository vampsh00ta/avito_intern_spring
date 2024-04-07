package psql

import "fmt"

func buildGetBannersQuery(tagID, featureID, limit, offset int32) (string, []any) {
	q := `select banner.*,banner_tag.tag_id,banner_tag.feature_id   from banner 
		join `

	argCount := 1
	args := []any{}
	if tagID != 0 {
		q += fmt.Sprintf(` (
			select distinct p1.banner_id from banner_tag p1 
			join banner_tag p2
			on p2.banner_id = p1.banner_id and p2.feature_id   = p1.feature_id and p2.tag_id   = $%d )`, argCount)
		argCount += 1
		args = append(args, tagID)

	}
	if featureID != 0 {
		if argCount == 1 {
			q += fmt.Sprintf(" (select distinct banner_id from banner_tag where feature_id = $%d)", argCount)
		} else {
			q = q[:len(q)-1]
			q += fmt.Sprintf(" feature_id = $%d) ", argCount)
		}
		argCount += 1
		args = append(args, featureID)

	}
	if featureID != 0 && tagID != 0 {
		q += fmt.Sprintf(" (select distinct banner_id from banner_tag)", argCount)
	}
	if limit > 0 {
		q = q[:len(q)-1]

		q += fmt.Sprintf(" limit  $%d)", argCount)
		argCount += 1
		args = append(args, limit)

	}
	if offset > 0 {
		q = q[:len(q)-1]

		q += fmt.Sprintf(" offset  $%d)", argCount)
		argCount += 1
		args = append(args, offset)

	}
	q += `
 banner_filter on banner_filter.banner_id = banner.id
	join banner_tag on banner_filter.banner_id = banner_tag.banner_id
`

	return q, args

}
