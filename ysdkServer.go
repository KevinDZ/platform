package main

import (
	_ "bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	_ "crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	_ "net/http/httputil"
	"net/url"
	"strconv"
	_ "strings"
	"time"
)

const (
	//TODO:
	CONST_YSDK_PLATFORM_ID int = 2

	//TODO:
	WeixinID      string = "wx596a3a5f4daf7a46"
	WeixinAppKey  string = "8db60cd46aac80afe5cdb5d628c82d2f"
	WeixinAuthUrl string = "http://ysdktest.qq.com/auth/wx_check_token"
	//TODO:
	QQID      string = "1105974649"
	QQAppKey  string = "W6GAJ2Ua2vIXa5Hx"
	QQAuthUrl string = "http://ysdktest.qq.com/auth/qq_check_token"

	//TODO
	MiAppID string = "1105974649"

	MiKey string = "P9a73urMdlnl65Voy0750Sn7u8vtJGqw&"

	//TODO
	MiGetBalanceSigUrl string = "/v3/r/mpay/get_balance_m"
	MiGetBalanceUrl    string = "https://ysdk.qq.com/mpay/get_balance_m"
	MiPaySigUrl        string = "/v3/r/mpay/pay_m"
	MiPayUrl           string = "https://ysdk.qq.com/mpay/pay_m"
)

/* ysdk get */
type YSDKAccountParam struct {
	OpenID      string `json:"openid"`
	AccessToken string `json:"accessToken"`
	PayToken    string `json:"payToken"`
	Pf          string `json:"pf"`
	Pfkey       string `json:"pfkey"`
	Type        string `json:"type"`
}

type YSDKPayRequest struct {
	YSDKAccountParam YSDKAccountParam `json:"thirdPay"`
}

type YSDKLoginRequest struct {
	YSDKAccountParam YSDKAccountParam `json:"thirdLogin"`
}

func HMacSha1(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(expectedMAC)
}

func MiCost(payorderID, stype, openID, payToken, pf, pfkey, MiKey string, fee int) (err error) {

	err = nil
	goto BEGIN
ERROR:
	return err

BEGIN:
	timestamp := fmt.Sprint(time.Now().Unix())

	billno := payorderID
	value := "amt=" + fmt.Sprint(fee/10) +
		"&appid=" + MiAppID +
		"&billno=" + billno +
		"&format=json&openid=" + openID +
		"&openkey=" + payToken +
		"&pf=" + pf +
		"&pfkey=" + pfkey +
		"&ts=" + timestamp +
		"&zoneid=1"

	valueEscape := url.QueryEscape(value)
	str := "GET&%2Fv3%2Fr%2Fmpay%2Fpay_m&" + valueEscape
	sig := HMacSha1(([]byte)(str), ([]byte)(MiKey))
	fmt.Println(sig)
	sig = url.QueryEscape(sig)
	fmt.Println(str)
	fmt.Println(value)
	fmt.Println(sig)

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	url := MiPayUrl + "?" + value + "&sig=" + sig
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		goto ERROR
	}

	if stype == "wx" {
		req.Header.Set("Cookie", "session_id=hy_gameid;session_type=wc_actoken;org_loc=/mpay/pay_m")
	} else if stype == "qq" {
		req.Header.Set("Cookie", "session_id=openid;session_type=kp_actoken;org_loc=/mpay/pay_m")
	} else {
		err = errors.New("type error")
		goto ERROR
	}

	fmt.Println(req.Cookie("session_id"))
	fmt.Println(req.Cookie("session_type"))
	fmt.Println(req.Cookie("org_loc"))

	response, err := client.Do(req)
	if err != nil {
		goto ERROR
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}

	fmt.Println(string(body))

	type ResultRequest struct {
		Msg string `json:"msg"`
		Ret int    `json:"ret"`
	}
	var result ResultRequest
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		err = errors.New("reulst json parse error")
		goto ERROR
	}
	if result.Ret != 0 {
		err = errors.New(result.Msg)
		goto ERROR
	}
	return nil

}

