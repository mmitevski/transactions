package db
import (
	"testing"
)

func connect() db.Database {
	config := &db.DatabaseConfig{
		Host: "localhost",
		Port: 5432,
		User: "test",
		Password: "test",
		Database: "test",
	}
	return db.NewDatabase(config)
}

func TestConnectAndReadValue(t *testing.T) {
	connect().Execute(func (tx db.Transaction){
		tx.Query("select 1", func(r db.Result) {
			values, err := r.Values()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", values[0])
		})
	})
}

func TestCreateTableCRUDandDropTable(t *testing.T) {
	t.Log("Connecting to database...")
	database := connect()
	t.Log("Creating table with initial data......")
	database.Execute(func (tx db.Transaction){
		tx.Execute("drop table if exists test")
		tx.Execute("create table test(id text)")
		tx.Execute("insert into test(id) values ($1)", "1")
		tx.Execute("insert into test(id) values ($1)", "2")
		tx.Execute("insert into test(id) values ($1)", "3")
	})
	ids := ""
	database.Execute(func (tx db.Transaction) {
		tx.Query("select id from test order by id", func (r db.Result) {
			var s string
			r.Scan(&s)
			ids += s
		})
	})
	t.Log("Checking returned datas validity...")
	if ids != "123" {
		t.Fatal("Expected data to be '123', but '" + ids + "' found.")
	}
	database.Execute(func (tx db.Transaction){
		tx.Execute("drop table test")
	})
}

