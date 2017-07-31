package main
import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/httputil"
	"net/url"
	_ "strings" 
	"strconv"
	"github.com/itsjamie/gin-cors"
	"time"
	"io/ioutil"
	"io"
	"regexp"
)

const (
	projectId string = "1"
	distributorId string = "1"
	path string = "./config/uc.conf"
	//apiKey string = "5b46e5be869c223a90f611fb01bd1e54"
	NotifyURL string = PAY_NOTIFY_HOST_URL + "/ucPayNotify"
	UCUrl string =  "http://sdk.9game.cn/cp/account.verifySession"
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

	type UC struct {
		ApiKey string `json:"apiKey"`
		//SecretKey string `json:"secretKey"`
	}
	var uc UC
	//获取配置文件的key

	keyJson := Pdp(projectId,distributorId ,path) 
	fmt.Println("key",keyJson)
	err := json.Unmarshal([]byte(keyJson),&uc)
	if err != nil {
		err = errors.New("UC Unmarshal error")
		fmt.Println(err)
		return
	}
	fmt.Println("uc",uc)
	apiKey := uc.ApiKey
	fmt.Println("apikey:",apiKey)


	// Apply the middleware to the router (works on groups too)
	router.Use(cors.Middleware(config))

	router.POST("/ucLogin", func(c *gin.Context) {
		UCThirdLogin(c, apiKey)
	})
	router.POST("/ucPay", func(c *gin.Context) {
		UCPay(c, apiKey)		
	})
	router.POST("/ucPayNotify", func (c *gin.Context) {
		UCPayNotify(c, apiKey)
		})
	router.Run(":9191")
}

func UCThirdLogin(c *gin.Context, apiKey string){

	type ThirdLoginRequest struct{	 
		ThirdLogin struct {
			Sid 	string 		`json:"sid"`	
			GameId 	string 	`json:"gameId"`
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
	}

	type Respond struct {
		Result     string `json:"result"`
		Message    string `json:"message"`
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
	fmt.Println(requestJson)
	var loginrequest ThirdLoginRequest

	err = json.Unmarshal(([]byte)(requestJson), &loginrequest)
	if err != nil {
		err = errors.New("tlrequest json parse error")
		goto ERROR
	}
	fmt.Println("UCRequest" , loginrequest)

	type UCLoginRequest struct{
	Id 	string 		`json:"id"`
	Data 	struct{
		Sid 	string 	`json:"sid"`	
	} 			`json:"data"`
	Game 	struct{
		GameId string 	`json:"gameId"`
	} 						`json:"game"`
	Sign  	string 	`json:"sign"`		
	}
	var ucrequest UCLoginRequest

	// id = time
	id := fmt.Sprint(time.Now().Unix())
	ucrequest.Id = id
	//Sid
	ucrequest.Data.Sid = loginrequest.ThirdLogin.Sid
	//gameId
	gameId := loginrequest.ThirdLogin.GameId
	ucrequest.Game.GameId = gameId
	//sign
	sign := "sid=" + loginrequest.ThirdLogin.Sid + apiKey
	fmt.Println("sign:",sign)
	ucrequest.Sign = MD5(([]byte)(sign))	
	fmt.Println("UCLoginRequest:",ucrequest)


	// add to parse
	jsons, err := json.Marshal(ucrequest)
	if err != nil {
		err = errors.New("UCLoginRequest json error")
		goto ERROR
	}
	fmt.Println("sdkrequest jsons:",string(jsons))
	reader := bytes.NewBuffer(jsons)
   
    newrequest, err := http.NewRequest("POST", UCUrl, reader)
    if err != nil {
    	err = errors.New("UCLoginRequest  new Request error")
	goto ERROR
    }
    fmt.Println("request:",newrequest)
    newrequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(newrequest)
    if err != nil {
    	err = errors.New("client do error")
	goto ERROR
    }
    fmt.Println("resp:",resp)
    fmt.Println("resp.Body:",resp.Body)
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	err = errors.New("respond readall error")
	goto ERROR
    }
    respStr := string(respBytes)
    fmt.Println("respStr:",respStr)
    defer resp.Body.Close()

    type LoginRequest struct{
	Id 		int64 	`json:"id"`
	State 	struct{
		Code int 			`json:"code"`
		Msg string 			`json:"msg"`
	}						`json:"state"`
	Data struct{
		//UC for account
		AccountId string 	`json:"accountId"`
		Creator string 		`json:"creator"`
		NickName string 	`json:"nickName"`
		} 			`json:"data"`
	}

	var lrequest LoginRequest
	err = json.Unmarshal([]byte(respStr), &lrequest)
	fmt.Println("tlrequest:",lrequest)
	if err != nil {
		err = errors.New("LoginRequest json parse error")
		goto ERROR
	}
	if lrequest.State.Code == 1 {		
	fmt.Println("respond success")	
	buf := lrequest.Data.AccountId
	fmt.Println("accountId:",buf)

	//userID
	type LoginRespond struct {
		Key string `json:"key"`
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
		//server  for UserID
		PlatformAccountID string `json:"platformAccountID"`
	}
	var loginresponds LoginRespond
	loginresponds.Key = loginrequest.GameAccessKey
	loginresponds.ChannelID = loginrequest.ChannelID
	loginresponds.PlatformID = loginrequest.PlatformID
	loginresponds.VersionID = loginrequest.VersionID
	loginresponds.DeviceInfo = loginrequest.DeviceInfo
	fmt.Println("loginrequest.DeviceInfo:",loginrequest.DeviceInfo)
	loginresponds.PlatformAccountID =  buf + "@uc"
	fmt.Println(loginresponds.PlatformAccountID)
	jsonstr, err := json.Marshal(loginresponds)
	if err != nil {
		err = errors.New("loginresponds error")
		goto ERROR
	}
	fmt.Println("loginrespond:",loginresponds)
	fmt.Println("send server jsons:",string(jsonstr))
	requestPkgs := url.Values{"request": {string(jsonstr)}}
	fmt.Println("requestPkg:",requestPkgs)

	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkgs)
	if err != nil {
		err = errors.New("PostForm error")
		goto ERROR
	}
	body1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("PostForm ReadAll error")
		goto ERROR
	}
	fmt.Println(response)
	fmt.Println(string(body1))
	defer response.Body.Close()


	//server resquest json 
	type AndrRespond struct {
		Result string `json:"result"`
		Message string `json:"message"`
		AccountId string 	`json:"accountId"`
		UserID string `json:"userID"`
	}
	var arespond AndrRespond
	err = json.Unmarshal(body1, &arespond)
	fmt.Println("request json:",arespond)
	if err != nil {
		err = errors.New("LoginRequest json parse error")
		goto ERROR
	}
	
	arespond.AccountId = lrequest.Data.AccountId
	fmt.Println("arespond.AccountId:",arespond.AccountId)
	var jsonbyte []byte
	jsonbyte, err = json.Marshal(arespond)
	if err != nil {
		err = errors.New("Marshal jsonbyte error")
		goto ERROR
	}
	c.String(http.StatusOK, string(jsonbyte))
	}else{
		err = errors.New("State.Code != 1")
		goto ERROR
	}
}

