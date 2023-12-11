package chill

import (
	"chilly_daze_gateway/graph/model"
	"context"
	"database/sql"
	"log"
	"time"
	"github.com/google/uuid"

	"chilly_daze_gateway/graph/services/lib/get_query"
)

type ChillService struct {
	Db *sql.DB
}

func (u *ChillService) StartChill(
	ctx context.Context,
	input model.StartChillInput,
) (*model.Chill, error) {
	chill := &model.Chill{
		ID: uuid.New().String(),
	}
	timestamp := time.Now()
	if input.Timestamp != nil {
		timestampRawValue, err := time.Parse(time.RFC3339, *input.Timestamp)
		if err != nil {
			log.Println("time.Parse error:", err)
			return nil, err
		}
		timestamp = timestampRawValue
	}

	query := get_query.GetInsertChillQuery(chill.ID, timestamp)

	_, err := u.Db.Exec(query)
	if err != nil {
		log.Println("u.Db.Exec error:", err)
		return nil, err
	}

	return chill, nil
}
