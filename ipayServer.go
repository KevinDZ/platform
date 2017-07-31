package main

import (
	"crypto/md5"
	_"database/sql"
	_"encoding/hex"
	"errors"
	"fmt"
	_ "strconv"
	_ "github.com/lib/pq"
	_ "bytes"
	"strings"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"net/http"
	"net/url"
	"time"
	"io"
	"io/ioutil"
)
const (
	IPayURL = "https://pay.ipaynow.cn" 

	//TODO:
	WeixinID      string = "148999970238152"
	WeixinAppKey  string = "ey12lW6mWurQAL73JMiSjZWj3T8p6rxK"
	
	//TODO:
	ZhifubaoID      string = "147868777472129"
	ZhifubaoAppKey  string = "1FZMAlAplOTamX6OARDVV8hrswhbGEVg"

)

const (
	funcode string = "WP001"
	notifyfuncode string = "N001"
	version string= "1.0.0"
	mhtOrderType string = "01"
	mhtCurrencyType string = "156"
	mhtCharset string = "UTF-8"
	deviceType string = "0601"
	mhtSignType string = "MD5"
	mhtOrderTimeOut string = "3600"
)



type ThirdRespond struct {
	Message string `json:"message"`
	Result string `json:"result"`
	PayOrderId string `json:"payOrderId"`
}

type ErrorRespond struct{
	Result string `json:"result"`
	Message string `json:"message"`
}

type UrlRespond struct {
	Result string `json:"result"`
	PayUrl string `json:"payUrl"`
}

func main() {
	gin.SetMode(GIN_MODE)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	ginUseLogger(router)

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
	
	router.POST("/iPayOrder", func (c *gin.Context) {
		IPayOrder(c)
		})
	router.POST("/iPayNotify", func (c *gin.Context) {
		IPayNotify(c)
		})

	router.Run(":11170")
}

