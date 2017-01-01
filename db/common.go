package db

import (
    "github.com/jackc/pgx"
)

type Result interface {

    Scan(dest ...interface{}) (err error)

    Values() ([]interface{}, error)

}

type Conn interface {

    Query(sql string, args ...interface{}) (*pgx.Rows, error)

}

type DataRequest struct {

    Start int64             `json:"start"`
    Limit int64             `json:"limit"`
    OrderBy string          `json:"orderBy"`
    OrderByAscending bool   `json:"orderByAscending"`
    Filter interface{}      `json:"filter"`

}

type DataResponse struct {

    Items *[]interface{}    `json:"items"`
    Count int64             `json:"count"`
    Aggregate interface{}   `json:"aggregate"`

}

type Transaction interface {

    Execute(sql string, args ...interface{}) int64

    Query(query string, handler func(Result), args ...interface{}) int32

}