func YSDKLogin(openID, accessToken, loginType string) (platformAccountID string, err error) {

	goto BEGIN
ERROR:
	return "", err

BEGIN:
	var appID, appKey, authUrl string

	if loginType == "wx" {
		appID, appKey, authUrl = WeixinID, WeixinAppKey, WeixinAuthUrl
	} else if loginType == "qq" {
		appID, appKey, authUrl = QQID, QQAppKey, QQAuthUrl
	} else {
		err = errors.New("type error")
		goto ERROR
	}

	timestamp := fmt.Sprint(time.Now().Unix())

	//  md5 to hex
	context := appKey + timestamp
	md5encode := md5.Sum([]byte(context))
	md5encodeStr := hex.EncodeToString(md5encode[:])

	// get and check accountid is login in
	v := url.Values{}
	v.Set("appid", appID)
	v.Set("timestamp", timestamp)
	v.Set("openid", openID)
	v.Set("sig", md5encodeStr)
	v.Set("openkey", accessToken)

	response, err := http.Get(authUrl + "?" + v.Encode())
	if err != nil {
	}

	// parse result
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
	}
	fmt.Println(string(body))
	defer response.Body.Close()

	type ResultRequest struct {
		Msg string `json:"msg"`
		Ret int    `json:"ret"`
	}

	var result ResultRequest
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		err = errors.New("third login json parse error")
		goto ERROR
	}

	if result.Ret != 0 {
		err = errors.New(result.Msg)
		goto ERROR
	}
	return openID + "@" + loginType, nil

}

func YSDKPay(c *gin.Context, payOrderTable PayOrderTable, gameIDMap map[int]GameIDRow, inputFee int, inputPayOrderID string) (err error) {

	type ThirdPayRequest struct {
		GameAccessKey string `json:"key"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		DeviceInfo    string `json:"deviceInfo"`
		//TODO:
		Fee        string `json:"fee"`
		PayOrderID string `json:"payOrderID"`
	}
	type PayRespond struct {
		Result     string `json:"result"`
		PayOrderID string `json:"payOrderId"`
		PayUrl     string `json:"payUrl"`
		Message    string `json:"message"`
	}

	goto BEGIN
ERROR:
	LogE(c, err.Error())
	c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: "申请订单失败。", PayUrl: ""})
	return

BEGIN:
	requestJson := c.PostForm("request")
	var request ThirdPayRequest
	var ysdkPayRequest YSDKPayRequest

	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		err = errors.New("request json parse error")
		goto ERROR
	}

	err = json.Unmarshal([]byte(requestJson), &ysdkPayRequest)
	if err != nil {
		err = errors.New("request json parse error")
		goto ERROR
	}
	ysdkAccountParam := ysdkPayRequest.YSDKAccountParam

	//fee
	var fee int
	if inputFee == 0 {
		fee, err = strconv.Atoi(request.Fee)
		if err != nil || fee < 10 {
			err = errors.New("fee not correct")
			goto ERROR
		}
	} else {
		fee = inputFee
	}
	var payOrderID string
	if inputPayOrderID != "" {
		payOrderID = inputPayOrderID
	} else {
		payOrderID = request.PayOrderID
	}
	fmt.Println(payOrderID)

	gameID, payChannelID, err := DecodePayOrderID(payOrderID)
	_ = payChannelID

	tableID := gameID
	fmt.Println(gameID)
	row, err := payOrderTable.Get(tableID, payOrderID)
	if err != nil {
		goto ERROR
	}
	if row == nil {
		err = errors.New("row nil")
		goto ERROR
	}
	fmt.Println(fee, row.Fee)
	if fee > row.Fee {
		err = errors.New("fee not correct")
		goto ERROR
	}

	err = MiCost(payOrderID, ysdkAccountParam.Type, ysdkAccountParam.OpenID,
		ysdkAccountParam.PayToken, ysdkAccountParam.Pf, ysdkAccountParam.Pfkey, MiKey, fee)
	if err != nil {
		goto ERROR
	}

	//pay success
	err = SaveOrderAndNotifyCP(gameIDMap, payOrderTable, payOrderID)
	if err != nil {
		goto ERROR
	}
	return nil

}

/* main */
func ParsePlatformAccountID(platformID int, requestJson string) (platformAccountID string, err error) {
	fmt.Println(requestJson)
	var request YSDKLoginRequest
	err = json.Unmarshal([]byte(requestJson), &request)
	fmt.Println("platformID:", platformID)
	if platformID == CONST_YSDK_PLATFORM_ID {
		platformAccountID, err = YSDKLogin(request.YSDKAccountParam.OpenID, request.YSDKAccountParam.AccessToken, request.YSDKAccountParam.Type)
		if err != nil {
			goto ERROR
		}
		return platformAccountID, nil
	}
ERROR:
	fmt.Println("Error:platformID:", platformID)
	return "", err
}

func ThirdLogin(c *gin.Context) {

	type ThirdLoginRequest struct {
		GameAccessKey string `json:"key"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		DeviceInfo    string `json:"deviceInfo"`
		// for account
		PlatformAccountID string `json:"platformAccountID"`
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	requestJson := c.PostForm("request")

	var request ThirdLoginRequest

	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		err = errors.New("request json parse error")
		goto ERROR
	}

	// platformID
	platformID, err := strconv.Atoi(request.PlatformID)
	if err != nil {
		err = errors.New("platformID not correct")
		goto ERROR
	}

	platformAccountIDStr, err := ParsePlatformAccountID(platformID, requestJson)
	if err != nil {
		goto ERROR
	}

	// add to parse
	request.PlatformAccountID = platformAccountIDStr
	json, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(json)}}
	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(body))
	defer response.Body.Close()
	c.String(http.StatusOK, string(body))
}