func IPayOrder(c *gin.Context) {

type ThirdOrderRequest struct {
	PayID string `json:"payID"`
	UserID string `json:"userID"`
	PlatFormID string `json:"platformID"`

	DeviceInfo struct {
		Os string `json:"os"`
		Ios struct{
			Ip string `json:"ip"`
			Mac string `json:"mac"`
			DeviceType string `json:"deviceType"`
			Idfv string `json:"idfv"`

			Idfa string `json:"idfa"`
			} `json:"ios"`
		}	`json:"deviceInfo"`
		Timestamp string `json:"timestamp"`
		Fee string `json:"fee"`
		PayInfo string `json:"payInfo"`
		ChannelID string `json:"channelID"`
		ServerID string `json:"serverID"`

		Key string `json:"key"`

}

	type IPayRespond struct {
		Funcode string `json:"funcode"`
		Version string `json:"version"`
		AppId string `json:"appId"`
		ResponseCode string `json:"responseCode"`
		ResponseTime string `json:"responseTime"`
		ResponseMsg string `json:"responseMsg"`
		MhtOrderNo string `json:"mhtOrderNo"`
		NowPayOrderNo string `json:"nowPayOrderNo"`
		TransStatus string `json:"transStatus"`
		Tn string `json:"tn"`
		MhtSubMchId string `json:"mhtSubMchId"`
		SignType string `json:"signType"`
		Signature string `json:"signature"`
	}
	var err error
	goto BEGIN
ERROR:
	LogE(c, err.Error())
	fmt.Println(err)
	c.JSON(http.StatusOK, ErrorRespond{Result: "FAIL", Message: err.Error()})
	return
BEGIN:
	//var respond IPayRespond
	var thirdorderrequest ThirdOrderRequest
	var thirdrespond ThirdRespond

	//
	requestJson := c.PostForm("request")
	fmt.Println("IPay request:",requestJson)
	//fmt.Println("c.Request.Form:", c.Request.Form) 
	err = json.Unmarshal(([]byte)(requestJson), &thirdorderrequest)
	if err != nil {
		err = errors.New("ThirdOrderRequest json parse error")
		goto ERROR
	}
	fmt.Println("ThirdOrderRequest:",thirdorderrequest)

	requestPkg := url.Values{"request": {string(requestJson)}}
	response, err := http.PostForm("http://127.0.0.1:8090/ipaynowOrder", requestPkg)
	if err != nil {
		err = errors.New("http.PostForm response error")
		goto ERROR
	}
	//var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("ioutil.ReadAll(response.Body) error")
		goto ERROR
	}
	fmt.Println(string(bodyBytes))
	defer response.Body.Close()

	err = json.Unmarshal(bodyBytes, &thirdrespond)
	if err != nil {
		err = errors.New("ThirdPayRespond json parse error")
		goto ERROR
	}
	fmt.Println("ThirdRespond:",thirdrespond)
	if thirdrespond.Result != "SUCCESS" {
		err = errors.New("Result is FAIL")
		goto ERROR
	}	
	//Order info
	//times , err := strconv.Atoi(thirdorderrequest.Timestamp)
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	t := tm.Format("20060102150405")
	fmt.Println("t:",t)

	mhtOrderStartTime := string(t)
	fmt.Println("mhtOrderStartTime:",mhtOrderStartTime)
	outputType := "0"
	mhtOrderNo := thirdrespond.PayOrderId
	//request.MhtSubMchId = 
	mhtOrderName := thirdorderrequest.PayInfo
	mhtOrderAmt := thirdorderrequest.Fee
	mhtOrderDetail := thirdorderrequest.PayInfo
	notifyUrl := PAY_NOTIFY_HOST_URL + ":11170/iPayNotify"
	frontNotifyUrl := PAY_NOTIFY_HOST_URL + ":11170/iPayNotify"
	mhtReserved := "test"
	mhtOrderDetail = "test" 
	mhtOrderName = "test"
	var payChannelType ,appId , appKey string
	//如果是微信
	if thirdorderrequest.PayID ==  "23" {
		outputType = "2"
		payChannelType = "13"
		appId = WeixinID
		appKey = WeixinAppKey
		fmt.Println("WeixinAppKey:",WeixinAppKey)
	}else if thirdorderrequest.PayID ==  "24"{
		appId = ZhifubaoID
		payChannelType = "12"
		appKey = ZhifubaoAppKey
		fmt.Println("ZhifubaoAppKey:",ZhifubaoAppKey)
	}

	//if request.OutputType == "0" && request.PayChannelType == "11" {		
	//	request.PayAccNo = 	
	//}
	var consumerCreateIp string
	if thirdorderrequest.PayID ==  "24" {		
		consumerCreateIp = c.ClientIP()
		//consumerCreateIp = thirdorderrequest.DeviceInfo.Ios.Ip
		fmt.Println("ZhifubaoIp:",consumerCreateIp)
	}
	if outputType == "2" && payChannelType == "13" {		
		consumerCreateIp = c.ClientIP()
		//consumerCreateIp = thirdorderrequest.DeviceInfo.Ios.Ip
		fmt.Println("WeixinIp:",consumerCreateIp)
	}
	var md5Str string
	if thirdorderrequest.PayID ==  "23" {
	//IPayRequest sign for MD5
		md5Str = "appId=" + appId + "&consumerCreateIp=" + consumerCreateIp + "&deviceType=" + deviceType + "&frontNotifyUrl=" + frontNotifyUrl +"&funcode=" + funcode + "&mhtCharset=" + mhtCharset + 
			"&mhtCurrencyType=" + mhtCurrencyType + "&mhtOrderAmt=" +mhtOrderAmt + "&mhtOrderDetail=" + mhtOrderDetail + "&mhtOrderName=" + mhtOrderName + 
			"&mhtOrderNo=" + mhtOrderNo + "&mhtOrderStartTime=" + mhtOrderStartTime + "&mhtOrderTimeOut=" + mhtOrderTimeOut + "&mhtOrderType=" + 
			mhtOrderType + "&mhtReserved=" + mhtReserved + "&mhtSignType=" + mhtSignType + "&notifyUrl=" + notifyUrl + "&outputType=" + 
			outputType + "&payChannelType=" + payChannelType + "&version=" + version 		
	}else {
		md5Str = "appId=" + appId + "&deviceType=" + deviceType + "&frontNotifyUrl=" + frontNotifyUrl +"&funcode=" + funcode + "&mhtCharset=" + mhtCharset + 
			"&mhtCurrencyType=" + mhtCurrencyType + "&mhtOrderAmt=" +mhtOrderAmt + "&mhtOrderDetail=" + mhtOrderDetail + "&mhtOrderName=" + mhtOrderName + 
			"&mhtOrderNo=" + mhtOrderNo + "&mhtOrderStartTime=" + mhtOrderStartTime + "&mhtOrderTimeOut=" + mhtOrderTimeOut + "&mhtOrderType=" + 
			mhtOrderType + "&mhtReserved=" + mhtReserved + "&mhtSignType=" + mhtSignType + "&notifyUrl=" + notifyUrl + "&outputType=" + 
			outputType + "&payChannelType=" + payChannelType + "&version=" + version 	
	}
	fmt.Println("md5Str:",md5Str)		
	md5Key := MD5([]byte(appKey))
	fmt.Println("appKey:",appKey)
	fmt.Println("md5Key:",md5Key)
	signature := MD5([]byte(md5Str +  "&" + md5Key))
	fmt.Println("signature:",signature)
	fmt.Println("---------------------------------------------")
 	v := url.Values{}
	v.Set("appId", appId)
	v.Set("deviceType",  deviceType)
	v.Set("frontNotifyUrl", frontNotifyUrl)
	v.Set("funcode", funcode)
	v.Set("mhtCharset", mhtCharset )
	v.Set("mhtCurrencyType", mhtCurrencyType)
	v.Set("mhtOrderAmt", mhtOrderAmt)
	v.Set("mhtOrderDetail", mhtOrderDetail)
	v.Set("mhtOrderName", mhtOrderName)
	v.Set("mhtOrderNo", mhtOrderNo)
	v.Set("mhtOrderStartTime", mhtOrderStartTime)
	v.Set("mhtOrderTimeOut", mhtOrderTimeOut)
	v.Set("mhtOrderType", mhtOrderType)
	v.Set("mhtReserved", mhtReserved)
	v.Set("mhtSignType", mhtSignType)
	v.Set("notifyUrl", notifyUrl)
	v.Set("outputType", outputType)
	v.Set("payChannelType", payChannelType)
	v.Set("version", version)
	v.Set("mhtSignature", signature)
	if thirdorderrequest.PayID ==  "24" {
		fmt.Println("http.Get:",IPayURL + "?" + v.Encode())
		httpget := IPayURL + "?" + v.Encode()
		c.JSON(http.StatusOK,UrlRespond{Result: "SUCCESS",PayUrl: httpget})
		return 		
	}else{
	v.Set("consumerCreateIp", consumerCreateIp)
		var resp *http.Response
		resp, err = http.Get(IPayURL + "?" + v.Encode())
		if err != nil {
			err = errors.New("IPayResp parse error")
			goto ERROR
		}
		fmt.Println("resp:",resp)
		fmt.Println("resp.Body:",resp.Body)
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.New("ioutil.ReadAll error")
			goto ERROR
		}
		fmt.Println("respStr:",string(respBytes))
		respStr := string(respBytes)
		defer resp.Body.Close()
		respArray := strings.Split(respStr,"&")
		fmt.Println("respArray:",respArray)
		var responseCode string
		for _ ,v := range respArray {
			if strings.Contains(v, "responseCode=") {
				codeStr := strings.Split(v,"=")
				fmt.Println("codeStr:",codeStr[1])
				responseCode = codeStr[1]
			}
			if strings.Contains(v, "tn=") && responseCode == "A001" {
				fmt.Println("Weixin is tn")
				fmt.Println("v:",v)
				tnStr := strings.Split(v,"=")
				fmt.Println("tnStr:",tnStr[1])
				weixinTn ,err := url.QueryUnescape(tnStr[1])
				if err != nil {
					err = errors.New("url.QueryUnescape error")
					goto ERROR
				}
				fmt.Println("weixinTn:",weixinTn)
				c.JSON(http.StatusOK,UrlRespond{Result: "SUCCESS",PayUrl: weixinTn})				
				return
			}
		}
	}
	err = errors.New("ResponseCode != A001")
	goto ERROR
}

