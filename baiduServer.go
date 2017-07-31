package main

import (
	"bytes"
	"crypto/md5"
	/*"crypto/hmac"
	"crypto/sha1"
	_ "crypto/sha256"
	"database/sql"*/
	"encoding/hex"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	"io/ioutil"
	//"io"
	"net/http"
	//"net/http/cookiejar"
	_ "net/http/httputil"
	"unsafe"
	"net/url"
	"strings"
	//"unicode/utf8"
	//"encoding/binary"
	//"strconv"
	//"time"
)
const (
	BaiduSecretKey string = "m6T6cGCYAkjLZgFMBxlM3RcnpNYB37p6"
)
var (
	buf string 
)
func main() {
	router := gin.Default()

	router.POST("/baiduLogin", func(c *gin.Context) {
		BaiduLogin(c)
	})

	router.POST("/baiduOrderQuery", func(c *gin.Context) {
		BaiduOrderQuery(c)
	})

	router.POST("/baiduPayReceive", func(c *gin.Context) {
		BaiduPayReceive(c)
	})

	router.Run(":8686")
}

func MD5(message []byte) string {
	/*w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
	md5 := ([]byte)(md5str2)
	fmt.Println("md5Sign:",string(md5))
	return string(md5)*/
	//md5 to hex
	/*md5encode := md5.Sum([]byte(message))
	return hex.EncodeToString(md5encode[:])*/
	md5Ctx := md5.New()
    	md5Ctx.Write(message)
    	cipherStr := md5Ctx.Sum(nil)
    	fmt.Println(hex.EncodeToString(cipherStr))
    	return hex.EncodeToString(cipherStr)
}

