 package main

import (
	//"bytes"
	//"crypto/md5"
	//_ "crypto/sha256"
	//"database/sql"
	"crypto/hmac"
	"crypto/sha1"
	//"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	"io/ioutil"
	"net/http"
	"net/url"
/*	"net/http/cookiejar"
	_ "net/http/httputil"
	_ "strings"
	"time"*/
	"strconv"
)

const (
	//TODO:
	CONST_MI_PLATFORM_ID int = 5
	//CONST_MI_TRADE_SUCCESS string = "TRADE_SUCCESS"
	)
var (
	bufId string
	bufKey string
	bufSecretKey string
)

 func main(){
	 router := gin.Default()
	 router.POST("/miLogin",func (c *gin.Context) {
		 MiLogin(c)
	 })
	 /*
	router.POST("/miPayNotify",func (c *gin.Context) {
		MiPayNotify(c)
	 })*/
	router.GET("/miPayNotify",func (c *gin.Context) {
		MiPayNotify(c)
	 })
	/*router.POST("/queryorder",func (c *gin.Context) {
		//MiQueryOrder(c)
	 })*/
	 router.Run(":8888")
 }

func HMacSha1(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	fmt.Println(expectedMAC)
	return string(expectedMAC)
	//return base64.StdEncoding.EncodeToString(expectedMAC)
	/*decodeurl,err := url.QueryUnescape(string(expectedMAC))
	fmt.Println(decodeurl)
	if err != nil {
		err = errors.New("URLEncoding response json parse error")
	}
	return decodeurl*/

}

/*func DesEncrypt(origData, key []byte) ([]byte, error) {
     block, err := des.NewCipher(key)
     if err != nil {
          return nil, err
     }
     origData = PKCS5Padding(origData, block.BlockSize())
     // origData = ZeroPadding(origData, block.BlockSize())
     blockMode := cipher.NewCBCEncrypter(block, key)
     crypted := make([]byte, len(origData))
      // 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
     // crypted := origData
     blockMode.CryptBlocks(crypted, origData)
     return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
     padding := blockSize - len(ciphertext)%blockSize
     padtext := bytes.Repeat([]byte{byte(padding)}, padding)
     return append(ciphertext, padtext...)
}*/

