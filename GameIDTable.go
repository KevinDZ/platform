package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "strconv"
	_ "time"
)

type GameIDTable struct {
	Db *sql.DB
}

type GameIDRow struct {
	GameID       int
	AccessKey    string
	PayKey       string
	PayNotifyUrl string
}

func (table GameIDTable) GetAll() (accessMap map[string]GameIDRow, err error) {
	var gameIDRow GameIDRow
	var rows *sql.Rows
	rows, err = table.Db.Query("SELECT game_id, access_key, pay_key, pay_notify_url FROM game_id_table")

	fmt.Println(err)

	accessMap = make(map[string]GameIDRow)

	for rows.Next() {
		err = rows.Scan(&gameIDRow.GameID,
			&gameIDRow.AccessKey,
			&gameIDRow.PayKey,
			&gameIDRow.PayNotifyUrl)

		if err != nil {
			goto ERROR
		}
		fmt.Println(gameIDRow)
		accessMap[gameIDRow.AccessKey] = gameIDRow
	}
	return accessMap, nil

ERROR:
	fmt.Println(err)
	return nil, err
}

func (table GameIDTable) GetAllByID() (accessMap map[int]GameIDRow, err error) {
	var gameIDRow GameIDRow
	var rows *sql.Rows
	rows, err = table.Db.Query("SELECT game_id, access_key, pay_key, pay_notify_url FROM game_id_table")

	fmt.Println(err)

	accessMap = make(map[int]GameIDRow)

	for rows.Next() {
		err = rows.Scan(&gameIDRow.GameID,
			&gameIDRow.AccessKey,
			&gameIDRow.PayKey,
			&gameIDRow.PayNotifyUrl)

		if err != nil {
			goto ERROR
		}
		fmt.Println(gameIDRow)
		accessMap[gameIDRow.GameID] = gameIDRow
	}
	return accessMap, nil

ERROR:
	fmt.Println(err)
	return nil, err
}

func (table GameIDTable) FixUrl(gameID int, url string) (err error) {

	var res sql.Result
	var affect int64

	stmt, err := table.Db.Prepare("UPDATE game_id_table SET pay_notify_url=$2 WHERE game_id=$1")
	if err != nil {
		goto ERROR
	}
	defer stmt.Close()

	res, err = stmt.Exec(gameID, url)
	if err != nil {
		goto ERROR
	}

	affect, err = res.RowsAffected()
	if err != nil {
		goto ERROR
	}
	if affect != 1 {
		err = errors.New("fix payNotify table return mutiple rows")
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
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "234567"
	DB_NAME     = "test"
	DB_HOST     = "120.77.84.118"
	DB_PORT     = "5432"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	fmt.Println(err)

	_ = db
	payOrderTable := GameIDTable{db}
	_ = payOrderTable
	keymap, idmap, _ := payOrderTable.GetAll()
	fmt.Println(keymap)
	fmt.Println(idmap)

	payOrderTable.FixUrl(1, "http://www.baidu.com")

}
*/
