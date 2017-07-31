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
	"time"
	"io/ioutil"
	"io"
	"unsafe"
	"regexp"
)

const (
	apiKey string = "5b46e5be869c223a90f611fb01bd1e54"
)
var (
	buf string  //存储AccountId的数据
	gameId string
)

func main() {
	router := gin.Default()

	router.POST("/ucLogin", func(c *gin.Context) {
		ThirdLogin(c)
	})
	router.POST("/ucPay", func(c *gin.Context) {
		UCPay(c)		
	})
	router.POST("/ucPayNotify", func (c *gin.Context) {
		UCPayCallback(c)
		})
	router.Run(":9191")
}

func ThirdLogin(c *gin.Context){

	type ThirdLoginRequest struct{	 
	 	ThirdLogin	struct{
		Sid 	string 		`json:"sid"`	
		GameId 	string 	`json:"gameId"`
		} 			`json:"thirdLogin"`	
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
	var request ThirdLoginRequest

	err = json.Unmarshal(([]byte)(requestJson), &request)
	if err != nil {
		err = errors.New("tlrequest json parse error")
		goto ERROR
		//fmt.Println(err.Error())
	}
	fmt.Println("UCRequest" , request)

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
	ucrequest.Data.Sid = request.ThirdLogin.Sid
	//gameId
	gameId = request.ThirdLogin.GameId
	ucrequest.Game.GameId = gameId
	//sign
	sign := "sid=" + request.ThirdLogin.Sid + apiKey
	fmt.Println("sign:",sign)
	ucrequest.Sign = MD5(([]byte)(sign))
	
	fmt.Println("UCLoginRequest:",ucrequest)

	// add to parse
	jsons, err := json.Marshal(ucrequest)
	if err != nil {
		goto ERROR
	}
	fmt.Println("sdkrequest jsons:",string(jsons))
	reader := bytes.NewBuffer(jsons)
    sdkurl := "http://sdk.9game.cn/cp/account.verifySession"
    newrequest, err := http.NewRequest("POST", sdkurl, reader)
    if err != nil {
        goto ERROR
    }
    fmt.Println("request:",newrequest)
    newrequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(newrequest)
    if err != nil {
        goto ERROR
    }
    fmt.Println("resp:",resp)
    fmt.Println("resp.Body:",resp.Body)
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        goto ERROR
    }
    respStr := string(respBytes)
    fmt.Println("respStr:",respStr)
    defer resp.Body.Close()
    str := (*string)(unsafe.Pointer(&respBytes))
    fmt.Println("str:",*str)

    type LoginRequest struct{
	Id 		int64 	`json:"id"`
	State 	struct{
		Code int 			`json:"code"`
		Msg string 			`json:"msg"`
	}						`json:"state"`
	Data struct{
		AccountId string 	`json:"accountId"`
		Creator string 		`json:"creator"`
		NickName string 	`json:"nickName"`
		} 			`json:"data"`
	// for account
	PlatformAccountID string `json:"platformAccountID"`
	}
	requestJsons1 := c.PostForm("request")

	var tlrequest LoginRequest

	err = json.Unmarshal([]byte(requestJsons1), &tlrequest)
	fmt.Println("tlrequest:",tlrequest)
	if err != nil {
		err = errors.New("LoginRequest json parse error")
		goto ERROR
	}

	buf := tlrequest.Data.AccountId
	fmt.Println("accountId:",buf)
	//c.String(http.StatusOK,"UserID: %s", buf)

	tlrequest.PlatformAccountID =  buf + "@uc"
	fmt.Println(tlrequest.PlatformAccountID)
	jsonstr, err := json.Marshal(tlrequest)
	if err != nil {
		goto ERROR
	}
	fmt.Println(tlrequest)
	requestPkg := url.Values{"request": {string(jsonstr)}}
	fmt.Println(string(jsons))

	response, err := http.PostForm("http://127.0.0.1:7777/accountThirdLogin", requestPkg)
	//response, err := http.PostForm("http://test.tanchenggame.com:7777/accountThirdLogin", requestPkg)
	if err != nil {
		goto ERROR
	}
	body1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(response)
	fmt.Println(string(body1))
	defer response.Body.Close()
	c.String(http.StatusOK, string(body1))

}

func UCPay(c *gin.Context){
	type ClientRequest struct {
		Amount string `json:"amount"`
		CallbackInfo string `json:"callbackInfo"`
		CpOrderId string `json:"cpOrderId"`
		NotifyUrl string `json:"notifyUrl"`

	}
	type PayRespond struct {
		Result     string `json:"result"`
		PayOrderID string `json:"payOrderId"`
		PayUrl     string `json:"payUrl"`
		Message    string `json:"message"`
	}
	var err error
	goto BEGIN
ERROR:
	LogE(c, err.Error())
	c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: "申请订单失败。", PayUrl: ""})
	return

