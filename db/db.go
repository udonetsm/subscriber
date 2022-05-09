package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	mod "subscriber/models"

	"github.com/udonetsm/help/helper"
	"github.com/udonetsm/help/models"
	_ "gorm.io/driver/postgres"
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

//Deletes all target_simbols from target_string
func Replace(target_string string, target_simbols ...interface{}) string {
	for _, target_simbol := range target_simbols {
		target_string = strings.Replace(target_string, target_simbol.(string), "", -1)
	}
	return target_string
}

func Set(data []byte, timestamp int64) (string, error) {
	order := mod.Order{}
	if err := json.Unmarshal(data, &order); err != nil {
		//if incoming data not order json
		return "", errors.New("data_not_valid")
	}
	d := string(data)
	// drop
	d = Replace(d, "\n")
	data = []byte(d)
	sdb := sqlDb()
	defer sdb.Close()
	_, err := sdb.Query("insert into orders(id, orderjson, pubdate) values($1, $2, $3)", order.Order_id, data, timestamp)
	if err != nil {
		//probably incoming message with last order but with new data. Trying to update order in db
		_, err = sdb.Query("update orders set orderjson=$1, pubdate=$2 where id=$3", data, timestamp, order.Order_id)
		if err != nil {
			return "", sql.ErrNoRows
		}
	}
	//if data valid and no other errors, data must be writtn in cache
	return order.Order_id, nil
}

// this func gets maxvalue of pubdate and generate delta between lats written date and now
func GetDelta() (delta int64) {
	sdb := sqlDb()
	defer sdb.Close()
	err := sdb.QueryRow("select max(pubdate) from orders").Scan(&delta)
	if err != nil {
		// if database contians nothing
		return 1
	}
	delta = time.Now().UnixNano() - delta
	return
}

// This func fill cache
func GetAll() ([]string, []string) {
	var key_list, value_list []string
	var key, value string
	sdb := sqlDb()
	query, err := sdb.Query("select id, orderjson from orders")
	helper.Errors(err, "query")
	for query.Next() {
		query.Scan(&key, &value)
		key_list = append(key_list, key)
		value_list = append(value_list, value)
	}
	return key_list, value_list
}