func UCPay(c *gin.Context, apiKey string){
	type ClientRequest struct {
		AccountId string `json:"accountId"`
		Amount string `json:"amount"`
		CallbackInfo string `json:"callbackInfo"`
		CpOrderId string `json:"cpOrderId"`
	}
	type PayRespond struct {
		Result     string `json:"result"`
		//PayOrderID string `json:"payOrderId"`
		//PayUrl     string `json:"payUrl"`
		Message    string `json:"message"`
	}
	var err error
	goto BEGIN
ERROR:
	LogE(c, err.Error())
	c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: err.Error()})
	return

BEGIN:

	clientrequestJson := c.PostForm("request")
	fmt.Println("requestJson:",clientrequestJson)
	var clientrequest ClientRequest
	
	err = json.Unmarshal([]byte(clientrequestJson), &clientrequest)
	if err != nil {
		err = errors.New("UCPay request json parse error")
		goto ERROR
	}
	fmt.Println("UCpay clientrequest:",clientrequest)
	//accountId
/*	fmt.Println("buf:",buf)
	AccountId := buf*/
	md5 := "accountId=" + clientrequest.AccountId + "amount=" + clientrequest.Amount + "callbackInfo=" + clientrequest.CallbackInfo + "cpOrderId=" + clientrequest.CpOrderId + "notifyUrl=" + NotifyURL + apiKey
	fmt.Println("AccountId",clientrequest.AccountId)
	fmt.Println("md5:",md5)


	md5Str := MD5(([]byte)(md5))
	fmt.Println("md5Str:",md5Str)
	
    type Respond struct {
			Result  string `json:"result"`
			Message string `json:"message"`
			AccountID string `json:"accountID"`			
			NotifyUrl string `json:"notifyUrl"`
			Sign string `json:"sign"`
			SignType string `json:"signType"`
		}
	// NotifyUrl、sign、signType
	c.JSON(http.StatusOK, Respond{Result: "SUCCESS", Message: "申请订单成功。 ",AccountID:  clientrequest.AccountId, NotifyUrl: NotifyURL, SignType: "MD5", Sign: md5Str})

	//}
}


func MD5(message []byte) string {
	w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
/*	md5 := ([]byte)(md5str2)
	fmt.Println(md5)
	fmt.Println(string(md5))*/
	return md5str2
}

