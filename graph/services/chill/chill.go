package chill

import (
	"chilly_daze_gateway/graph/model"
	"context"
	"database/sql"
	"log"
	"time"
)

type ChillService struct {
	Db *sql.DB
}

func (u *ChillService) StartChill(
	ctx context.Context,
	input model.StartChillInput,
) (*model.Chill, error) {
	timestamp := time.Now()
	if input.Timestamp != nil {
		timestamp, err := time.Parse(time.RFC3339, *input.Timestamp)
		if err != nil {
			log.Println("time.Parse error:", err)
			return nil, err
		}
	}

}
