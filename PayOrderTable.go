package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "strconv"
	"time"
)

/*
CREATE TABLE "pay0000"
(
  order_id character varying(256) PRIMARY KEY,
  account_id bigint NOT NULL,
  server_id integer NOT NULL,
  device_id character varying(256),
  platform_id integer NOT NULL,
  channel_id integer NOT NULL,
  fee integer NOT NULL,
  status integer NOT NULL,
  order_time timestamp without time zone NOT NULL,
  pay_time timestamp without time zone,
  cp_notify_date character varying(256),
  third_pay_order_id character varying(256) NOT NULL,
  third_pay_channel_id integer NOT NULL,
  third_pay_url character varying(256) NOT NULL
)


CREATE TABLE "payOrder"
(
  order_id character varying(256) NOT NULL,
  user_id integer NOT NULL,
  fee integer NOT NULL,
  status integer NOT NULL,
  order_time timestamp without time zone NOT NULL,
  pay_time timestamp without time zone,
  cp_notify_date character varying(256),
  third_pay_order_id character varying(256) NOT NULL,
  third_pay_channel_id integer NOT NULL,
  third_pay_url character varying(256) NOT NULL,
  CONSTRAINT order_pkey PRIMARY KEY (order_id)
)
*/
//game_id integer NOT NULL,

type PayOrderTable struct {
	Db *sql.DB
}

type PayOrderRow struct {
	OrderID           string
	AccountID         int64
	ServerID          int
	DeviceID          string
	PlatformID        int
	ChannelID         int
	Fee               int
	Status            int
	OrderTime         string
	PayTime           string
	CpNotifyDate      string
	ThirdPayOrderID   string
	ThirdPayChannelID int
	ThirdPayUrl       string
}

func (table PayOrderTable) getTableName(tableID int) string {
	return "pay_" + fmt.Sprintf("%04d", tableID) + "_table"
}

func (table PayOrderTable) GetCount(tableID int, accountID int64, platformID int) (count int, err error) {
	sqlStr := `select count(*) from ` + table.getTableName(tableID) +
		` where account_id = '` + fmt.Sprint(accountID) + `'` +
		` and platform_id = '` + fmt.Sprint(platformID) + `'` +
		` and status > '0'`
	fmt.Println(sqlStr)
	rows, err := table.Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		//goto ERROR
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil || rows.Next() != false {
			fmt.Println(err)
			return -1, err
		}
	}
	return count, nil
}

func (table PayOrderTable) GetOrderList(tableID int, accountID int64, platformID int, startTime, endTime string) (list []PayOrderRow, err error) {

	var tableName = table.getTableName(tableID)
	var slice []PayOrderRow
	var payOrderRow PayOrderRow
	sqlStr := "SELECT order_id, account_id, fee, status,COALESCE(to_char(order_time, 'YYYY-MM-DD|HH24:MI:SS'), '') , COALESCE(to_char(pay_time, 'YYYY-MM-DD|HH24:MI:SS'), ''), cp_notify_date , third_pay_order_id , third_pay_id , third_pay_url " +
		" from " + tableName +
		` where account_id = '` + fmt.Sprint(accountID) + `'` +
		` and platform_id = '` + fmt.Sprint(platformID) + `'` +
		` and order_time >= '` + startTime + `'` +
		` and order_time < '` + endTime + `'` +
		` and status = '0'` +
		` order by order_time desc`

	fmt.Println(sqlStr)
	rows, err := table.Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		goto ERROR
	}
	defer rows.Close()
	list = slice

	for rows.Next() {
		err = rows.Scan(&payOrderRow.OrderID,
			&payOrderRow.AccountID,
			&payOrderRow.Fee,
			&payOrderRow.Status,
			&payOrderRow.OrderTime,
			&payOrderRow.PayTime,
			&payOrderRow.CpNotifyDate,
			&payOrderRow.ThirdPayOrderID,
			&payOrderRow.ThirdPayChannelID,
			&payOrderRow.ThirdPayUrl)
		if err != nil {
			goto ERROR
		}

		slice = append(slice, payOrderRow)
	}
	goto SUCCESS

ERROR:
	fmt.Println(err)
	return list, err

SUCCESS:
	list = slice
	return list, nil

}

