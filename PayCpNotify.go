package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (cp GameIDRow) PostToCpNotify(userID string, fee int, orderID string, payTime string, cpDate string) error {
	var err error
	goto BEGIN

ERROR:
	fmt.Println(err)
	return err

BEGIN:
	feeStr := strconv.Itoa(fee)
	key := cp.PayKey         //"44F3411D4C4A490F91D2D16C893C5145"
	cpUrl := cp.PayNotifyUrl //"http://120.77.84.118:80"
	//cpUrl:= "http://127.0.0.1:5723"
	//  unorder for send
	requestPkg := url.Values{"userID": {userID}, "fee": {feeStr}, "orderID": {orderID}, "payTime": {payTime}, "cpDate": {cpDate}}
	// for calculte md5
	md5Pkg := url.Values{"cpDate": {cpDate}, "fee": {feeStr}, "orderID": {orderID}, "payTime": {payTime}, "userID": {userID}}
	md5content, err := url.QueryUnescape(md5Pkg.Encode())
	if err != nil {
		goto ERROR
	}
	md5sign := md5.Sum([]byte(md5content + key))
	md5str := hex.EncodeToString(md5sign[:])
	requestPkg.Set("signValue", md5str)

	fmt.Println(requestPkg)
	md5content, err = url.QueryUnescape(requestPkg.Encode())

	fmt.Println(md5content)
	fmt.Println(cpUrl)
	response, err := http.PostForm(cpUrl, requestPkg)
	if err != nil {
		goto ERROR
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}

	fmt.Println("reponse:", string(body))
	if string(body) != "success" {
		err = errors.New("Reponse Error")
		goto ERROR
	}

	return nil

}

func SaveOrderAndNotifyCP(gameIDMap map[int]GameIDRow, payOrderTable PayOrderTable, orderID string) (err error) {

	gameID, payChannelID, err := DecodePayOrderID(orderID)
	fmt.Println(orderID, gameID, payChannelID)
	///*
	_ = payChannelID
	tableID := gameID
	_ = tableID

	///*
	row, err := payOrderTable.Get(tableID, orderID)
	if err != nil {
		return err
	}
	fmt.Println(row)

	// get had pay notify
	err = payOrderTable.SavePayOrder(tableID, orderID, 1)
	fmt.Println(orderID)
	if err != nil {
		return err
	}
	//fmt.Println("success save")
	//fmt.Println("Send CP Date" + row.CpNotifyDate)

	accountEncode, err := CPUserIDEncode(row.AccountID, gameID)
	if err != nil {
		return err
	}
	//XXX:
	for i := 0; i < 6; i++ {
		fmt.Println(accountEncode, row.Fee, row.OrderID, row.OrderTime, row.CpNotifyDate)
		//err = cp.PostToCpNotify(accountEncode, row.Fee, row.OrderID, row.OrderTime, row.CpNotifyDate)
		cp := gameIDMap[gameID]
		err = cp.PostToCpNotify(accountEncode, row.Fee, row.OrderID, row.OrderTime, row.CpNotifyDate)
		if err == nil {
			fmt.Println("success post")
			// success send pay cp
			err = payOrderTable.SavePayOrder(tableID, orderID, 2)

			if err != nil {
				return err
			}

			return nil
		}
		time.Sleep(time.Minute * 1)
	}

	fmt.Println("fail post")

	return errors.New("FAIL POST to CP")
	//*/
	return err

}