func UCPayNotify(c *gin.Context, apiKey string){
	/*type UCCallBackRespond struct{
	Ver string 				`json:"ver"`
	Data struct{
		OrderId string 	 	`json:"orderId"`
		GameId string 	 	`json:"gameId"`
		AccountId string 	`json:"accountId"`
		Creator string  	`json:"creator"`
		PayWay string  	 	`json:"payWay"`
		Amount string 	 	`json:"amount"`
		CallBackInfo string `json:"callbackInfo"`
		OrderStatus string `json:"orderStatus"`
		FailedDesc string `json:"failedDesc"`
		CpOrderId int32 `json:"cporderId"`
	} `json:"data"`
	Sign string `json:"sign"`
}*/
	type UCCallBackRequest struct{
		Result string `json:"result"`
		Message string `json:"message"`
		// for account
		//PlatformAccountID string `json:"platformAccountID"`
	}
	var err error
	goto BEGIN

ERROR:
	LogE(c, err.Error())
	c.JSON(http.StatusOK, UCCallBackRequest{Result: "FAILURE"})
	return 
BEGIN:	
	//UCcallback send ucServer 
	type ServerRequest struct {
	Ver string 				`json:"ver"`
	Data struct{
		OrderId string 	 	`json:"orderId"`
		GameId string 	 	`json:"gameId"`
		AccountId string 	`json:"accountId"`
		Creator string  	`json:"creator"`
		PayWay string  	 	`json:"payWay"`
		Amount string 	 	`json:"amount"`
		CallBackInfo string `json:"callbackInfo"`
		OrderStatus string `json:"orderStatus"`
		FailedDesc string `json:"failedDesc"`
		CpOrderId string `json:"cporderId"`
	} `json:"data"`
	Sign string `json:"sign"`
	// for account
	//PlatformAccountID string `json:"platformAccountID"`
	}
	var srequest ServerRequest
	//获取UC数据
	buffer := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buffer)
	defer c.Request.Body.Close()
	body := string(buffer[:n])
	fmt.Println(string(buffer[0:n]))
	bodystring := string(body)
	fmt.Println("bodystring:",bodystring)
	//regexp  
	reg := regexp.MustCompile(`{(.*)}`)
	result := reg.FindAllStringSubmatch(bodystring,-1)
	//var bodyreg []string
	for _,v := range result {
		fmt.Println("result:",v[1])
	
	requeststring := "{" + v[1] + "}"
	fmt.Println("requeststring:",requeststring)
	err = json.Unmarshal(([]byte)(requeststring), &srequest)
	if err != nil {
		err = errors.New("ServerRequet json parse error")
		goto ERROR
	}
	/*srequest.PlatformAccountID = buf +"@uc"*/
	fmt.Println("srequest:",srequest)

	//orderstatus == s
	if srequest.Data.OrderStatus == "S" {
			fmt.Println("OrderStatus is S")

			type ServerRespond struct {
				CpOrderId string `json:"cpOrderId"`
				PayFee string `json:"payFee"`
				PayTime string `json:"payTime"`
			}
			fmt.Println("amountFee:",srequest.Data.Amount)			
			fmt.Println("srequest:",srequest)
			// add to parse
			var srespond ServerRespond
			
			
			srequest.Data.FailedDesc = ""
			failedDesc := "failedDesc=" + srequest.Data.FailedDesc
			accountId := "accountId=" + srequest.Data.AccountId
			amount := "amount=" + srequest.Data.Amount
			fmt.Println("amount:",srequest.Data.Amount)
			callbackInfo := "callbackInfo=" + srequest.Data.CallBackInfo
			creator := "creator=" + srequest.Data.Creator
			gameId := "gameId=" + srequest.Data.GameId
			orderStatus := "orderStatus=" + srequest.Data.OrderStatus
			payWay := "payWay=" + srequest.Data.PayWay
			orderId := "orderId=" + srequest.Data.OrderId
			cpOrderID := string(srequest.Data.CpOrderId)

			cpOrderId := "cpOrderId=" + cpOrderID				
			md5 := accountId + amount + callbackInfo + cpOrderId + creator + failedDesc + gameId + orderId+ orderStatus + payWay + apiKey

			fmt.Println("md5:",md5)
			sign := MD5(([]byte)(md5))
			fmt.Println("sign:",sign)
			fmt.Println("srequest.Sign:",srequest.Sign)	
			if srequest.Sign == sign {	
				fmt.Println("sign pass")
				//amount * 100
				amountFee ,err:= strconv.ParseFloat(srequest.Data.Amount, 64)
				if err != nil {
					err = errors.New("strconv.Atoi error")
					goto ERROR
				}
				fee := amountFee*100
				fmt.Println("fee:",fee)
				srespond.PayFee = strconv.Itoa(int(fee))
				srespond.CpOrderId = srequest.Data.CpOrderId
				fmt.Println("srequest.Data.Amount:",srequest.Data.Amount)
				jsons, err := json.Marshal(srespond)
				if err != nil {
					err = errors.New("json.Marsha error")
					goto ERROR
				}
				fmt.Println("server respond jsons:",string(jsons))

				fmt.Println("srespond:",srespond)				
				requestPkg := url.Values{"request": {string(jsons)}}
				fmt.Println("requestPkg:",requestPkg)
				response, err := http.PostForm("http://127.0.0.1:8090/thirdPay", requestPkg)
				if err != nil {
					err = errors.New("PostForm error")
					goto ERROR
				}
				var bodyBytes []byte
				bodyBytes, err = ioutil.ReadAll(response.Body)
				if err != nil {
					err = errors.New("PostForm readall error")
					goto ERROR
				}
				defer response.Body.Close()
				fmt.Println(string(bodyBytes))
				c.String(http.StatusOK, "SUCCESS")
				return 
			}
		}
	}
	c.String(http.StatusOK, "FAILURE")
	return 
}