func main() {

	// database
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println(err)
		return
	}

	var payOrderTable PayOrderTable
	payOrderTable = PayOrderTable{db}

	/* get gameID */
	gameIDTable := GameIDTable{db}
	idGameIDMap, err := gameIDTable.GetAllByID()
	if err != nil {
		fmt.Println(err)
		return
	}
	accessGameIDMap, err := gameIDTable.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = accessGameIDMap

	gin.SetMode(GIN_MODE)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	ginUseLogger(router)

	// Set up CORS middleware options
	config := cors.Config{
		Origins:         "*",
		RequestHeaders:  "Authorization",
		Methods:         "GET, POST, PUT",
		Credentials:     false,
		ValidateHeaders: false,
		MaxAge:          1 * time.Minute,
	}

	// Apply the middleware to the router (works on groups too)
	router.Use(cors.Middleware(config))

	router.POST("/ysdkPay", func(c *gin.Context) {

		type Request struct {
			PlatformID string `json:"platformID"`
		}

		type Respond struct {
			Result  string `json:"result"`
			Message string `json:"message"`
		}

		requestJson := c.PostForm("request")
		fmt.Println(requestJson)
		var request Request

		var err error
		goto BEGIN
	ERROR:
		LogE(c, err.Error())
		return

	BEGIN:
		err = json.Unmarshal([]byte(requestJson), &request)
		if err != nil {
			err = errors.New("request json parse error")
			goto ERROR
		}

		platformID, err := strconv.Atoi(request.PlatformID)
		if err != nil {
			err = errors.New("platformID not correct")
			goto ERROR
		}

		_ = platformID
		if platformID == CONST_YSDK_PLATFORM_ID {
			_, _ = payOrderTable, idGameIDMap
			err = YSDKPay(c, payOrderTable, idGameIDMap, 0, "")
			if err != nil {
				goto ERROR
			}
		}
		c.JSON(http.StatusOK, Respond{Result: "SUCCESS", Message: ""})
	})

	///*
	router.POST("/checkPay", func(c *gin.Context) {
		///*
		goto BEGIN
	ERROR:
		fmt.Println(err)
		//LogE(c, err)
		return
	BEGIN:
		fmt.Println("checkPay")

		requestJson := c.PostForm("request")
		accountID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo, err := ParseCommonPayRequest(accessGameIDMap, requestJson)
		if err != nil {
			goto ERROR
		}
		_, _, _, _, _, _, _, _ = accountID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo

		t1 := (time.Now()).Format("2006-01-02 15:04:05")
		t2 := time.Now().Local().Add(-40 * time.Minute).Format("2006-01-02 15:04:05")

		list, err := payOrderTable.GetOrderList(gameID, accountID, platformID, t2, t1)
		if err != nil {
			goto ERROR
		}

		fmt.Println(list)
		for _, row := range list {
			_ = row
			fmt.Println(row.Fee, row.OrderID)
			err = YSDKPay(c, payOrderTable, idGameIDMap, row.Fee, row.OrderID)
		}
		//_ = list
		//*/
	})
	//*/

	router.POST("/ysdkLogin", func(c *gin.Context) {
		ThirdLogin(c)
	})
	router.POST("/thirdLogin", func(c *gin.Context) {
		ThirdLogin(c)
	})
	// TODO:
	//router.POST("/ysdkPay", func(c *gin.Context) {
	//router.POST("/thirdPay", func(c *gin.Context) {
	//	ThirdPay(c)
	//})

	http.ListenAndServe(":7373", router)
}