func BaiduLogin(c *gin.Context) {
	type ThirdLoginRequest struct {
		ThirdLogin	struct{
		AppID string `json:"AppID"`
		AccessToken string `json:"AccessToken"`
		} 	`json:"thirdLogin"`
	}
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	var tlrequest ThirdLoginRequest
	requestJson := c.PostForm("request")
	fmt.Println("ThirdLoginRequest Json:",requestJson)
	err = json.Unmarshal([]byte(requestJson), &tlrequest)
	if err != nil {
		err = errors.New("ThirdLoginRequest json parse error")
		goto ERROR
	}
	//fmt.Println("Baidu login request:",tlrequest)
	//发送post请求百度json结构体
	type BaiduLoginRequest struct {
		AppID string `json:"AppID"`
		AccessToken string `json:"AccessToken"`
		Sign string `json:"Sign"`
	} 
	var request BaiduLoginRequest
	//发送百度sdk请求
	request.AppID = tlrequest.ThirdLogin.AppID
	fmt.Println("AppID:",request.AppID)
	request.AccessToken = tlrequest.ThirdLogin.AccessToken
	fmt.Println("AccessToken:",request.AccessToken)
	//sign签名字符串
	//signString := "appId" +  tlrequest.ThirdLogin.AppID + "accessToken" + tlrequest.ThirdLogin.AccessToken + BaiduSecretKey
	signString := tlrequest.ThirdLogin.AppID + tlrequest.ThirdLogin.AccessToken + BaiduSecretKey
	//fmt.Println("signString:",signString)
	signString = strings.Replace(signString,"-","",-1)
	fmt.Println("signString:",signString)
	//字符串转整型有问题
	/*signutf, _ := strconv.Atoi(signString)
	fmt.Println("signutf:",signutf)*/
	//字节数组
	/*signbyte := []byte(signString)
	fmt.Println("signbyte:",signbyte)
	//字节转整型
	 bytesBuffer := bytes.NewBuffer(signbyte)  
    	var tmp int32  
    	binary.Read(bytesBuffer, binary.BigEndian, &tmp)  
    	fmt.Println(int32(tmp))
	//signString := "AppID=" +  tlrequest.ThirdLogin.AppID + "&AccessToken=" + tlrequest.ThirdLogin.AccessToken + BaiduSecretKey
	//sign字符串的unicode码转为UTF-8
	utfbyte := make([]byte, utf8.UTFMax)
	n := utf8.EncodeRune(utfbyte, int32(tmp))
	fmt.Println("signString:",signString)
	fmt.Println("utfbyte:",n,utfbyte)*/
	//md5sign
	md5Sign := MD5([]byte(signString))
	//md5Sign := MD5(utfbyte)
	//md5sign 变加密结果,均转换为小写字符
	//md5Sign = strings.ToLower(md5Sign)
	request.Sign = md5Sign
	//fmt.Println("request:",request)

	// add to parse
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	fmt.Println("Baidu Login request jsons:",string(jsons))
	reader := bytes.NewBuffer(jsons)
	//发送给百度sdk服务器，登陆状态查询
	/*if {
		//sdk版本号<3.6.0
		baiduURL := "http://querysdkapi.91.com/cploginstatequery.ashx?"		
	}else {
		//sdk版本号>=3.6.0
		baiduURL := "http://querysdkapi.baidu.com/query/cploginstatequery?"
	}*/
	baiduURL := "http://querysdkapi.baidu.com/query/cploginstatequery?"
	//POST请求
	/*resp, err := http.Post(baiduURL, "application/json;charset=utf-8", nil)
	if err != nil { 
      		goto ERROR 
	}
	fmt.Println("resp:",resp)*/
	bdrequest, err := http.NewRequest("POST", baiduURL, reader)
    if err != nil {
        goto ERROR
    }
    //fmt.Println("baidu request:",bdrequest)
    //bdrequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(bdrequest)
    if err != nil {
        goto ERROR
    }
   //fmt.Println("resp:",resp)
  // fmt.Println("resp.Body:",resp.Body)
    respBytes, err := ioutil.ReadAll(resp.Body)
    fmt.Println(respBytes)
    if err != nil {
        goto ERROR
    }
    respStr := string(respBytes)
   // fmt.Println("respStr:",respStr)
    defer resp.Body.Close()
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("str:",*str)
	/*//百度返回参数
	type BaiduLoginRespond struct {
		AppID string `json:"AppID"`
		ResultCode string `json:"ResultCode"`
		ResultMsg string `json:"ResultMsg"`
		Sign string `json:"Sign"`
		Content struct {
			UID string `json:"UID"`
		} `json:"Content"`
		
	}
	var  loginrespond BaiduLoginRespond
	err = json.Unmarshal(([]byte)(respStr), &loginrespond)
	if err != nil {
		err = errors.New("BaiduLoginRespond json parse error")
		goto ERROR
	}
	fmt.Println("BaiduLoginRespond:",loginrespond)

	//TODO---解析  1，UrlEncode 2，Base64
	//判断ResultCode是否成功
	if loginrespond.ResultCode == "1" {
	//content URLEncode base64
	//先把content定义成string，通过urlEncode转，再base64 ，用解析json方式获得
	contentbytes, err := base64.URLEncoding.DecodeString(string(loginrespond.Content))
	if err != nil {
		err = errors.New("contentbytes Content json parse error")
		goto ERROR
	}
	fmt.Println("uid:",string(contentbytes))
	content := string(contentbytes)
	fmt.Println("content:",content)
	err = json.Unmarshal(([]byte)(contentbytes), &respond.Content)
	if err != nil {
		err = errors.New("BaiduQueryRespond Content json parse error")
		goto ERROR
	}
	fmt.Println("respond.Content:",respond.Content)
	//platformID,拿出content的UID
	buf = respond.Content.UID
	fmt.Println("UIDbuf:",buf)
	//respond MD5加密
	contentencode := base64.StdEncoding.EncodeToString(contentbytes)
	fmt.Println("content:",contentencode)
	md5content :=  loginrespond.AppID + loginrespond.ResultCode + contentencode + BaiduSecretKey
	fmt.Println("md5content:",md5content)
	loginsign := MD5([]byte(md5content))
	fmt.Println("loginsign:",loginsign)
	//sign validation
	if loginsign == loginrespond.Sign {
		fmt.Println("登陆请求成功")
		return
	}
	fmt.Println("登陆请求失败")
	goto ERROR
	}
	fmt.Println("ResultCode不等于1")
	goto ERROR*/


	type BaiduLoginRespond struct {
		AppID string `json:"AppID"`
		ResultCode string `json:"ResultCode"`
		ResultMsg string `json:"ResultMsg"`
		Sign string `json:"Sign"`
		Content string `json:"Content"`		
	}
	
	var  loginrespond BaiduLoginRespond
	
	err = json.Unmarshal(([]byte)(respStr), &loginrespond)
	if err != nil {
		err = errors.New("BaiduLoginRespond json parse error")
		goto ERROR
	}
	fmt.Println("BaiduLoginRespond:",loginrespond)

	//TODO---解析  1，UrlEncode 2，Base64
	/*urldecode,err := url.Parse(loginrespond.Content.UID)
	if err != nil {
		err = errors.New("urldecode json parse error")
		goto ERROR
	}
	//buf := uid 
	decode := urldecode.Query().Decode()
	fmt.Println("decode:",decode)*/
	//判断ResultCode是否成功
	if loginrespond.ResultCode == "1" {
	//content URLEncode base64
	//先把content定义成string，通过urlEncode转，再base64 ，用解析json方式获得
	urldecodestring ,err := url.QueryUnescape(string(loginrespond.Content))
	if err != nil {
		err = errors.New("urldecode error")
		goto ERROR
	}
	//url转小写字符
	urldecode := strings.ToLower(urldecodestring)
	contentbytes, err := base64.URLEncoding.DecodeString(urldecode)
	if err != nil {
		err = errors.New("contentbytes Content json parse error")
		goto ERROR
	}
	fmt.Println("uid:",string(contentbytes))
	type Content struct {
			UID string `json:"uid"`
	}
	var content Content
	err = json.Unmarshal(([]byte)(contentbytes), &content)
	if err != nil {
		err = errors.New("BaiduQueryRespond Content json parse error")
		goto ERROR
	}
	fmt.Println("Content json:",content)
	//platformID,拿出content的UID
	buf = content.UID
	fmt.Println("UIDbuf:",buf)
	//respond MD5加密
	contentencode := base64.StdEncoding.EncodeToString(contentbytes)
	fmt.Println("content:",contentencode)
	md5content :=  loginrespond.AppID + loginrespond.ResultCode + contentencode + BaiduSecretKey
	fmt.Println("md5content:",md5content)
	loginsign := MD5([]byte(md5content))
	fmt.Println("loginsign:",loginsign)
	//sign validation
	if loginsign == loginrespond.Sign {
		fmt.Println("登陆请求成功")
		type ServerRequest struct {
		AppID string `json:"AppID"`
		ResultMsg string `json:"ResultMsg"`
		Sign string `json:"Sign"`
		PlatformAccountID string `json:"platformAccountID"`		
		}
		var srequest ServerRequest
		srequest.AppID = loginrespond.AppID
		srequest.ResultMsg = loginrespond.ResultMsg
		srequest.Sign = loginrespond.Sign
		srequest.PlatformAccountID = buf + "@baidu"
		fmt.Println("srequest:",srequest)
	// add to parse
	jsons, err := json.Marshal(srequest)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(jsons)}}
	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(bodyBytes))
	defer response.Body.Close()
	c.String(http.StatusOK, string(bodyBytes))

	body, err := ioutil.ReadAll(resp.Body) 
	defer resp.Body.Close()
	fmt.Println(body)

	type ServerRespond struct{
		Result string `json:"result"`
		Message string `json:"message"`
	}
	var respond ServerRespond
	err = json.Unmarshal(bodyBytes, &respond)
	if err != nil {
		err = errors.New("ServerRespond json parse error")
		goto ERROR
	}
	fmt.Println("respond:",respond)
	//c.JSON(http.StatusOK, respond{Result: "SUCCESS", Message: ""})
	//c.JSON(http.StatusOK,respond{Result:respond.Result,Message:respond.Message})
		return
	}
	fmt.Println("登陆请求失败")
	goto ERROR
	}
	fmt.Println("ResultCode不等于1")
	goto ERROR
}

