package main

import (
	"database/sql"
	_ "errors"
	"fmt"
	_ "strconv"
	"time"

	_ "github.com/lib/pq"
)

type BaiduTable struct {
	Db *sql.DB
}

type BaiduRow struct {
	MuidMd5 string
	/*Appid      string
	AppType    string
	ClickId    string
	AdvertseId string
	ClickTime  string*/
}

func (table BaiduTable) GetTableName(tableID int) string {
	return "baidu_" + fmt.Sprintf("%04d", tableID) + "_table"
}

func (table BaiduTable) Update(muidMd5 string, gameID, platformID int) (err error) {
	// TODO:
	var sqlStr string

	sqlStr = `update ` + table.GetTableName(gameID) +
		` set hadsend = 'true' ` +
		` where muid_md5='` + muidMd5 + `'` +
		` and platform_id='` + fmt.Sprint(platformID) + `'`

	fmt.Println(sqlStr)
	rows, err := table.Db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		goto ERROR
	}

	goto SUCCESS

ERROR:
	fmt.Println(err)
	return err

SUCCESS:
	return nil
}

func (table BaiduTable) GetLastChannel(muidMd5 string, gameID, platformID int) (channelID, lasttime string, hadsend bool, err error) {
	sql := `select cast(extract(epoch from receive_time) as bigint), channel_id, hadsend from ` + table.GetTableName(gameID) +
		` where muid_md5= '` + muidMd5 + `'` +
		` and platform_id = '` + fmt.Sprint(platformID) + `'` +
		` order by hadsend desc, receive_time desc `
	fmt.Println("sql:", sql)
	rows, err := table.Db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return "-1", "-1", false, err
	}
	defer rows.Close()
	fmt.Println("rows:", rows)
	for rows.Next() {
		err = rows.Scan(&lasttime, &channelID,  &hadsend)
		if err != nil {
			fmt.Println(err)
			return "-1", "-1", false, err
		}
		break
	}
	fmt.Println("channelID:",channelID)
	fmt.Println("lasttime:",lasttime)
	return channelID, lasttime, hadsend, err
}

func (table BaiduTable) IsHadSend(muidMd5 string, gameID, platformID, channelID int) (hadSend bool, err error) {
	var sqlStr string

	sqlStr = `select * from ` + table.GetTableName(gameID) +
		` where muid_md5= '` + muidMd5 + `'` +
		` and hadsend= 'true'` +
		` and platform_id = '` + fmt.Sprint(platformID) + `'` +
		` order by receive_time desc `

	fmt.Println(sqlStr)
	rows, err := table.Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return true, err
	}

	defer rows.Close()
	for rows.Next() {
		fmt.Println("had send")
		return true, nil

	}
	fmt.Println("not send")
	return false, nil
}

func (table BaiduTable) Insert(gameID, platformID, channelID int, muidMd5 string) (err error) {
	var retCallbackUrl string
	receiveTime := (time.Now()).Format("2006-01-02 15:04:05")
	sql := "INSERT INTO " + table.GetTableName(gameID) + " (platform_id, channel_id, muid_md5, receive_time, hadsend) VALUES(" +
		"'" + fmt.Sprint(platformID) + "'," +
		"'" + fmt.Sprint(channelID) + "'," +
		"'" + muidMd5 + "'," +
		"'" + receiveTime + "'," +
		"'false') RETURNING muid_md5"
	fmt.Println("sql:", sql)
	err = table.Db.QueryRow(sql).Scan(&retCallbackUrl)
	if err != nil {
		goto ERROR

	}
	fmt.Println("retCallbackUrl:",retCallbackUrl)
	goto SUCCESS
ERROR:
	fmt.Println(err)
	return err
SUCCESS:
	return nil
}