func MiLogin(c *gin.Context) {
	/*message := c.PostForm("request")
	c.String(http.StatusOK,message)
	fmt.Println(message)*/

	type MiLoginRequest struct {
		ThirdLogin struct {
			Session string `json:"session"`
			Uid string `json:"uid"`
			AppId string `json:"appId"`
			AppKey string `json:"appKey"`
			AppSecretKey string `json:"appSecret"`
		}	`json:"thirdLogin"`
		Key string `json:"key"`
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	requestJson := c.PostForm("request")
	fmt.Println(requestJson)
	//c.String(http.StatusOK,requestJson)
	var loginrequest MiLoginRequest
	err = json.Unmarshal(([]byte)(requestJson), &loginrequest)
	if err != nil {
		err = errors.New("loginrequest Unmarshal json parse error")
		goto ERROR
	}
	//反序列 json
	fmt.Println("loginrequest Unmarshal: ",loginrequest)
	//c.String(http.StatusOK,loginrequest)
	//TODO   problem
	//小米签名
	sig := "appId=" + loginrequest.ThirdLogin.AppId + "&session=" + loginrequest.ThirdLogin.Session + "&uid=" + loginrequest.ThirdLogin.Uid //带签名字符串
	//sign := base32.URLEncoding.EncodeToString([]byte(sig))
	//sign := base64.URLEncoding.EncodeToString([]byte(sig))
	
	//UrlEncode is sign
	//sign := url.QueryEscape(sig)		//参与签名不需要URLencoding

	//appId...AppKey
	bufId = loginrequest.ThirdLogin.AppId
	bufKey = loginrequest.ThirdLogin.AppKey
	
	//AppSecretKey hmac-sha1 key
	MiKey := loginrequest.ThirdLogin.AppSecretKey
	
	//AppSecretKey
	bufSecretKey = MiKey

	signhmac := HMacSha1(([]byte)(sig), ([]byte)(bufSecretKey))
	//signhmac,_ := DesEncrypt(([]byte)(sig), ([]byte)(MiKey))

	fmt.Println("sig: ",sig)
	//fmt.Println("UrlEncode: ", sign)
	fmt.Println("appSecretKey:", bufSecretKey)
	fmt.Println("HMacSha1: ", signhmac)
	//signature 16位 进制输出
	signature := hex.EncodeToString([]byte(signhmac))
	fmt.Println("signature:",signature)

	type MiRequest struct {
		AppId string  `json:"appId"`
		Session string `json:"session"`
		Uid string `json:"uid"`
		Signature string `json:"signature"`
		//platformUserID
		PlatformUserID string `json:"platformUserID"`
	}

	//send MiLogin Get
	var request MiRequest
	request.AppId = loginrequest.ThirdLogin.AppId
	request.Session = loginrequest.ThirdLogin.Session
	request.Uid = loginrequest.ThirdLogin.Uid
	request.Signature = signature
	fmt.Println(request)	//MiRequest json struct
	urls := "http://mis.migc.xiaomi.com/api/biz/service/verifySession.do?appId=" + request.AppId + "&session=" + request.Session + "&uid=" + request.Uid + "&signature=" + request.Signature
	fmt.Println("urls: ",urls)
	//urlencode := url.QueryEscape(urls)
	//fmt.Println("urlencode:",urlencode)
	/*req,err := http.NewRequest("GET",url,bloginrequest)
	if err != nil {
		err = errors.New("NewRequest request json parse error")
		goto ERROR
	}
	req.Header.Set("Content-Type", "application/json")    
	client := &http.Client{}                      
	response, err := client.Do(req)
	if err != nil {
		err = errors.New("NewRequest request json parse error")
		goto ERROR
	}*/

	//resp, err := http.Get(urlencode)
	resp, err := http.Get(urls)
	if err != nil {
		err = errors.New("Get response json parse error")
		goto ERROR
	}
	fmt.Println("HTTP GET:",resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)		//[]byte
	fmt.Println(string(body))
	
	type MiLoginRespond struct {
		Errcode int32  `json:"errcode"`
		ErrMsg string  `json:"errMsg"`
	}

	var mirespond MiLoginRespond
	err = json.Unmarshal(body, &mirespond)
	if err != nil {
		err = errors.New("MiLoginRespond json parse error")
		goto ERROR
	}
	fmt.Println("MiLoginRespond json struct:",mirespond)
	//c.JSON(http.StatusOK,mirespond{Errcode: mirespond.Errcode, ErrMsg: mirespond.ErrMsg})
	/*
	// platformID
	platformUserID := request.Uid + "@mi" 

	// add to parse
	request.PlatformUserID = platformUserID
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(jsons)}}
	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	reqbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(reqbody))
	defer response.Body.Close()

	c.JSON(http.StatusOK,request)



	//接收本地服务器的数据
	type miClientRespond struct {
		Result string `json:"result"`
		Message string `json:"message"`
		UserID string `json:"userID"`
	}
	var mcrespond miClientRespond
	err = json.Unmarshal([]byte(reqbody), &mcrespond)
	if err != nil {
		err = errors.New("miClientRespond json parse error")
		goto ERROR
	}
	c.JSON(http.StatusOK,mcrespond)*/
	type ThirdLoginRequest struct {
		GameAccessKey string `json:"key"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		DeviceInfo   string  `json:"deviceInfo"`
		/*DeviceInfo    struct {
			DeviceId string `json:"DeviceIdIMEI"`
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
			SubscriberId string `json:"SubscriberIdIMSI"`
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
			}	 `json:"deviceInfo"`*/
		// for account
		PlatformAccountID string `json:"platformAccountID"`
	}

	/*var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:*/
	 


	// add to parse
	tlrequest.PlatformAccountID =  request.Uid + "@im"
	fmt.Println(tlrequest.PlatformAccountID)
	json, err := json.Marshal(tlrequest)
	if err != nil {
		goto ERROR
	}
	fmt.Println(tlrequest)
	requestPkg := url.Values{"request": {string(json)}}
	fmt.Println(string(json))

	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	//response, err := http.PostForm("http://test.tanchenggame.com:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	body1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	//fmt.Println(response)
	fmt.Println(string(body1))
	defer response.Body.Close()
	c.String(http.StatusOK, string(body1))
	//c.JSON(http.StatusOK, string(body1))
}

func MiPayNotify(c *gin.Context) {
	/*	
	//moblie pay
	type MiPayRequest struct {
		CpOrderId string `json:"cpOrderId"`
		CpUserInfo string `json:"cpUserInfo"`
		Mibi string `json:"mibi"`
	}*/
	//SUCCESS pay result notify request
	type MiPayRequest struct {
		AppId string `json:"appId"`
		CpOrderId string `json:"cpOrderId"`
		CpUserInfo string `json:"cpUserInfo"`
		Uid string `json:"uid"`
		OrderId string `json:"orderId"`
		OrderStatus string `json:"orderStatus"`
		//pay fee
		PayFee string `json:"payFee"`
		ProductCode string `json:"productCode"`
		ProductName string `json:"productName"`
		ProductCount string `json:"productCount"`
		PayTime string `json:"payTime"`
		OrderConsumeType string `json:"orderConsumeType"`
		PartnerGiftConsume string `json:"partnerGiftConsume"`
		Signature string `json:"signature"`
	}
	//FAIL pay notify respond
	type MiPayRespond struct {
		Errcode int32  `json:"errcode"`
		ErrMsg string  `json:"errMsg"`
	}
	
	//message := c.PostForm("request")
	//result := c.PostForm("result")
	//fmt.Println("MiUniPayOnline Request: ",message)
	//fmt.Println("MiUniPayOnline Result: ",result)
	//c.JSON(http.StatusOK, string(message))
	
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	//buffer := make([]byte, 1024)
	//n, _ := c.Request.Body.Read(buffer)
	var request MiPayRequest
	request.AppId = c.Request.URL.Query().Get("appId")
	request.CpOrderId = c.Request.URL.Query().Get("cpOrderId")
	request.CpUserInfo = c.Request.URL.Query().Get("cpUserInfo")
	request.Uid = c.Request.URL.Query().Get("uid")
	request.OrderId = c.Request.URL.Query().Get("orderId")
	request.OrderStatus = c.Request.URL.Query().Get("orderStatus")
	request.PayFee = c.Request.URL.Query().Get("payFee")
	request.ProductCode =  c.Request.URL.Query().Get("productCode")
	request.ProductName =  c.Request.URL.Query().Get("productCode")
	request.ProductCount =  c.Request.URL.Query().Get("productCode")
	request.PayTime =  c.Request.URL.Query().Get("PayTime")
	//request.OrderConsumeType =  c.Request.URL.Query().Get("OrderConsumeType")
	//request.PartnerGiftConsume =  c.Request.URL.Query().Get("PartnerGiftConsume")
	request.Signature =  c.Request.URL.Query().Get("signature")
	//defer c.Request.URL.Close()
	//URL := string(buffer[:n])
	miurl := "appId="+request.AppId+"&cpOrderId="+request.CpOrderId+"&cpUserInfo="+request.CpUserInfo+"&uid="+request.Uid + "&orderId=" + request.OrderId + "&orderStatus=" + request.OrderStatus + "&payFee=" + 
	request.PayFee + "&productCode=" + request.ProductCode +"&productName=" + request.ProductName + "&productCount=" + request.ProductCount + "&PayTime=" + request.PayTime +/*"&OrderConsumeType=" + 
	request.OrderConsumeType + "&partnerGiftConsume=" + request.PartnerGiftConsume +*/ "&signature=" + request.Signature
	fmt.Println("URL: ",miurl)
	fmt.Println("URLrequest: ",request)
	//fmt.Println(string(buffer[0:n]))
	//fmt.Println("URL:",URL)
/*
	//var respond MiPayRespond
	err = json.Unmarshal([]byte(message), &request)
	fmt.Println("request:",request)
	if err != nil {
		err = errors.New("MiPayRequest json parse error")
		goto ERROR
	}
	request.AppId = c.Param("appId")
	request.CpOrderId = c.Param("cpOrderId")
	name := c.Param("name")
	fmt.Println(name)
	c.String(http.StatusOK ,name)*/
	
	//if request.OrderStatus == CONST_MI_TRADE_SUCCESS {
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(jsons)}}
	fmt.Println(string(jsons))

	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	//response, err := http.PostForm("http://test.tanchenggame.com:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	mibody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	//fmt.Println(response)
	fmt.Println(string(mibody))
	defer response.Body.Close()
	//c.JSON(http.StatusOK, respond{Errcode: respond.Errcode})	
	c.String(http.StatusOK, string(mibody))
	c.JSON(http.StatusOK, string(mibody))
	//}

}

func POST(c *gin.Context){
	//{"cpOrderId":"20170626180830_0001_0005_88fYNl","payFee":"100","payTime":"2017-06-26 18:09:05"}
	type ThirdPayNotifyRequest struct {
		OrderID string `json:"cpOrderId"`
		Fee string `json:"payFee"`
		PayTime string `json:"payTime"`
	}
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	requestJson := c.PostForm("request")
	fmt.Println("MiPayNotify:",payrequest)
	
	var request ThirdPayNotifyRequest 
	err = json.Unmarshal(requestJson,&request)
	if err != nil {
		err = errors.New("MiPayNotifyRequest json parse error")
		goto ERROR
	}

	payOrderID := request.OrderID
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
	payTime := request.PayTime
	//储存mi的数据
	err = SaveOrderAndNotifyCP(gameIDMap, payOrderTable, payOrderID)
	if err != nil {
		goto ERROR
	}
	return nil
}
/*
func  MiQueryOrder(c *gin.Context) {
	msg:= c.PostForm("request")
	fmt.Println("MiQueryOrder Request:",msg)
	type  MIQueryOrderRequest struct {
		AppId string `json:"appId"`
		CpOrderId string `json:"cpOrderId"`
		Uid string `json:"uid"`
		Signature string `json:"signature"`
	}
	//SUCCESS CpOrder query notify respond
	type Respond struct {
		AppId string `json:"appId"`
		CpOrderId string `json:"cpOrderId"`
		CpUserInfo string `json:"cpUserInfo"`
		Uid string `json:"uid"`
		OrderId string `json:"orderId"`
		OrderStatus string `json:"orderStatus"`
		PayFee string `json:"payFee"`
		ProductCode string `json:"productCode"`
		ProductName string `json:"productName"`
		ProductCount string `json:"productCount"`
		PayTime string `json:"payTime"`
		OrderConsumeType string `json:"orderConsumeType"`
		PartnerGiftConsume string `json:"partnerGiftConsume"`
		Signature string `json:"signature"`
	}	

	//FAIL  CpOrder query notify respond
	type MiQueryOrderRespond struct {
		Errcode int32  `json:"errcode"`
		ErrMsg string  `json:"errMsg"`
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	var request MIQueryOrderRequest

	err = json.Unmarshal([]byte(msg), &request)
	fmt.Println("request:",request)
	if err != nil {
		err = errors.New("MiOrderQueryRequest json parse error")
		goto ERROR
	}
	url := "http://mis.migc.xiaomi.com/api/biz/service/queryOrder.do"
	resp,err := http.Get(url)
	if err != nil {
		err = errors.New("MiOrderQuery Get json parse error")
		goto ERROR
	}
	fmt.Println("MiOrderQuery resp:",resp)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(bodyBytes))
	defer resp.Body.Close()
	c.String(http.StatusOK, string(bodyBytes))
}*/