func BaiduOrderQuery(c *gin.Context) {
	//接收客户端的json数据
	type ThirdQueryRequest struct {
		AppID string `json:"appID"`
		//订单号
		CpOrderId  string `json:"cpOrderId "`		
	}
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	var tqrequest ThirdQueryRequest
	requestJson := c.PostForm("request")
	fmt.Println("requestJson:",requestJson)
	err = json.Unmarshal(([]byte)(requestJson), &tqrequest)
	if err != nil {
		err = errors.New("ThirdQueryRequest json parse error")
		goto ERROR
	}
	fmt.Println("tqrequest:",tqrequest)
	//发送订单信息查询
	type BaiduQueryRequest struct {
		AppID string `json:"AppID"`
		//订单号orderId
		CooperatorOrderSerial  string `json:"CooperatorOrderSerial "`
		OrderType int `json:"OrderType"`		
		Sign string `json:"Sign"`
		Action int `json:"Action "`
	}
	var request BaiduQueryRequest
	request.AppID = tqrequest.AppID
	request.CooperatorOrderSerial = tqrequest.CpOrderId
	request.OrderType = 1
	request.Action = 1002
	//订单查询签名sign
	signString := request.AppID + request.CooperatorOrderSerial + BaiduSecretKey
	fmt.Println("signString:",signString)
	//md5签名加密
	md5Sign := MD5([]byte(signString))
	//md5sign 变加密结果均转换为小写字符
	md5Sign = strings.ToLower(md5Sign)
	request.Sign = md5Sign

	fmt.Println("request:",request)

	// add to parse
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	fmt.Println("BaiduOrderQuery request jsons:",string(jsons))
	reader := bytes.NewBuffer(jsons)
	orderurl := "http://querysdkapi.91.com/CpOrderQuery.ashx"
	orderrequest, err := http.NewRequest("POST", orderurl, reader)
    if err != nil {
        goto ERROR
    }
    fmt.Println("order request:",orderrequest)
    orderrequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(orderrequest)
    if err != nil {
        goto ERROR
    }
    fmt.Println("resp:",resp)
    fmt.Println("resp.Body:",resp.Body)
    respBytes, err := ioutil.ReadAll(resp.Body)
    fmt.Println(respBytes)
    if err != nil {
        goto ERROR
    }
    respStr := string(respBytes)
    fmt.Println("respStr:",respStr)
    defer resp.Body.Close()
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("str:",*str)

	/*//接收订单信息查询
	type BaiduQueryRespond struct {
		AppID string `json:"AppID"`
		ResultCode string `json:"ResultCode"`
		ResultMsg string `json:"ResultMsg"`
		Sign string `json:"Sign"`
		Content struct {
			UID int64 `json:"UID"`
			OrderSerial  string `json:"OrderSerial "`
			MerchandiseName  string `json:"MerchandiseName "`
			OrderMoney float64  `json:"OrderMoney "`
			CooperatorOrderSerial string `json:"CooperatorOrderSerial"`
			OrderStatus int `json:"OrderStatus "`
			StatusMsg string `json:"StatusMsg "`
			StartDateTime string `json:"StartDateTime "`
			VoucherMoney int `json:"VoucherMoney "`
		}	`json:"Content"`
	}
	var respond BaiduQueryRespond
	err = json.Unmarshal(([]byte)(respStr), &respond)
	if err != nil {
		err = errors.New("respStr json parse error")
		goto ERROR
	}
	fmt.Println("respond:",respond)
	//content解析urlencode，再base64
	contentbytes, err := base64.URLEncoding.DecodeString(string(respond.Content))
	content := string(contentbytes)
	fmt.Println("contentbytes:",string(contentbytes))
	err = json.Unmarshal(([]byte)(contentbytes), &respond.Content)
	if err != nil {
		err = errors.New("respond Content json parse error")
		goto ERROR
	}
	fmt.Println("respond.Content:",respond.Content)
	//respond MD5加密
	//判断订单是否成功
	if respond.Content.OrderStatus == 1 {
		contentencode := base64.StdEncoding.EncodeToString(contentbytes)
		fmt.Println("content:",contentencode)
		md5content :=  respond.AppID + respond.ResultCode + contentencode + BaiduSecretKey
		fmt.Println("md5content:",md5content)
		ordersign := MD5([]byte(md5content))
		fmt.Println("ordersign:",ordersign)
		//判断签名是否一致
		if respond.Sign == ordersign {
		fmt.Println("支付请求成功")
		return
		}
	fmt.Println("支付请求失败")
	goto ERROR
	}
	fmt.Println("ResultCode不等于1")
	goto ERROR*/


	//接收订单信息查询
	type BaiduQueryRespond struct {
		AppID string `json:"appID"`
		ResultCode string `json:"resultCode"`
		ResultMsg string `json:"resultMsg"`
		Sign string `json:"sign"`
		Content string  `json:"content"`
	}
	type OrderContent struct {
			UID int64 `json:"uid"`
			OrderSerial  string `json:"orderSerial "`
			MerchandiseName  string `json:"merchandiseName "`
			OrderMoney float64  `json:"orderMoney "`
			CooperatorOrderSerial string `json:"cooperatorOrderSerial"`
			OrderStatus int `json:"orderStatus "`
			StatusMsg string `json:"statusMsg "`
			StartDateTime string `json:"startDateTime "`
			VoucherMoney int `json:"voucherMoney "`
		}
	var respond BaiduQueryRespond
	var ordercontent OrderContent
	err = json.Unmarshal(([]byte)(respStr), &respond)
	if err != nil {
		err = errors.New("respStr json parse error")
		goto ERROR
	}
	fmt.Println("respond:",respond)
	//content解析urlencode，再base64
	//添加了urldecode---不能确定是否正确
	urldecodestring ,err:= url.QueryUnescape(string(respond.Content))
	if err != nil {
		err = errors.New("urldecode error")
		goto ERROR
	}
	//url转小写字符
	urldecode := strings.ToLower(urldecodestring)
	//使用base64对url进行解码		
	contentbytes, err := base64.URLEncoding.DecodeString(urldecode)
	//变成json结构，有什么用
	err = json.Unmarshal(([]byte)(contentbytes), &ordercontent)
	if err != nil {
		err = errors.New("Content json parse error")
		goto ERROR
	}
	fmt.Println("Order Content:",ordercontent)
	//respond MD5加密
	//判断订单是否成功
	if ordercontent.OrderStatus == 1 {
		contentencode := base64.StdEncoding.EncodeToString(contentbytes)
		fmt.Println("content:",contentencode)
		md5content :=  respond.AppID + respond.ResultCode + contentencode + BaiduSecretKey
		fmt.Println("md5content:",md5content)
		ordersign := MD5([]byte(md5content))
		fmt.Println("ordersign:",ordersign)
		//判断签名是否一致
		if respond.Sign == ordersign {
		fmt.Println("支付请求成功")
		type ServerRequest struct {
		AppID string `json:"AppID"`
		ResultMsg string `json:"ResultMsg"`
		Sign string `json:"Sign"`
		Content struct {
			UID int64 `json:"uid"`
			OrderId  string `json:"orderId "`		//sdk内部订单号
			MerchandiseName  string `json:"merchandiseName "`
			OrderMoney float64  `json:"orderMoney "`
			CpOrderId string `json:"CpOrderId"`		//CP 订单号
			OrderStatus int `json:"orderStatus "`
			StatusMsg string `json:"statusMsg "`		//订单创建时描述
			StartDateTime string `json:"startDateTime "`	//订单创建时间 
			VoucherMoney int `json:"voucherMoney "`		//代金卷
		}	`json:"content"`
		PlatformAccountID string `json:"platformAccountID"`		
		}
		var srequest ServerRequest
		srequest.AppID = respond.AppID
		srequest.ResultMsg = respond.ResultMsg
		srequest.Sign = respond.Sign
		srequest.PlatformAccountID = buf + "@baidu"
		srequest.Content.UID = ordercontent.UID
		srequest.Content.OrderId = ordercontent.OrderSerial
		srequest.Content.MerchandiseName = ordercontent.MerchandiseName
		srequest.Content.OrderMoney = ordercontent.OrderMoney
		srequest.Content.CpOrderId = ordercontent.CooperatorOrderSerial
		srequest.Content.StatusMsg = ordercontent.StatusMsg
		srequest.Content.StartDateTime = ordercontent.StartDateTime
		srequest.Content.VoucherMoney = ordercontent.VoucherMoney

		fmt.Println("srequest:",srequest)
	// add to parse
	jsons, err := json.Marshal(srequest)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(jsons)}}
	response, err := http.PostForm("http://127.0.0.1:7777/accountQueryRequest", requestPkg)
	if err != nil {
		goto ERROR
	}
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(bodyBytes))
	defer response.Body.Close()
	//c.String(http.StatusOK, string(bodyBytes))

	body, err := ioutil.ReadAll(resp.Body) 
	defer resp.Body.Close()
	fmt.Println(body)

	type ServerRespond struct{
		Result string `json:"result"`
		Message string `json:"message"`
	}
	var respond ServerRespond
	err = json.Unmarshal(bodyBytes, &respond)
	if err != nil {
		err = errors.New("ServerRespond json parse error")
		goto ERROR
	}
	fmt.Println("respond:",respond)
	//c.String(http.StatusOK,respond)
	//不能返回数据，respond is no a type
	//c.JSON(http.StatusOK, respond{Result: "SUCCESS", Message: ""})
		return
		}
	fmt.Println("支付请求失败")
	goto ERROR
	}
	fmt.Println("ResultCode不等于1")
	goto ERROR
}

