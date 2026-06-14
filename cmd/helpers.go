package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Emin-07/final-project/internal/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const tablesCheckSqliteQuery = `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`

func (app *application) writeJson(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (app *application) jsonError(w http.ResponseWriter, error string, status int) {
	app.writeJson(w, map[string]string{"error": error}, status)
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	key := os.Getenv("JWT_KEY")
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(key), nil
}

func createTablesAndOpenDB(dsn string) (*sql.DB, error) {
	// dsn == db name if it's sqlite you're using
	dbName := os.Getenv("TODO_DBFILE")
	if dbName == "" {
		dbName = "scheduler.db"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Adding table into db if it's empty
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	isEmpty, err := models.HasAnyTables(ctx, db, tablesCheckSqliteQuery)
	if err != nil {
		return nil, err
	}

	if isEmpty {
		err = models.InitScheduleTable(ctx, db)
		if err != nil {
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