/* tableID == gameID */
func (table PayOrderTable) Get(tableID int, orderID string) (retRow *PayOrderRow, err error) {
	fmt.Println(tableID, orderID)
	var tableName = table.getTableName(tableID)
	var payOrderRow PayOrderRow

	var rows *sql.Rows
	retRow = nil

	sqlStr := "SELECT order_id, account_id, fee, status,COALESCE(to_char(order_time, 'YYYY-MM-DD|HH24:MI:SS'), '') , " +
		"COALESCE(to_char(pay_time, 'YYYY-MM-DD|HH24:MI:SS'), ''), cp_notify_date , third_pay_order_id , third_pay_id , third_pay_url FROM " +
		tableName + " WHERE order_id='" + orderID + "'"
	fmt.Println(orderID)
	rows, err = table.Db.Query(sqlStr)
	if err != nil {
		goto ERROR
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&payOrderRow.OrderID,
			&payOrderRow.AccountID,
			&payOrderRow.Fee,
			&payOrderRow.Status,
			&payOrderRow.OrderTime,
			&payOrderRow.PayTime,
			&payOrderRow.CpNotifyDate,
			&payOrderRow.ThirdPayOrderID,
			&payOrderRow.ThirdPayChannelID,
			&payOrderRow.ThirdPayUrl)
		if err != nil {
			goto ERROR
		}

		// if more rows return err
		if rows.Next() != false {
			err = errors.New("return mutiple rows")
			goto ERROR
		}

		retRow = &payOrderRow
	}
	goto SUCCESS

ERROR:
	fmt.Println(err)
	return nil, err

SUCCESS:
	return retRow, nil
}

/* tableID == gameID */
func (table PayOrderTable) NewBookOrder(tableID int,
	orderID string,
	orderTime string,
	accountID int64,
	serverID int,
	deviceID string,
	platformID int,
	channelID int,
	fee int,
	cpNotifyDate string,
	thirdPayOrderID string,
	thirdPayChannelID int,
	thirdPayUrl string) (err error) {

	var tableName = table.getTableName(tableID)
	var status int = 0

	var retOrderID string
	err = table.Db.QueryRow("INSERT INTO "+tableName+" ( "+
		"order_id, "+
		"account_id, "+
		"server_id, "+
		"device_id, "+
		"platform_id, "+
		"channel_id, "+
		"fee, "+
		"status, "+
		"order_time, "+
		//"pay_time, "+
		"cp_notify_date, "+
		"third_pay_order_id, "+
		"third_pay_id, "+
		"third_pay_url ) "+
		" VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13  ) RETURNING order_id",
		orderID,
		accountID,
		serverID,
		deviceID,
		platformID,
		channelID,
		fee,
		status,
		orderTime,
		//PayTime,
		cpNotifyDate,
		thirdPayOrderID,
		thirdPayChannelID,
		thirdPayUrl).Scan(&retOrderID)

	if err != nil {
		goto ERROR
	}
	goto SUCCESS

ERROR:
	fmt.Println(err)
	return err

SUCCESS:
	return nil
}

//pay
/* tableID == gameID */
func (table PayOrderTable) SavePayOrder(tableID int, orderID string, status int) (err error) {

	payTime := (time.Now()).Format("2006-01-02 15:04:05")
	var tableName = table.getTableName(tableID)

	var res sql.Result
	var affect int64

	stmt, err := table.Db.Prepare("UPDATE " + tableName + " SET status=$2, pay_time=$3 WHERE order_id=$1")
	if err != nil {
		goto ERROR
	}
	defer stmt.Close()

	res, err = stmt.Exec(orderID, status, payTime)
	if err != nil {
		goto ERROR
	}

	affect, err = res.RowsAffected()
	if err != nil {
		goto ERROR
	}
	if affect != 1 {
		err = errors.New("return mutiple rows")
		goto ERROR
	}

	goto SUCCESS

ERROR:
	fmt.Println(err)
	return err

SUCCESS:
	return nil

}

/*
func main() {
  dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
  DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
  db, err := sql.Open("postgres", dbinfo)

  if err != nil {
	fmt.Println(err)
	return
  }

  gameID := 1
  orderTime := (time.Now()).Format("2006-01-02 15:04:05")
  // table
  var payOrderTable PayOrderTable
  payOrderTable = PayOrderTable{db}
  _, _ = gameID, orderTime
  orderID := "11111"

  err = payOrderTable.NewBookOrder(gameID, orderID, orderTime, 1, 0, "device", 0, 0, 100, "cpdate", "third pay Order id", 1, "http://pay url")
  if err != nil {
	fmt.Println(err)
  }
  err = payOrderTable.SavePayOrder(gameID, orderID, 1)
  if err != nil {
	fmt.Println(err)
  }
  //fmt.Println(d)
  q, err := payOrderTable.Get(gameID, orderID)
  if err != nil {
	fmt.Println(err)
  }
  fmt.Println(q)
  list, err := payOrderTable.GetOrderList(1, 54, 0, "2017-01-01", "2018-01-01")
  fmt.Println(list)

  defer db.Close()

}

*/
