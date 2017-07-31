package main
import (
	"bytes"
	//"crypto/hmac"
	"crypto/md5"
	//"crypto/sha1"
	//_ "crypto/sha256"
	//"database/sql"
	//"encoding/base64"
	//"encoding/base32"
	//"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	//"io/ioutil"
	"net/http"
	//"net/http/cookiejar"
	_ "net/http/httputil"
	"net/url"
	//"strconv"
	_ "strings" 
	"time"
	//"os"
	"io/ioutil"
	"io"
	"unsafe"
	"regexp"
	//"log"
)
/*
func main(){
	router := gin.Default()
	router.POST("/post",func (this *gin.Context){
		id := this.Query("id")
		page := this.DefaultQuery("page","0")
		name := this.PostForm("name")
		message := this.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		})
	router.Run(":4400")
}*/

const (
	apiKey string = "5b46e5be869c223a90f611fb01bd1e54"
)
var (
	buf string  //存储AccountId的数据
	gameId string
)

func main() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		ThirdLogin(c)
	})

	router.POST("/ucpay", func(c *gin.Context) {
		UCPay(c)		
	})
	router.POST("/ucpaycallback", func (c *gin.Context) {
		UCPayCallback(c)
		})

	router.POST("/submitroledata", func (c *gin.Context) {
		SubmitRoleData(c)
		})

	router.Run(":4040")
}
//未进行判断创建者：JY、ALI、PP、WDJ
func ThirdLogin(c *gin.Context){

	type ThirdLoginRequest struct{	 
	 	ThirdLogin	struct{
		Sid 	string 		`json:"sid"`	
		GameId 	string 		`json:"gameId"`
		} 					`json:"thirdLogin"`	
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
	var tlrequest ThirdLoginRequest

	err = json.Unmarshal(([]byte)(requestJson), &tlrequest)
	if err != nil {
		err = errors.New("tlrequest json parse error")
		goto ERROR
		//fmt.Println(err.Error())
	}
	fmt.Println("-----" , tlrequest , "-------")

	type SDKLoginRequest struct{
	Id 			string 		`json:"id"`
	Data 		struct{
		Sid 	string 		`json:"sid"`	
	} 						`json:"data"`
	Game 		struct{
		GameId 	string 		`json:"gameId"`
	} 						`json:"game"`
	Sign  		string 		`json:"sign"`		
	}
	var sdkrequest SDKLoginRequest
	// id
	id := fmt.Sprint(time.Now().Unix())
	sdkrequest.Id = id
	//Sid
	sdkrequest.Data.Sid = tlrequest.ThirdLogin.Sid
	//gameId
	gameId = tlrequest.ThirdLogin.GameId
	sdkrequest.Game.GameId = gameId
	//sign
	sign := "sid=" + tlrequest.ThirdLogin.Sid + apiKey
	fmt.Println("sign:",sign)
	sdkrequest.Sign = MD5(([]byte)(sign))
	
	fmt.Println("sdkrequest:",sdkrequest)

	// add to parse
	jsons, err := json.Marshal(sdkrequest)
	if err != nil {
		goto ERROR
	}
	fmt.Println("sdkrequest jsons:",string(jsons))
	reader := bytes.NewBuffer(jsons)
    sdkurl := "http://sdk.9game.cn/cp/account.verifySession"
    //测试完全不能用--直接报解析错误
    //sdkurl := "http://sdk.test4.g.uc.cn/cp/account.verifySession"
    request, err := http.NewRequest("POST", sdkurl, reader)
    if err != nil {
        goto ERROR
    }
    fmt.Println("request:",request)
    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        goto ERROR
    }
    fmt.Println("resp:",resp)
    fmt.Println("resp.Body:",resp.Body)
    respBytes, err := ioutil.ReadAll(resp.Body)
    //fmt.Println(respBytes)
    if err != nil {
        goto ERROR
    }
    respStr := string(respBytes)
    fmt.Println("respStr:",respStr)
    defer resp.Body.Close()
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("str:",*str)
	
	type sdkAccountParam struct{
	Id 		int64 	`json:"id"`
	State 	struct{
		Code int 			`json:"code"`
		Msg string 			`json:"msg"`
	}						`json:"state"`
	Data struct{
		AccountId string 	`json:"accountId"`
		Creator string 		`json:"creator"`
		NickName string 	`json:"nickName"`
		} 					`json:"data"`
}
	var accountParam sdkAccountParam
	//fmt.Println(respStr)
	///*
	err = json.Unmarshal(([]byte)(respStr), &accountParam)
	if err != nil {
		err = errors.New("accountParam request json parse error")
		goto ERROR
	}
	
	accountId := accountParam.Data.AccountId
	fmt.Println("accountId:",accountId)
	buf = accountId
	//c.JSON(http.StatusOK, accountParam{Data:{AccountId:accountId}})


	/*
	requestPkg := url.Values{"request": {string(json)}}
	response, err := http.PostForm("http://sdk.9game.cn/cp/account.verifySession", requestPkg)
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
	*/

	//目前不需要
	/*
	type SDKLoginResponse struct{
	Id 		int64 	`json:"id"`
	//响应码说明（state.code）
	//	
	//响应码	说明	错误原因
	//1	成功	--
	//10	请求参数错误	请求内容格式有误、gameID有误、签名校验失败等
	//11	用户未登录	sid不存在，请求地址有误等
	//99	服务器内部错误	接口名有误、请求地址有误等
	// 
	State 	struct{
		Code int 			`json:"code"`
		Msg string 			`json:"msg"`
	}						02`json:"state"`
	Data struct{02
		AccountId string 	02`json:"accountId"`
		Creator string 		02`json:"creator"`
		NickName string 	`json:"nickName"`
		} 					`json:"data"`
}

	var sdkresponse SDKLoginResponse*/
	//return ""

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
	//fmt.Println(request)
	//accountId
	AccountId := buf
	//端口号：4040
	notifyUrl := "http://120.77.202.87:4040/ucpaycallback"
	//notifyUrl := "http://172.16.60.61/ucpaycallback"
	md5 := "accountId=" + AccountId + "amount=" + request.Amount + "callbackInfo=" + request.CallbackInfo + "cpOrderId=" + request.CpOrderId + "notifyUrl=" + notifyUrl + apiKey
	fmt.Println("md5:",md5)

	/*err = json.Unmarshal([]byte(requestJson), &ysdkPayRequest)
	if err != nil {
		err = errors.New("request json parse error")
		goto ERROR
	}
	*/

	md5Str := MD5(([]byte)(md5))
	fmt.Println("md5Str:",md5Str)
	/*
	jsonstr, err := json.Marshal(md5Str)
	if err != nil {
		goto ERROR
	}
	jsonStr := string(jsonstr)
	*/

	/*reader := bytes.NewBuffer(json)
    resp, err := http.Post(url, "application/x-www-form-urlencoded", reader)
    if err != nil {
        goto ERROR
    }

    client := &http.Client{

    }
    resp, err := client.Do(request)
    if err != nil {
        goto ERROR
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        goto ERROR
    }
    defer resp.Body.Close()
	str := (*string)(unsafe.Pointer(&respBytes))
    fmt.Println(*str)*/
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
	//c.JSON(http.StatusOK, Respond{Result: "SUCCESS", Message: "申请订单成功。 ",AccountID: AccountId, NotifyUrl: notifyUrl, SignType: "MD5", Sign: jsonStr})
    //c.JSON(http.StatusOK, Respond{Result: "SUCCESS", Message: fmt.Sprintf(string(json))})
/*
	//pay success
	err = SaveOrderAndNotifyCP(gameIDMap, payOrderTable, payOrderID)
	if err != nil {
		goto ERROR
	}
	return nil*/

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
	//requestJson := c.PostForm("request")
	//var request UCCallBackRequest
	//fmt.Println("callBack requestJson:",requestJson)
	/*//设置时间戳
	timestamp := fmt.Sprint(time.Now().Unix())
	*/
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
	
	//UCcallback的数据发送到服务器
	/*type ServerRequest struct {
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
	// for account
	PlatformAccountID string `json:"platformAccountID"`
	}*/
	//v[1] join {}
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

func SubmitRoleData(c *gin.Context) {
/*
	type UCSubmitRoleDataRequest struct{
	Id string `json:"id"`
	Service string `json:"service"`
	Data struct{
		AccountId string `json:"accountId"`
		
		//假如gameData的数据为:
		//{"category":"loginGameRole","content":{"roleLevel":"88","roleName":" 请∝ 再给我一支烟","zoneName":"终南山下-兵临城下","roleId":"53568193","zoneId":2705,"roleCTime":1353271378,"os":"android","roleLevelMTime":1456380919}}
		//那么经过UrlEncode后的字符串为（注：该字符串也是作为计算MD5签名的原文内容）
		//%7b%22category%22%3a%22loginGameRole%22%2c%22content%22%3a%7b%22roleLevel%22%3a%2288%22%2c%22roleName%22%3a%22+%e8%af%b7%e2%88%9d+%e5%86%8d%e7%bb%99%e6%88%91%e4%b8%80%e6%94%af%e7%83%9f%22%2c%22zoneName%22%3a%22%e7%bb%88%e5%8d%97%e5%b1%b1%e4%b8%8b-%e5%85%b5%e4%b8%b4%e5%9f%8e%e4%b8%8b%22%2c%22roleId%22%3a%2253568193%22%2c%22zoneId%22%3a2705%2c%22roleCTime%22%3a1353271378%2c%22os%22%3a%22android%22%2c%22roleLevelMTime%22%3a1456380919%7d%7d
		 
		GameData string `json:"gameData"`
		}`json:"data"`
	Game struct{
		GameId string `json:"gameId"`
		}`json:"game"`
	sign string `json:"sign"`
	}

	

	
	//1.gameData原本是json对象
	//2.gameData经过json encode成为一个字符串
	//3.该字符串再经过urlencode

	//假如gameData的数据为:
	//{"category":"loginGameRole","content":{"roleLevel":"88","roleName":" 请∝ 再给我一支烟","zoneName":"终南山下-兵临城下","roleId":"53568193","zoneId":2705,"roleCTime":1353271378,"os":"android","roleLevelMTime":1456380919}}
	//那么经过UrlEncode后的字符串为（注：该字符串也是作为计算MD5签名的原文内容）---编码base64
	//%7b%22category%22%3a%22loginGameRole%22%2c%22content%22%3a%7b%22roleLevel%22%3a%2288%22%2c%22roleName%22%3a%22+%e8%af%b7%e2%88%9d+%e5%86%8d%e7%bb%99%e6%88%91%e4%b8%80%e6%94%af%e7%83%9f%22%2c%22zoneName%22%3a%22%e7%bb%88%e5%8d%97%e5%b1%b1%e4%b8%8b-%e5%85%b5%e4%b8%b4%e5%9f%8e%e4%b8%8b%22%2c%22roleId%22%3a%2253568193%22%2c%22zoneId%22%3a2705%2c%22roleCTime%22%3a1353271378%2c%22os%22%3a%22android%22%2c%22roleLevelMTime%22%3a1456380919%7d%7d
		
	type GameData struct{
	Category string `json:"category"`
	Content  struct{ 
		RoleLevel 		string  `json:"roleLevel"`
		RoleName 		string  `json:"roleName"`
		ZoneName 		string  `json:"zoneName"`
		RoleId 			string  `json:"roleId"`
		ZoneId 			string  `json:"zoneId"`
		RoleCTime 		string  `json:"roleCTime"`
//		Os 				string	`json:"os"`
//		RoleLevelMTime 	string 	`json:"roleLevelMTime"`
	} `json:"content"`
}

	type LoginGameRole struct {
		ZoneId string `json:"zoneId"`
		ZzoneName string `json:"zoneName"`
		RoleId string `json:"roleId"`
		RoleName string `json:"roleName"`
		RoleCTime string `json:"roleCTime"`
		RoleLevel string `json:"roleLevel"`
//		Os string `json:"os"`
//		RoleLevelMTime string `json:"roleLevelMTime"`
	}
	var err error
	goto BEGIN

ERROR:

	LogE(c, err.Error())
	//c.JSON(http.StatusOK, UCCallBackRequest{Result: "FAILURE"})
	return

BEGIN:
 	var request UCSubmitRoleDataRequest	//请求数据结构
 	var role loginGameRole 	//游戏据类型（category）
	
	request.Id = fmt.Sprint(time.Now().Unix())
	request.Service = "ucid.game.gameData"
	request.Data.AccountId = buf
	request.Game.GameId = gameId

	//gameData
	GameData.Category = "LoginGameRole"
	//游戏服务器的角色数据--存储于数据库
	//TODO   如果玩家数据存在换行符，则要替换成空串
	if GameData.Content != "\n" {
		GameData.Content.ZoneId = 
		GameData.Content.ZoneName =
		GameData.Content.RoleId =
		GameData.Content.RoleName =
		GameData.Content.RoleLevel =
		GameData.Content.RoleCTime =
	} else {
		GameData.Content = ""
	}

	//  可选数据
	//	RoleLevelMTime----当用户的角色等级发生变化后调用GameData.Content.RoleLevelMTime
	//	GameData.Content.Os = "android"
	//	GameData.Content.RoleLevelMTime =


	//GameData经过json encode变成字符串
	request.Data.GameData = base64.URLEncoding.EncodeToString([]byte(GameData))
	fmt.Println(request.Data.GameData)

	//sign
	sign := "accountId=" + request.Data.AccountId + "gameData=" + request.Data.GameData
	fmt.Println(sign)
	request.Sign = MD5(([]byte)(sign))

	// add to parse
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	roleReader := bytes.NewBuffer(jsons)
    roleUrl := "http://collect.sdkyy.9game.cn:8080/ng/cpserver/gamedata/ucid.game.gameData"
    request, err := http.NewRequest("POST", roleUrl, roleReader)
    if err != nil {
        goto ERROR
    }
    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        goto ERROR
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    //fmt.Println(respBytes)
    respStr := string(respBytes)
    fmt.Println(respStr)
    if err != nil {
        goto ERROR
    }
    defer resp.Body.Close()
    str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)

	type UCSubmitRoleDataRespond struct {
		Id string `json:"id"`
		State struct {
			Code string `json:"code"`
			Msg string `json:"msg"`
		}	`json:"state"`
		Data struct {}	`json:"data"`
	}
	var respond UCSubmitRoleDataRespond //响应数据结构
	err = json.Unmarshal(([]byte)(respStr), &respond)
	if err != nil {
		err = errors.New("UCSubmitRoleDataRespond request json parse error")
		goto ERROR
	}
	//   
	//请求与返回的时间戳判断是否一致的判断
	//if respond.Id == request.Id {
	//}
*/
}