package services

import (
	"context"
	"log"
	"time"

	"backend/config"
	"backend/models"
)

func SyncClicksFromRedis(slug string, clicks int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := config.Db.Exec(
		ctx,
		"UPDATE short_links SET clicks = $1, updated_at = CURRENT_TIMESTAMP WHERE slug = $2",
		clicks, slug,
	)
	return err
}

func StartClicksSyncService() {
	ticker := time.NewTicker(5 * time.Minute) 
	defer ticker.Stop()

	log.Println("Click sync service started")

	for range ticker.C {
		if err := syncClicksToDatabase(); err != nil {
			log.Printf("Error syncing clicks: %v", err)
		}
	}
}

func syncClicksToDatabase() error {
	ctx := context.Background()
	
	slugs, err := models.GetAllSlugs()
	if err != nil {
		return err
	}

	syncCount := 0
	for _, slug := range slugs {
		key := "link:" + slug + ":clicks"
		clicks, err := config.Rdb.Get(ctx, key).Int()
		if err != nil {
			continue 
		}

		// sync ke database
		if err := SyncClicksFromRedis(slug, clicks); err != nil {
			log.Printf("Failed to sync slug %s: %v", slug, err)
			continue
		}
		
		syncCount++
	}

	log.Printf("Synced %d click counters to database", syncCount)
	return nil
}