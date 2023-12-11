package get_query

import (
	"fmt"
	"time"
)

func GetInsertChillQuery(
	id string,
	timestamp time.Time,
) string {
	query := fmt.Sprintf(
		`INSERT INTO chill (
			id,
			created_id
	)
	VALUES (
		%s,
		%s
	)`,
		id,
		timestamp,
	)

	return query
}
