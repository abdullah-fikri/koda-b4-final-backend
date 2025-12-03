package utils

import (
	"backend/config"
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


const (
	slugLength    = 7
	maxRetries    = 5
	
)
// generate slug
func generateSlug(length int) (string, error) {
	letterBytes   := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	max := big.NewInt(int64(len(letterBytes)))

	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[n.Int64()]
	}
	return string(b), nil
}


//cek slug di db
func slugExists(slug string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := config.Db.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM short_links WHERE slug = $1)",
		slug,
	).Scan(&exists)

	return exists, err
}


// generate unik slug
func GenerateUniqueSlug() (string, error) {
	for range maxRetries {
		slug, err := generateSlug(slugLength)
		if err != nil {
			return "", err
		}

		exists, err := slugExists(slug)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}
	}
	return "", gin.Error{Err: http.ErrHandlerTimeout, Type: gin.ErrorTypePrivate}
}