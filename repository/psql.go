package repository

import (
	"database/sql"
	"fmt"
	"log"
	"stratum-server/config"

	_ "github.com/lib/pq"
)

const (
	postgresDriver = "postgres"
)

type request struct {
	query string
	args  []interface{}
}

type QueryRequest struct {
	Query string
	Args  []interface{}
}

type InsertRequest struct {
	Query string
	Args  []interface{}
}

type UpdateRequest struct {
	Query string
	Args  []interface{}
}

// Repository describes interface to deal with repository.
type Repository interface {
	Query(input QueryRequest, destinationArgs ...interface{}) error
	Insert(input InsertRequest, destinationArgs ...interface{}) error
	Update(input UpdateRequest, destinationArgs ...interface{}) error
}

type postgres struct {
	db *sql.DB
}

// NewService creates new instance for devices service.
func NewRepository(cfg config.PostgreSQLConfig) *postgres {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DB)

	db, err := sql.Open(postgresDriver, psqlInfo)
	if err != nil {
		log.Panicf("error establishing connection with postgres: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicf("error executing ping with postgres: %v", err)
	}

	return &postgres{
		db: db,
	}
}

func (psql *postgres) Query(input QueryRequest, destinationArgs ...interface{}) error {
	req := request {
		query: input.Query,
		args: input.Args,
	}
	if err := psql.queryRow(req, destinationArgs...); err != nil {
		log.Printf("error performing Query: %v", err)
		return err
	}

	return nil
}

func (psql *postgres) Insert(input InsertRequest, destinationArgs ...interface{}) error {
	req := request {
		query: input.Query,
		args: input.Args,
	}
	if err := psql.queryRow(req, destinationArgs...); err != nil {
		log.Printf("error performing Insert: %v", err)
		return err
	}

	return nil
}

func (psql *postgres) Update(input UpdateRequest, destinationArgs ...interface{}) error {
	req := request {
		query: input.Query,
		args: input.Args,
	}
	if err := psql.queryRow(req, destinationArgs...); err != nil {
		log.Printf("error performing Update: %v", err)
		return err
	}

	return nil
}

func (psql *postgres) queryRow(req request, destinationArgs ...interface{}) error {
	if err := psql.db.QueryRow(req.query, req.args...).Scan(destinationArgs...); err != nil {
		log.Printf("error performing queryRow: %v", err)
		return err
	}
	return nil
}
