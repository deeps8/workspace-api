package cronjob

import (
	"context"
	"encoding/json"
	"log"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/robfig/cron"
)

func SyncDb() {
	c := cron.New()

	// define the cron job function
	c.AddFunc("0 0 * * *", func() {
		GetAllRedisData()
	})

	c.Start()
}

func GetAllRedisData() {
	var cursor uint64
	var keys []string
	var err error

	for {

		keys, cursor, err = database.Rdb.Scan(context.Background(), cursor, "*", 10).Result()
		if err != nil {
			log.Fatalf("Error scanning keys: %v", err)
		}

		// Loop through the fetched keys and get their values
		for _, key := range keys {
			val, err := database.Rdb.Get(context.Background(), key).Result()
			if err != nil {
				log.Fatalf("Error getting value for key %s: %v", key, err)
			} else {
				var brdData utils.RdbDataType
				err := json.Unmarshal([]byte(val), &brdData)
				if err != nil {
					log.Fatalf("Error unmarshaling JSON from Redis: %v", err)
				}

				// populate the db if sync false
				if !brdData.Synced {
					log.Default().Printf("Updating the BOARD data for %v", key)
					err = database.UpdateBoardData(key, brdData.Data)
					if err != nil {
						log.Fatalf("Error updating the DB: %v", err)
					}

					// update the redis sync to true;
					brdData.Synced = true
					brdJson, err := json.Marshal(brdData)
					if err != nil {
						log.Fatalf("Error marshaling struct to JSON: %v", err)
					}
					err = database.Rdb.Set(context.Background(), key, string(brdJson), 0).Err()
					if err != nil {
						log.Fatalf("Error setting value in Redis: %v", err)
					}
				}
			}
		}

		// If cursor is 0, we've finished iterating
		if cursor == 0 {
			break
		}
	}

}
