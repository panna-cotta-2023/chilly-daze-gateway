package get_query

import (
	"fmt"
	"time"
)

func GetInsertChillQuery(
	timestamp time.Time,
) string {
	query := fmt.Sprintf(
		`INSERT INTO chill (
			created_id
	)
	VALUES (
		%s
	)`,
		timestamp,
	)

	return query
}
