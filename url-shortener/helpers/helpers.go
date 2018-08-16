package helpers

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	LongURL  string
	ShortURL string
}

func GetDatabaseConnector(path string) (*sql.DB, error) {
	log.Println("Connecting to database")
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return database, nil
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		pos := seededRand.Intn(len(charset))
		b[i] = charset[pos]
	}
	return string(b)
}

func CheckAndGetShortCode(longURL string, sqlitePath string, shortCodeLength int) (code string, err error) {
	database, err := GetDatabaseConnector(sqlitePath)
	if err != nil {
		return "", err
	}
	defer database.Close()

	// check if the long URL already is in the DB
	row := database.QueryRow("SELECT short FROM urls WHERE long = ?", longURL)
	var existingShortCode string
	err = row.Scan(&existingShortCode)
	if err == sql.ErrNoRows {
		// the URL is not in the DB yet - generate a new short code
		code = generateRandomString(shortCodeLength)
		log.Printf("New URL: %s -> %s", longURL, code)
		statement, _ := database.Prepare("INSERT INTO urls (long, short) VALUES (?, ?)")
		statement.Exec(longURL, code)
	} else if err != nil {
		return "", err
	} else {
		// the URL is already in the DB - serve the existing short code
		log.Printf("%s already in DB: Serving existing code %s", longURL, existingShortCode)
		code = existingShortCode
	}
	return code, nil
}

func ConstructRedirectResponse(longURL string, domain string, port int, code string) ShortenResponse {
	// omit port 80 in short URL link for convenience
	var host string
	if port != 80 {
		host = domain + ":" + strconv.Itoa(port)
	} else {
		host = domain
	}
	host = host + "/" + code

	return ShortenResponse{LongURL: longURL, ShortURL: host}
}

func GetLongURLByShortCode(sqlitePath string, code string) (string, error) {
	database, err := GetDatabaseConnector(sqlitePath)
	if err != nil {
		return "", err
	}
	defer database.Close()

	// QueryRow, because we only expect at most one result to be returned
	// NOTE: Browsers might automatically request /favicon.ico or similar
	// In that case the result will be empty and ignored
	row := database.QueryRow("SELECT long FROM urls WHERE short = ?", code)
	var url string
	err = row.Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}
