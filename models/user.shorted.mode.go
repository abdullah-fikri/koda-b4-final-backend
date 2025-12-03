package models

import (
	"context"
	"time"

	"backend/config"
)


type DailyAnalytic struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Date        string    `json:"date"`
	TotalLinks  int       `json:"total_links"`
	TotalVisits int       `json:"total_visits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DashboardStats struct {
	TotalLinks      int              `json:"total_links"`
	LinksThisWeek   int              `json:"links_this_week"`
	TotalVisits     int              `json:"total_visits"`
	VisitsGrowth    float64          `json:"visits_growth"`
	AvgClickRate    float64          `json:"avg_click_rate"`
	ClickRateChange float64          `json:"click_rate_change"`
	Last7Days       []VisitorChart   `json:"last_7_days"`
}

type VisitorChart struct {
	Date   string `json:"date"`
	Visits int    `json:"visits"`
}


// get all link user
func GetUserLinksPaginated(userID int, page int, limit int) ([]ShortLink, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset := (page - 1) * limit

	// get total count
	var totalCount int
	err := config.Db.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM short_links WHERE user_id = $1`,
		userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// fetch paginated data
	rows, err := config.Db.Query(
		ctx,
		`SELECT id, slug, url, clicks, created_at, updated_at
		 FROM short_links
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var links []ShortLink
	for rows.Next() {
		var link ShortLink
		if err := rows.Scan(
			&link.ID,
			&link.Slug,
			&link.URL,
			&link.Clicks,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		links = append(links, link)
	}

	return links, totalCount, nil
}


// detail link by slug
func GetUserLinkBySlug(userID int, slug string) (*ShortLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var link ShortLink
	err := config.Db.QueryRow(
		ctx,
		`SELECT id, slug, url, clicks, created_at, updated_at 
		 FROM short_links 
		 WHERE slug = $1 AND user_id = $2`,
		slug, userID,
	).Scan(
		&link.ID,
		&link.Slug,
		&link.URL,
		&link.Clicks,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &link, nil
}

// update 
func UpdateUserLink(userID int, slug string, newURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := config.Db.Exec(
		ctx,
		`UPDATE short_links 
		 SET url = $1, updated_at = CURRENT_TIMESTAMP 
		 WHERE slug = $2 AND user_id = $3`,
		newURL, slug, userID,
	)
	if err != nil {
		return err
	}

	// cek jika gada
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// delete 
func DeleteUserLink(userID int, slug string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := config.Db.Exec(
		ctx,
		`DELETE FROM short_links 
		 WHERE slug = $1 AND user_id = $2`,
		slug, userID,
	)
	if err != nil {
		return err
	}

	// cek jika gada
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// custom error for not found
var ErrNotFound = &NotFoundError{}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "short link not found or you don't have permission"
}


// dashboard
func GetUserDashboardStats(userID int) (*DashboardStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := &DashboardStats{}

	// total links
	err := config.Db.QueryRow(ctx,
		"SELECT COUNT(*) FROM short_links WHERE user_id = $1",
		userID,
	).Scan(&stats.TotalLinks)
	if err != nil {
		return nil, err
	}

	// link dalam 7 hari
	err = config.Db.QueryRow(ctx,
		`SELECT COUNT(*) FROM short_links 
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '7 days'`,
		userID,
	).Scan(&stats.LinksThisWeek)
	if err != nil {
		return nil, err
	}

	// total visit
	err = config.Db.QueryRow(ctx,
		"SELECT COALESCE(SUM(clicks), 0) FROM short_links WHERE user_id = $1",
		userID,
	).Scan(&stats.TotalVisits)
	if err != nil {
		return nil, err
	}

	// visit last week
	var visitsLastWeek int
	err = config.Db.QueryRow(ctx,
		`SELECT COALESCE(SUM(clicks), 0) FROM short_links 
		 WHERE user_id = $1 AND updated_at < NOW() - INTERVAL '7 days'`,
		userID,
	).Scan(&visitsLastWeek)
	if err != nil {
		visitsLastWeek = 0
	}

	// kalkulasi perhitungan kuncjungan
	if visitsLastWeek > 0 {
		stats.VisitsGrowth = ((float64(stats.TotalVisits) - float64(visitsLastWeek)) / float64(visitsLastWeek)) * 100
	} else if stats.TotalVisits > 0 {
		stats.VisitsGrowth = 100.0
	}

	//  average click rate (total visits / total links)
	if stats.TotalLinks > 0 {
		stats.AvgClickRate = float64(stats.TotalVisits) / float64(stats.TotalLinks)
	}

	//  click rate from last week for comparison
	var linksLastWeek int
	err = config.Db.QueryRow(ctx,
		`SELECT COUNT(*) FROM short_links 
		 WHERE user_id = $1 AND created_at < NOW() - INTERVAL '7 days'`,
		userID,
	).Scan(&linksLastWeek)
	if err != nil {
		linksLastWeek = 0
	}

	var avgClickRateLastWeek float64
	if linksLastWeek > 0 {
		avgClickRateLastWeek = float64(visitsLastWeek) / float64(linksLastWeek)
	}

	// kalkulasi click rate persentase
	if avgClickRateLastWeek > 0 {
		stats.ClickRateChange = ((stats.AvgClickRate - avgClickRateLastWeek) / avgClickRateLastWeek) * 100
	} else if stats.AvgClickRate > 0 {
		stats.ClickRateChange = 100.0
	}

	// 7 days data 
	rows, err := config.Db.Query(ctx,
		`SELECT 
			DATE(created_at) as date,
			COUNT(*) as visits
		 FROM short_links
		 WHERE user_id = $1 
		   AND created_at >= NOW() - INTERVAL '7 days'
		 GROUP BY DATE(created_at)
		 ORDER BY date ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	
	chartData := make(map[string]int)
	for rows.Next() {
		var date time.Time
		var visits int
		if err := rows.Scan(&date, &visits); err != nil {
			continue
		}
		chartData[date.Format("2006-01-02")] = visits
	}

	// tetep buat data meski 0
	stats.Last7Days = make([]VisitorChart, 7)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		stats.Last7Days[6-i] = VisitorChart{
			Date:   dateStr,
			Visits: chartData[dateStr],
		}
	}

	return stats, nil
}

// saveDailyAnlitics db
func SaveDailyAnalytics(userID int, date string, totalLinks, totalVisits int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := config.Db.Exec(ctx,
		`INSERT INTO daily_analytics (user_id, date, total_links, total_visits)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, date) 
		 DO UPDATE SET 
		   total_links = EXCLUDED.total_links,
		   total_visits = EXCLUDED.total_visits,
		   updated_at = NOW()`,
		userID, date, totalLinks, totalVisits,
	)
	return err
}