func BaiduPayReceive(c *gin.Context) {
	type BaiduPayRequest struct {
		AppID string `json:"appID"`
		OrderSerial  string `json:"orderSerial "`
		CooperatorOrderSerial string `json:"cooperatorOrderSerial"`
		Sign string `json:"sign"`
		Content string `json:"content"`
	}
	type BaiduPayContentRequest struct {
		UID int64 `json:"uid"`
		OrderSerial  string `json:"orderSerial "`
		MerchandiseName  string `json:"merchandiseName "`
		OrderMoney float64  `json:"orderMoney "`
		CooperatorOrderSerial string `json:"cooperatorOrderSerial"`
		OrderStatus int `json:"orderStatus "`
		StatusMsg string `json:"statusMsg "`
		StartDateTime string `json:"startDateTime "`
		VoucherMoney int `json:"voucherMoney "`
	}
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	var bprequest BaiduPayRequest
	var paycontentrequest BaiduPayContentRequest
	//receiveurl := "http://39.108.53.234:8686/"		//服务器地址
	requestJson := c.PostForm("request")
	//以request接收，直接使用base64解码
	err = json.Unmarshal([]byte(requestJson), &bprequest)
	if err != nil {
		err = errors.New("BaiduPayReceive request json parse error")
		goto ERROR
	}
	fmt.Println("BaiduPayReceive:",bprequest)
	contentbytes, err := base64.URLEncoding.DecodeString(bprequest.Content)
	fmt.Println("BaiduPay contentbytes:",contentbytes)
	//变成json结构，有什么用
	err = json.Unmarshal(contentbytes, &paycontentrequest)
	if err != nil {
		err = errors.New("Content json parse error")
		goto ERROR
	}
	fmt.Println("BaiduPay Content:",paycontentrequest)
	if paycontentrequest.OrderStatus == 1 {
		fmt.Println("OrderStatus == 1 , 支付请求成功")
		contentencode := base64.StdEncoding.EncodeToString(contentbytes)
		fmt.Println("content Encode:",contentencode)
		md5content :=  bprequest.AppID + bprequest.OrderSerial + bprequest.CooperatorOrderSerial + contentencode + BaiduSecretKey
		fmt.Println("md5content:",md5content)
		paysign := MD5([]byte(md5content))
		fmt.Println("Pay Sign:",paysign)
		if paysign == bprequest.Sign {
			fmt.Println("签名正确")
			type BaiduPayRespond struct {
				AppID string `json:"AppID"`
				ResultCode int `json:"ResultCode"`
				ResultMsg string `json:"ResultMsg"`
				Sign string `json:"Sign"`
				Content string `json:"Content"`
			}
			var respond BaiduPayRespond
			respond.AppID = bprequest.AppID
			respond.ResultCode = 1
			respond.ResultMsg = "发货通知成功接收"
			respond.Sign = paysign
			respond.Content = ""

			jsons, err := json.Marshal(respond)
			if err != nil {
			goto ERROR
			}
			requestPkg := url.Values{"request": {string(jsons)}}
			fmt.Println("requestPkg:",requestPkg)
			//c.JSON(http.StatusOK,requestPkg)
			
			//百度存在四个ip地址：59.56.20.61 、 59.56.17.125 、 36.250.11.72 、 36.250.11.87 
			//四个网段 ： 61.135.190.1-61.135.190.254 、 111.13.102.1-111.13.102.254 、
 			//		180.149.130.1-180.149.130.254 、 220.181.50.1-220.181.50.254 
			
			response, err := http.PostForm("http://59.56.20.61", requestPkg)
			//response, err := http.PostForm("http://127.0.0.1:7777/baiduPayReceive", requestPkg)
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
			fmt.Println("签名错误")
			goto ERROR
		}
	fmt.Println("OrderStatus != 1 , 支付请求失败")
	goto ERROR
}