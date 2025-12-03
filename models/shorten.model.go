package models

import (
	"context"
	"time"

	"backend/config"
)

type ShortLink struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	URL       string    `json:"url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// create short link
func CreateShortLinkModel(slug string, url string, userID *int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if userID == nil {
        _, err := config.Db.Exec(
            ctx,
            "INSERT INTO short_links (slug, url) VALUES ($1, $2)",
            slug, url,
        )
        return err
    }

    _, err := config.Db.Exec(
        ctx,
        "INSERT INTO short_links (slug, url, user_id) VALUES ($1, $2, $3)",
        slug, url, *userID,
    )
    return err
}


//get url by slug
func GetURLBySlug(slug string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var url string
	err := config.Db.QueryRow(
		ctx,
		"SELECT url FROM short_links WHERE slug = $1 LIMIT 1",
		slug,
	).Scan(&url)

	return url, err
}


// getCurrent click dari db untuk disinkronkan di redis
func GetCurrentClicks(slug string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var clicks int
	err := config.Db.QueryRow(
		ctx,
		"SELECT clicks FROM short_links WHERE slug = $1",
		slug,
	).Scan(&clicks)

	return clicks, err
}


// GetAllSlugs returns all slugs
func GetAllSlugs() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.Db.Query(ctx, "SELECT slug FROM short_links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			continue
		}
		slugs = append(slugs, slug)
	}

	return slugs, rows.Err()
}