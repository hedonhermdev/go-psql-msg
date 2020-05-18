package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	DRIVER  = "postgres"
	SSLMODE = "disable"
)

type DBConnInfo struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

type DB *sql.DB

func LoadFromEnv() *DBConnInfo {
	host := os.Getenv("PG_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PG_PORT")
	if port == "" {
		port = "5432"
	}
	dbName := os.Getenv("PG_DB_NAME")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	return &DBConnInfo{
		Host:     host,
		Port:     port,
		DBName:   dbName,
		User:     user,
		Password: password,
	}
}

func (info *DBConnInfo) ConnString() string {
	conn_string := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", info.Host, info.Port, info.DBName, info.User, info.Password, SSLMODE)

	return conn_string
}

func ConnectToDB(conn_string string) (DB, error) {
	db_conn, err := sql.Open(DRIVER, conn_string)
	if err != nil {
		return nil, fmt.Errorf("db: %q", err)
	}

	err = db_conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("db: database closed connection unexpectedly")
	}
	return DB(db_conn), nil
}