func IPayNotify(c *gin.Context) {
	type IPayNotifyRequest struct {
		Funcode string `json:"funcode"`
		Version string `json:"version"`
		AppId string `json:"appId"`
		//MhtSubMchId string `json:"mhtSubMchId"`
		MhtOrderNo string `json:"mhtOrderNo"`
		MhtOrderName string `json:"mhtOrderName"`
		MhtOrderType string `json:"mhtOrderType"`
		MhtCurrencyType string `json:"mhtCurrencyType"`
		MhtOrderAmt string `json:"mhtOrderAmt"`
		MhtOrderTimeOut string `json:"mhtOrderTimeOut"`
		MhtOrderStartTime string `json:"mhtOrderStartTime"`
		PayTime string `json:"payTime"`
		MhtCharset string `json:"mhtCharset"`
		NowPayOrderNo string `json:"nowPayOrderNo"`
		ChannelOrderNo string `json:"channelOrderNo"`
		DeviceType string `json:"deviceType"`
		PayChannelType string `json:"payChannelType"`
		TransStatus string `json:"transStatus"`
		PayConsumerId string `json:'payConsumerId'`
		MhtReserved string `json:"mhtReserved"`
		SignType string `json:"signType"`
		Signature string `json:"signature"`
	}

	type IPayNotifyRespond struct {
		Success string `json:"success"`
	}
	var err error
	goto BEGIN
ERROR:
	LogE(c, err.Error())
	c.JSON(http.StatusOK, ErrorRespond{Result: "FAIL", Message: err.Error()})
	return
BEGIN:
	var request IPayNotifyRequest

	
	buffer := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buffer)
	defer c.Request.Body.Close()
	body := string(buffer[:n])
	fmt.Println(string(buffer[0:n]))
	fmt.Println("body:",body)

	respArray := strings.Split(body,"&")
	fmt.Println("respArray:",respArray)
	for _ ,v := range respArray {		
		resultStr := strings.Split(v,"=")
		fmt.Println("resultStr:",resultStr[1])
		if strings.Contains(v, "channelOrderNo=") {
				channelOrderNo := strings.Split(v,"=")
				fmt.Println("channelOrderNo:",channelOrderNo[1])
				request.ChannelOrderNo = channelOrderNo[1]		
		}	
		if strings.Contains(v, "mhtOrderAmt=") {
				mhtOrderAmt := strings.Split(v,"=")
				fmt.Println("funcode:",mhtOrderAmt[1])
				request.MhtOrderAmt = mhtOrderAmt[1]				
		}
		if strings.Contains(v, "mhtOrderName=") {
				mhtOrderName := strings.Split(v,"=")
				fmt.Println("mhtOrderName:",mhtOrderName[1])
				request.MhtOrderName = mhtOrderName[1]				
		}
		if strings.Contains(v, "mhtOrderNo=") {
				mhtOrderNo := strings.Split(v,"=")
				fmt.Println("mhtOrderNo:",mhtOrderNo[1])
				request.MhtOrderNo = mhtOrderNo[1]				
		}
		if strings.Contains(v, "mhtOrderStartTime=") {
				mhtOrderStartTime := strings.Split(v,"=")
				fmt.Println("mhtOrderStartTime:",mhtOrderStartTime[1])
				request.MhtOrderStartTime = mhtOrderStartTime[1]				
		}
		if strings.Contains(v, "mhtReserved=") {
				mhtReserved := strings.Split(v,"=")
				fmt.Println("mhtReserved:",mhtReserved[1])
				request.MhtReserved = mhtReserved[1]				
		}
		if strings.Contains(v, "nowPayOrderNo=") {
				nowPayOrderNo := strings.Split(v,"=")
				fmt.Println("nowPayOrderNo:",nowPayOrderNo[1])
				request.NowPayOrderNo = nowPayOrderNo[1]				
		}		
		if strings.Contains(v, "payChannelType=") {
				payChannelType := strings.Split(v,"=")
				fmt.Println("payChannelType:",payChannelType[1])
				request.PayChannelType = payChannelType[1]				
		}
		if strings.Contains(v, "payConsumerId=") {
				payConsumerId := strings.Split(v,"=")
				fmt.Println("payConsumerId:",payConsumerId[1])
				request.PayConsumerId = payConsumerId[1]				
		}
		if strings.Contains(v, "payTime=") {
				payTime := strings.Split(v,"=")
				fmt.Println("payTime:",payTime[1])
				request.PayTime = payTime[1]				
		}
		if strings.Contains(v, "signature=") {
				signature := strings.Split(v,"=")
				fmt.Println("signature:",signature[1])
				request.Signature = signature[1]				
		}
		if strings.Contains(v, "transStatus=") {
				transStatus := strings.Split(v,"=")
				fmt.Println("transStatus:",transStatus[1])
				request.TransStatus = transStatus[1]				
		}
		if strings.Contains(v, "mhtOrderType=") {
				mhtOrderType := strings.Split(v,"=")
				fmt.Println("mhtOrderType:",mhtOrderType[1])
				request.MhtOrderType = mhtOrderType[1]				
		}
		if strings.Contains(v, "funcode=") {
				funcode := strings.Split(v,"=")
				fmt.Println("funcode:",funcode[1])
				request.Funcode = funcode[1]				
		}

	}
	//调用银联，退出
	var appID, appKey string
	if request.PayChannelType == "12" {
		appID, appKey = ZhifubaoID, ZhifubaoAppKey
		fmt.Println("Zhifubao Pay",appID,appKey)
	}else if request.PayChannelType == "13" {
		fmt.Println("Weixin Pay",appID,appKey)
		appID, appKey = WeixinID, WeixinAppKey
	}else if request.PayChannelType == "11" {
		fmt.Println("Bank Pay")
		return
	}


	//sign string
	//"appId=" + request.AppId
	signString := "appId=" + appID + "&channelOrderNo=" + request.ChannelOrderNo + "&deviceType=" + 
		deviceType + "&funcode=" + request.Funcode + "&mhtCharset=" + mhtCharset + 
		"&mhtCurrencyType=" + mhtCurrencyType +"&mhtOrderAmt=" + request.MhtOrderAmt + 
		"&mhtOrderName=" + request.MhtOrderName + "&mhtOrderNo=" + request.MhtOrderNo + 
		"&mhtOrderStartTime=" + request.MhtOrderStartTime + "&mhtOrderTimeOut=" + 
		mhtOrderTimeOut + "&mhtOrderType=" + request.MhtOrderType + "&mhtReserved=" + 
		request.MhtReserved + "&nowPayOrderNo=" + request.NowPayOrderNo+ "&payChannelType=" + 
		request.PayChannelType + "&payConsumerId=" + request.PayConsumerId +
		"&payTime=" + request.PayTime + "&signType=" + mhtSignType + "&transStatus=" + 
		request.TransStatus + "&version=" + version
	fmt.Println("Sign String:",signString)
	//md5
	md5Key := MD5([]byte(appKey))
	fmt.Println("appKey:",appKey)
	md5Str := signString + "&" + md5Key
	fmt.Println("md5Str:",md5Str)
	//md5
	signmd5 := MD5([]byte(md5Str))
	fmt.Println("sign:",signmd5)
	fmt.Println("Sign:",request.Signature)
	if signmd5 == request.Signature && request.TransStatus == "A001" {
		fmt.Println("sign pass")
		
		type ThirdRequest struct {
			PayFee string `json:"payFee"`
			CpOrderId string `json:"cpOrderId"`			
		}
		//third pay for request
		var thirdrequest ThirdRequest
			
		//payFee is string
		thirdrequest.PayFee =  request.MhtOrderAmt 
		fmt.Println("payfee:",thirdrequest.PayFee)
		thirdrequest.CpOrderId = request.MhtOrderNo
		fmt.Println("CpOrderId:",thirdrequest.CpOrderId)
		jsons, err := json.Marshal(thirdrequest)
		if err != nil {
			err = errors.New("thirdrequest json parse error")
			goto ERROR
		}
		fmt.Println("Third respond jsons:",string(jsons))

		requestPkg := url.Values{"request": {string(jsons)}}
		fmt.Println("requestPkg:",requestPkg)
					
		response, err := http.PostForm("http://127.0.0.1:8090/thirdPay", requestPkg)
		if err != nil {
			err = errors.New("http.PostForm error")
			goto ERROR
		}
		var bodyBytes []byte
		bodyBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			err = errors.New(" ioutil.ReadAll error")
			goto ERROR
		}
		defer response.Body.Close()
		fmt.Println("response:",response)
		fmt.Println("bodyBytes:",string(bodyBytes))

		c.JSON(http.StatusOK,IPayNotifyRespond{ Success: "Y"})
		return
	}
	c.JSON(http.StatusOK,IPayNotifyRespond{ Success: "N"})
	return
}

func MD5(message []byte) string {
	w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}