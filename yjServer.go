package main

import (
	_ "bytes"
	_ "crypto/hmac"
	 "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_"database/sql"
	_"encoding/base64"
	_"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"io/ioutil"
	"net/http"
	_"net/http/cookiejar"
	_ "net/http/httputil"
	"net/url"
	_ "strconv"
	_ "strings"
	"io"	
	"time"
)
const (
	//{C51C59F1-401B9EA0}
	AppID = "C51C59F1-401B9EA0"
	YJURL = "http://sync.1sdk.cn/login/check.html"
	YJKey = ""
)

func main() {
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

	router.POST("/yjLogin", func(c *gin.Context) {
		YJLogin(c)
	})
	router.POST("/yjPayNotify", func (c *gin.Context) {
		YJPayNotify(c)
		})
	router.Run(":8787")
}

func YJLogin(c *gin.Context) {
	type ThirdLoginRequest struct {
		ThirdLogin struct{
			SDK string `json:"sdk"`
			App string `json:"app"`
			Uin string `json:"uin"`
			Sess string `json:"sess"`			
		}	`json:"thirdLogin"`
		GameAccessKey string `json:"key"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		DeviceInfo   struct {
			DeviceId string `json:"DeviceId(IMEI)"`
			Mac string `json:"Mac"`
			DeviceSoftwareVersion string `json:"DeviceSoftwareVersion"`
			Line1Number string `json:"Line1Number"`
			NetworkCountryIso string `json:""NetworkCountryIso"`
			NetworkOperator string `json:""NetworkOperator"`
			NetworkOperatorName string `json:"NetworkOperatorName"`
			NetworkType string `json:"NetworkType"`
			PhoneType string `json:"PhoneType"`
			SimSerialNumber string `json:"SimSerialNumber"`
			SimState string `json:"SimState"`
			SubscriberId string `json:"SubscriberId(IMSI)"`
			VoiceMailNumber string `json:"VoiceMailNumber"`
			Product string `json:"Product"`
			CPU_ABI string `json:"CPU_ABI"`
			TAGS string `json:"TAGS"`
			VersionCodesBase string `json:"VERSION_CODES.BASE"`
			MODEL string `json:"MODEL"`
			SDK string `json:"SDK"`
			VersionRelease string `json:"VERSION.RELEASE"`
			DEVICE string `json:"DEVICE"`
			DISPLAY string `json:"DISPLAY"`
			BRAND string `json:"BRAND"`
			BOARD string `json:"BOARD"`
			FINGERPRINT string `json:"FINGERPRINT"`
			ID string `json:"ID"`
			MANUFACTURER string `json:"MANUFACTURER"`
			USER string `json:"USER"`
			OS string `json:"OS"`
			}	 `json:"deviceInfo"`	
		// for account
		PlatformAccountID string `json:"platformAccountID"`
	}

	type Respond struct {
		Result string `json:"result"`
		Message string `json:"message"`
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: err.Error()})
	return 

BEGIN:
	requestJson := c.PostForm("request")
	fmt.Println("requestJson:",requestJson)

	var loginrequest ThirdLoginRequest
	err = json.Unmarshal(([]byte)(requestJson), &loginrequest)
	if err != nil {
		err = errors.New("ThirdLoginRequest json parse error")
		goto ERROR
	}
	fmt.Println("ThirdLoginRequest:",loginrequest)

	v := url.Values{}
	v.Set("app", loginrequest.ThirdLogin.App)
	v.Set("sdk", loginrequest.ThirdLogin.SDK)
	v.Set("sess", loginrequest.ThirdLogin.Sess)
	v.Set("uin", loginrequest.ThirdLogin.Uin)

	resp, err := http.Get(YJURL + "?" + v.Encode())
	if err != nil {
		goto ERROR
	}
	fmt.Println("http.GET:",YJURL + "?" + v.Encode())
   	fmt.Println("response:",resp)
   	fmt.Println("resp.Body:",resp.Body)
    	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto ERROR
	}
	respStr := string(respBytes)
	fmt.Println("respStr:",respStr)
	defer resp.Body.Close()

	type YJRespond struct {
		Code int `json:"code"`
		Sess string `json:"sess"`
		Reason string `json:"reason"`
	}
	var yjrespond YJRespond
	err = json.Unmarshal(respBytes, &yjrespond)
	if err != nil {
		err = errors.New("BoloRespond request json parse error")
		goto ERROR
	}
	fmt.Println("YJRespond:",yjrespond)
	fmt.Println("结果码验证:",yjrespond.Code)
	if yjrespond.Code == 0 {
		fmt.Println("易接登陆正常")

		//易接的userID 使用哪个sess/uin/app
		// add to parse
		loginrequest.PlatformAccountID = yjrespond.Sess + "@bolo"
		jsons, err := json.Marshal(loginrequest)
		if err != nil {
			goto ERROR
		}
		requestPkg := url.Values{"request": {string(jsons)}}
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
	}else {
		fmt.Println("易接登陆异常")
		c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: err.Error()})
	}
}

