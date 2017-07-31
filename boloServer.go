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
	"strconv"
	_ "strings"
	"io"	
	"time"
	//"regexp"
)
/*APPID：YjA2M2
秘钥：NzQ1YzAwMGZhNjg0MDM4*/
const (
	AppID = "YjA2M2"
	Key = "NzQ1YzAwMGZhNjg0MDM4"
	BoloURL = "https://i.fengei.com/api/5.0/user/verify"
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

	router.POST("/boloLogin", func(c *gin.Context) {
		BoloLogin(c)
	})
	router.POST("/boloPayNotify", func (c *gin.Context) {
		BoloPayNotify(c)
		})
	router.Run(":8686")
}

func BoloLogin(c *gin.Context) {
	type ThirdLoginRequest struct {
		ThirdLogin struct{
			OpenID string `json:"openid"`
			Token string `json:"token"`
			Platform string `json:"platform"`			
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
	v.Set("appid", AppID)
	v.Set("openid", loginrequest.ThirdLogin.OpenID)
	v.Set("token", loginrequest.ThirdLogin.Token)
	v.Set("platform", loginrequest.ThirdLogin.Platform)

	resp, err := http.Get(BoloURL + "?" + v.Encode())
	if err != nil {
		goto ERROR
	}
	fmt.Println("http.GET:",BoloURL + "?" + v.Encode())
   	fmt.Println("response:",resp)
    	fmt.Println("resp.Body:",resp.Body)
    	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto ERROR
	}
	respStr := string(respBytes)
	fmt.Println("respStr:",respStr)
	defer resp.Body.Close()
	type BoloRespond struct {
		Code int `json:"code"`
		OpenID string `json:"openid"`
		Reason string `json:"reason"`
	}
	var bolorespond BoloRespond
	//err = json.Unmarshal(([]byte)(respStr), &bolorespond)
	err = json.Unmarshal(respBytes, &bolorespond)
	if err != nil {
		err = errors.New("BoloRespond request json parse error")
		goto ERROR
	}
	fmt.Println("BoloRespond:",bolorespond)
	fmt.Println("结果码验证:",bolorespond.Code)
	if bolorespond.Code == 0 {
		fmt.Println("Bolo凭据有效")

		// add to parse
		loginrequest.PlatformAccountID = bolorespond.OpenID + "@bolo"
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
	}else if bolorespond.Code == -1 {
		fmt.Println("Bolo凭据无效")
		c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: "凭据无效"})
	}else if bolorespond.Code == -2 {
		fmt.Println("Bolo凭据对用户无效")
		c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: "凭据对用户无效"})
	}else if bolorespond.Code == -3 {
		fmt.Println("Bolo应用不存在")
		c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: "应用不存在"})
	}else if bolorespond.Code == -4 {
		fmt.Println("Bolo凭据过期")
		c.JSON(http.StatusOK,Respond{Result: "FAIL", Message: "凭据过期"})
	}
}

func BoloPayNotify(c *gin.Context) {
	type BoloPayRespond struct {
		Money string `json:"money"`
		No string `json:"no"`
		OrderId string `json:"orderid"`
		Params string `json:"params"`
		Time int `json:"time"`
		Sign string `json:"sign"`		
	}

	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())	
	c.String(http.StatusOK, "FAIL")
	return 

BEGIN:	
	fmt.Println("c.Request.Form:", c.Request.Form) 
	
	if c.Request.Form == nil { 
		c.Request.ParseMultipartForm(32 << 20) 
	} 
	for k, _ := range c.Request.Form { 
		fmt.Println("c.Request.Form.k:", k)


		var payrespond BoloPayRespond
		err = json.Unmarshal(([]byte)(k), &payrespond)
		if err != nil {
			err = errors.New("BoloPayRespond json parse error")
			goto ERROR
		}
		fmt.Println("BoloPayRespond:",payrespond)

		time := strconv.Itoa(payrespond.Time)
		sign := "money=" + payrespond.Money + "&no=" + payrespond.No + "&orderid=" + payrespond.OrderId + "&params=" + payrespond.Params + "&time=" + time
		secretsign := Key + "-" + sign + "-" + Key + "-" + AppID + "-" + Key
		fmt.Println("signstring:",sign)
		fmt.Println("secretsign:",secretsign)
		//md5 to hex
		md5sign := MD5([]byte(secretsign))
		fmt.Println("sign:",md5sign)
		fmt.Println("Sign:",payrespond.Sign)
		//sign verification
		if md5sign == payrespond.Sign {
			fmt.Println("数据返回正确，支付完成")		
			type ThirdRequest struct {
				PayFee string `json:"payFee"`
				CpOrderId string `json:"cpOrderId"`			
			}
			//third pay for request
			var request ThirdRequest
			//fee is float64
			Fee, err :=strconv.ParseFloat(payrespond.Money , 64)
			if err != nil {
				goto ERROR
			}
			//payFee is string
			request.PayFee = strconv.Itoa(int(Fee))
			fmt.Println("payfee:",request.PayFee)
			request.CpOrderId = payrespond.OrderId
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
			fmt.Println("response:",response)
			fmt.Println("bodyBytes:",string(bodyBytes))
			c.String(http.StatusOK,"SUCCESS")
		}		
	} 

}

func MD5(message []byte) string {
	w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //md5 to hex
	return md5str2
}