BEGIN:
	requestJson := c.PostForm("request")
	fmt.Println(requestJson)
	var request ClientRequest
	//var ysdkPayRequest YSDKPayRequest

	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		err = errors.New("request json parse error")
		goto ERROR
	}
	//accountId
	AccountId := buf
	//端口号：4040
	notifyUrl := "http://120.77.202.87:4040/ucpaycallback"
	//notifyUrl := "http://172.16.60.61/ucpaycallback"
	md5 := "accountId=" + AccountId + "amount=" + request.Amount + "callbackInfo=" + request.CallbackInfo + "cpOrderId=" + request.CpOrderId + "notifyUrl=" + notifyUrl + apiKey
	fmt.Println("md5:",md5)


	md5Str := MD5(([]byte)(md5))
	fmt.Println("md5Str:",md5Str)
	
    type Respond struct {
			Result  string `json:"result"`
			Message string `json:"message"`
			AccountID string `json:"accountID"`
			//TODO  增加了notifyurl作为sdk服务器的回调方法参数
			NotifyUrl string `json:"notifyUrl"`
			//TODO  还缺少sign和signType
			Sign string `json:"sign"`
			SignType string `json:"signType"`
		}

		//添加了NotifyUrl、sign、signType
	c.JSON(http.StatusOK, Respond{Result: "SUCCESS", Message: "申请订单成功。 ",AccountID: AccountId, NotifyUrl: notifyUrl, SignType: "MD5", Sign: md5Str})

}


func MD5(message []byte) string {
	w := md5.New()
	msg := string(message)
	io.WriteString(w, msg)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
	md5 := ([]byte)(md5str2)
	fmt.Println(md5)
	fmt.Println(string(md5))
	return string(md5)
}

func UCPayCallback(c *gin.Context) string{
	type UCCallBackRespond struct{
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
}
	type UCCallBackRequest struct{
		Result string `json:"result"`
		Message string `json:"message"`
	}
	var err error
	goto BEGIN

ERROR:

	LogE(c, err.Error())
	c.JSON(http.StatusOK, UCCallBackRequest{Result: "FAILURE"})
	return "FAILURE"

BEGIN:
	//TODO 从sdk服务端获取POST请求
	sdkurllogin := "http://sdk.9game.cn/cp/account.verifySession"
	resp, err := http.Post(sdkurllogin, "application/json;charset=utf-8", nil)
	if err != nil { 
      goto ERROR 
}
	fmt.Println(resp.Body)
	
	//UCcallback的数据发送到服务器
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
	PlatformAccountID string `json:"platformAccountID"`
	}
	var srequest ServerRequest
	//获取UC数据
	buffer := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buffer)
	defer c.Request.Body.Close()
	body := string(buffer[:n])
	fmt.Println("c:",c)
	fmt.Println(string(buffer[0:n]))
	fmt.Println("body:",body)
	//转成字符串处理
	bodystring := string(body)
	fmt.Println("bodystring:",bodystring)
	//regexp  正则表达式处理
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
	srequest.PlatformAccountID = buf +"@uc"
	fmt.Println("srequest:",srequest)
	// add to parse
	jsons, err := json.Marshal(srequest)
	if err != nil {
		goto ERROR
	}
	requestPkg := url.Values{"request": {string(jsons)}}
	response, err := http.PostForm("http://127.0.0.1:7777/ThirdPay", requestPkg)
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
	var respond UCCallBackRespond
	err = json.Unmarshal(bodyBytes, &respond)
	if err != nil {
		err = errors.New("UCCallBackRespond json parse error")
		goto ERROR
	}
	fmt.Println("Callback:",respond)
	//判断订单状态
	if respond.Data.OrderStatus == "S" {
			respond.Data.FailedDesc = ""
			failedDesc := respond.Data.FailedDesc
			accountId := "accountId=" + respond.Data.AccountId
			amount := "amount=" + respond.Data.Amount
			callbackInfo := "callBackInfo=" + respond.Data.CallBackInfo
			creator := "creator=" + respond.Data.Creator
			gameId := "gameId=" + respond.Data.GameId
			orderStatus := "orderStatus" + respond.Data.OrderStatus
			payWay := "payWay" + respond.Data.PayWay
			orderId := "orderId" + respond.Data.OrderId
			cpOrderID := string(respond.Data.CpOrderId)
			if cpOrderID != "" {
				cpOrderId := "cpOrderId=" + cpOrderID				
				md5 := accountId + amount + callbackInfo + cpOrderId + creator + failedDesc + gameId + orderId+ orderStatus + payWay + apiKey
				sign := MD5(([]byte)(md5))
			respond.Sign = sign
			fmt.Println(sign)
			fmt.Println(respond.Sign)	
			c.JSON(http.StatusOK, UCCallBackRequest{Result: "SUCCESS"})
			return "SUCCESS"
			}else{
				md5 := accountId + amount + callbackInfo + creator + failedDesc + gameId + orderId+ orderStatus + payWay + apiKey
			sign := MD5(([]byte)(md5))
			respond.Sign = sign
			fmt.Println(sign)
			fmt.Println(respond.Sign)	
			c.JSON(http.StatusOK, UCCallBackRequest{Result: "SUCCESS"})
			return "SUCCESS"	
			}
		}
	goto ERROR
	}
	return "FAILURE"
}