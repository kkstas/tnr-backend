package testutils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kkstas/tnr-backend/internal/app"
	"github.com/kkstas/tnr-backend/internal/database"
	_ "modernc.org/sqlite"
)

func NewTestApplication(t testing.TB) (newApp http.Handler, cleanup func(), db *sql.DB) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	db, cleanupDb := OpenTestDB(t)

	cleanup = func() {
		cleanupDb()
		cancel()
	}

	getenv := func(k string) string {
		switch k {
		case "JWT_SECRET_KEY":
			return "secret-key"
		default:
			return ""
		}
	}

	newApp = app.NewApplication(ctx, getenv, db, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})))

	return newApp, cleanup, db
}

func OpenTestDB(t testing.TB) (db *sql.DB, cleanup func()) {
	t.Helper()
	dbName := fmt.Sprintf("%s.db", RandomString(32))
	db, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(1)&_time_format=sqlite")
	if err != nil {
		t.Fatalf("failed to open sql db: %v", err)
	}

	err = database.InitDBTables(context.Background(), db)
	if err != nil {
		t.Fatalf("failed to init db tables: %v", err)
	}

	cleanup = func() {
		db.Close()
		if err := os.Remove(dbName); err != nil {
			t.Fatalf("failed to remove test database file %s: %v", dbName, err)
		}
	}

	return db, cleanup
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func ToJSONBuffer(t testing.TB, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func DecodeJSON[T any](t testing.TB, body io.Reader) T {
	t.Helper()
	var result T
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	return result
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status: got %d, want %d", got, want)
	}
}

func AssertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertValidDate(t testing.TB, dateStr string) {
	t.Helper()
	layout := "2006-01-02T15:04:05Z"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		t.Errorf("string %s is not valid date in format %s: %v", dateStr, layout, err)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't expect one: %v", err)
	}
}

func AssertNotEmpty(t testing.TB, got string) {
	t.Helper()
	if len(got) == 0 {
		t.Error("expected a non-empty string but didn't get one")
	}
}
