package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/udonetsm/help/helper"
	"github.com/udonetsm/help/models"
	mod "github.com/udonetsm/publisher/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connStr(path string) string {
	conn_data := models.Postgres_conf{}
	conn_data = conn_data.StoreConf(path)
	return conn_data.Dbname + conn_data.Dbpassword + conn_data.SslMode +
		conn_data.Dbhost + conn_data.Dbport + conn_data.Dbuser
}

func sqlDb() *sql.DB {
	defer helper.PanicCapture("sqlDb")
	path := helper.Home() + "/.confs/conn_config.yaml"
	connstr := connStr(path)
	sdb, err := sql.Open("pgx", connstr)
	helper.Errors(err, "sqlopen")
	sdb.SetMaxIdleConns(5)
	sdb.SetMaxOpenConns(6)
	sdb.SetConnMaxIdleTime(2)
	sdb.SetConnMaxLifetime(3)
	return sdb
}

func gormDb() *gorm.DB {
	sdb := sqlDb()
	defer helper.PanicCapture("gormDb")
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sdb}), &gorm.Config{})
	helper.Errors(err, "gormOpen")
	return gdb
}

func New(data []byte, timestamp int64) {
	sdb := sqlDb()
	order := mod.Order{}
	if err := json.Unmarshal(data, &order); err != nil {
		return
	}
	_, err := sdb.Query("insert into orders(id, orderjson, pubdate) values($1, $2, $3)", order.Order_id, string(data), timestamp)
	if err != nil {
		_, err = sdb.Query("update orders set orderjson=$1, pubdate=$2 where id=$3", string(data), timestamp, order.Order_id)
		if err != nil {
			return
		}
	}
}

func GetMaxValue() (maxpubdate int64) {
	sdb := sqlDb()
	err := sdb.QueryRow("select max(pubdate) from orders").Scan(&maxpubdate)
	log.Println(err)
	helper.Errors(err, "GetMax")
	return
}
