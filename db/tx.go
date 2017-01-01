package db
import (
    "github.com/jackc/pgx"
    "log"
    "time"
)

type txWrapper struct {
    conn *pgx.ConnPool
    tx * pgx.Tx
}

type resultWrapper struct {

    rows *pgx.Rows

}

func (self *txWrapper) Execute(sql string, args ...interface{}) int64 {
    self.begin()
    tag, err := self.tx.Exec(sql, args...)
    if err != nil {
	log.Printf(`[tx] %s "%s"`, "Exec", sql)
        log.Printf(`[tx] %s "%s"`, "ERROR", err)
        panic(err)
    }
    return tag.RowsAffected()
}

func (self *txWrapper) begin() error {
    if self.tx == nil {
        var err error
        self.tx, err = self.conn.Begin()
        if err != nil {
            log.Printf("[tx] %s %v\n", "START FAILED", time.Now())
            self.Rollback()
            return err
        }
        log.Printf("[tx] %s %v\n", "STARTED", time.Now())
    }
    return nil;
}

func (self *txWrapper) Query(query string, handler func(Result), args ...interface{}) int32 {
    self.begin()
    var rows *pgx.Rows
    rows, err := self.tx.Query(query, args...)
    if err != nil {
	log.Printf(`[tx] %s "%s"`, "Query", query)
        log.Printf(`[tx] %s "%s"`, "ERROR", err)
        panic(err)
    }
    defer func() {
        err := recover()
        rows.Close()
        if err != nil {
            log.Printf("Query: %s", err)
            panic(err)
        }
    }()
    var c int32
    var r resultWrapper
    r.rows = rows
    for rows.Next() {
        handler(&r)
        c++
    }
    return c
}

func (self *resultWrapper) Scan(dest ...interface{}) (err error) {
    e := self.rows.Scan(dest...)
    if e != nil {
        self.rows.Values()
        panic(e)
    }
    return e
}

func (self *resultWrapper) Values() ([]interface{}, error) {
    return self.rows.Values()
}

func (self *txWrapper) Select(query string, handler func(Result), args ...interface{}) (int32, error) {
    var rows *pgx.Rows
    rows, err := self.conn.Query(query, args...)
    if err != nil {
        panic(err)
    }
    defer func() {
        err := recover()
        rows.Close()
        if err != nil {
            panic(err)
        }
    }()
    var c int32
    var r * resultWrapper
    r.rows = rows
    for rows.Next() {
        handler(r)
        c++
    }
    return c, nil
}

func (self *txWrapper) Commit() error {
    if self.tx != nil {
        log.Printf("[tx] %s %v\n", "COMMITTING", time.Now())
        tx := self.tx
        self.tx = nil
        err := tx.Commit()
        if err != nil {
            tx.Rollback()
            return err
        }
        log.Printf("[tx] %s %v\n", "COMMITTED", time.Now())
    }
    return nil
}

func (self *txWrapper) Rollback() {
    if self.tx != nil {
        defer log.Printf("[tx] %s %v\n", "ROLLBACK", time.Now())
        tx := self.tx
        self.tx = nil
        tx.Rollback()
    }
}