func YJPayNotify(c *gin.Context) {
	type YJRespond struct {
		App string `json:"app"`
		Cbi string `json:"cbi"`
		Ct string `json:"ct"`
		Fee string `json:"fee"`
		Pt string `json:"pt"`
		SDK string `json:"sdk"`
		Ssid string `json:"ssid"`
		St string `json:"st"`
		Tcd string `json:"tcd"`
		Uid string `json:"uid"`
		Ver string `json:"ver"`
		Sign string `json:"sign"`
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())	
	return 

BEGIN:
	var respond YJRespond
	respond.App = c.Request.URL.Query().Get("app")
	respond.Cbi = c.Request.URL.Query().Get("cbi")
	respond.Ct = c.Request.URL.Query().Get("ct")		//支付完成时间
	respond.Fee = c.Request.URL.Query().Get("fee")
	respond.Pt = c.Request.URL.Query().Get("pt")		//付费时间，订单创建服务器UTC时间戳（毫秒）
	respond.SDK = c.Request.URL.Query().Get("sdk")
	respond.Ssid = c.Request.URL.Query().Get("ssid")	//订单在渠道平台上的流水号
	respond.St =  c.Request.URL.Query().Get("st")
	respond.Tcd =  c.Request.URL.Query().Get("tcd")	//订单在易接服务器上的订单号
	respond.Uid =  c.Request.URL.Query().Get("uid")	//付费用户在渠道平台上的唯一标记
	respond.Ver =  c.Request.URL.Query().Get("ver")
	respond.Sign =  c.Request.URL.Query().Get("sign")
	//sign
	sign := "app=" + respond.App + "&cbi=" + respond.Cbi + "&ct=" + respond.Ct + "&Fee=" + respond.Fee + "&pt=" + respond.Pt + "&sdk=" + respond.SDK + "&ssid=" +
		respond.Ssid + "&st=" + respond.St + "&tcd=" + respond.Tcd + "&uid=" + respond.Uid + "&Ver=" + respond.Ver

	//MD5
	md5str := MD5([]byte(sign + YJKey))
	cpSign := PAY_NOTIFY_HOST_URL + sign + "&sign=" + md5str
	fmt.Println("sign:",cpSign)
	fmt.Println("Sign:",respond.Sign)
	if respond.Sign == cpSign && respond.St == "1" {
		fmt.Println("签名验证通过")
		type ThirdRequest struct {
			PayFee string `json:"payFee"`
			PayTime string `json:"payTime"`
			CpOrderId string `json:"cpOrderId"`			
		}
		//third pay for request
		var request ThirdRequest
		request.PayFee = respond.Fee
		request.PayTime = respond.Ct
		request.CpOrderId = respond.Tcd
		jsons, err := json.Marshal(request)
		if err != nil {
			goto ERROR
		}
		fmt.Println("Third respond jsons:",string(jsons))

		requestPkg := url.Values{"request": {string(jsons)}}
		fmt.Println("requestPkg:",requestPkg)
				
		response, err := http.PostForm("http://127.0.0.1:8090/thirdPay", requestPkg)
		if err != nil {
			goto ERROR
		}
		var bodyBytes []byte
		bodyBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			goto ERROR
		}
		defer response.Body.Close()
		fmt.Println(string(bodyBytes))

		c.String(http.StatusOK,"SUCCESS")
	}
}


func MD5(message []byte) string {
	w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //md5 to hex
	return md5str2
}