package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	constShenfutongWeixin = 21
	constShenfutongAliPay = 22
	//constPayChannelIDmax  = 3

	NotifyUrl             string = PAY_NOTIFY_HOST_URL
	NotifyPort            string = ":8090"
	constShenfutongNotify string = "/ShenfutongNotify"
)
const (
	/*
		constShenfutongPaykey   string = "44F3411D4C4A490F91D2D16C893C5145"
		constShenfutongAppId    string = "1681"
		constShenfutongMerId    string = "21100600801"
		constPostOrderWeixinUrl string = "http://payment-test.szhuyu.com/wap/weixin/wft"
		constPostOrderAliPayUrl string = "http://payment-test.szhuyu.com/wap/alipay/v2"
	*/
	constPostOrderWeixinUrl string = "http://payment.szhuyu.com/wap/weixin/wft"
	constPostOrderAliPayUrl string = "http://payment.szhuyu.com/wap/alipay/v2"
	constShenfutongPaykey   string = "B61997B5FB664B299C5E00418FB1BECF"
	constShenfutongAppId    string = "17410017"
	constShenfutongMerId    string = "1704101722644663"
)

/* order ID */
/* book order by thirdPayChannelID */
func PayOrder(thirdPayChannelID, fee int, encodeOrderId, itemDesc, callbackUrl string) (thirdPayOrder, thirdPayUrl string, err error) {
	err = nil
	if thirdPayChannelID == constShenfutongWeixin {
		var sftPay = ShenfutongWeiXin{constShenfutongPaykey, constShenfutongAppId, constShenfutongMerId}
		PostOrderUrl := constPostOrderWeixinUrl
		thirdPayOrder, thirdPayUrl, err = sftPay.PostShenfutongWeixinOrder(PostOrderUrl, itemDesc, fee, "用户名称", callbackUrl+constShenfutongNotify, encodeOrderId)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
		fmt.Println(thirdPayOrder, thirdPayUrl)
		return thirdPayOrder, thirdPayUrl, err

	} else if thirdPayChannelID == constShenfutongAliPay {
		var sftPay = ShenfutongAlipay{constShenfutongPaykey, constShenfutongAppId, constShenfutongMerId}
		PostOrderUrl := constPostOrderAliPayUrl
		thirdPayOrder, thirdPayUrl, err = sftPay.PostShenfutongAlipayOrder(PostOrderUrl, itemDesc, fee, "用户名称", callbackUrl+constShenfutongNotify, encodeOrderId)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}

		fmt.Println(thirdPayUrl)
		return "", thirdPayUrl, err
	}
	err = errors.New("payorder")
	return "", "", err
}

/* new order */
func Order(c *gin.Context, payOrderTable PayOrderTable, accessGameIDMap map[string]GameIDRow) {

	type PayRequest struct {
		PayID         string `json:"payID"`
		CpUserID      string `json:"userID"`
		GameAccessKey string `json:"key"`
		ServerID      string `json:"serverID"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		//DeviceInfo    string `json:"deviceInfo"`
		CpDate string `json:"cpDate"`
		Fee    string `json:"fee"`
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
	fmt.Println(err)
	LogE(c, err.Error())
	c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: "申请订单失败。", PayUrl: ""})
	return

BEGIN:
	var request PayRequest
	requestJson := c.PostForm("request")
	fmt.Println("request=" + requestJson)

	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		err = errors.New("json parse error")
		goto ERROR
	}
	//XXX:

	platformID, err := strconv.Atoi(request.PlatformID)
	if err != nil {
		err = errors.New("platformID not correct")
		goto ERROR
	}
	fmt.Println("-----------")

	thirdPayChannelID, err := strconv.Atoi(request.PayID)
	if err != nil || thirdPayChannelID <= 0 {
		err = errors.New("payID not correct")
		goto ERROR
	}
	fmt.Println("-----------", thirdPayChannelID)

	fmt.Println(platformID)
	if thirdPayChannelID < 20 {
		if 1 < platformID && platformID < 20 {
			thirdPayChannelID = platformID
		}
	}
	fmt.Println("-----------", thirdPayChannelID)

	accountID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo, err := ParseCommonPayRequest(accessGameIDMap, requestJson)
	//fmt.Println(accountID, thirdPayChannelID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo, err)
	_ = deviceInfo
	fmt.Println()
	if err != nil {
		goto ERROR
	}
	fmt.Println("-----------")
	//request.DeviceInfo = deviceInfo
	fmt.Println("-----------")

	/*
			gameIDDate := accessGameIDMap[request.GameAccessKey]
			if request.GameAccessKey == "" || gameIDDate.AccessKey != request.GameAccessKey {
		/*
			gameIDDate := accessGameIDMap[request.GameAccessKey]
			if request.GameAccessKey == "" || gameIDDate.AccessKey != request.GameAccessKey {
				err = errors.New("key error")
				goto ERROR
			}
			gameID := gameIDDate.GameID

			request.DeviceInfo = FixDeviceInfo(request.DeviceInfo)
			deviceID, err := GetDeviceID(request.DeviceInfo)

			platformID, err := strconv.Atoi(request.PlatformID)
			if err != nil {
				err = errors.New("platformID not correct")
				goto ERROR
			}
			// channlID

			channelID, err := strconv.Atoi(request.ChannelID)
			if err != nil {
				err = errors.New("channelID not correct")
				goto ERROR
			}

			// serverID

			serverID, err := strconv.Atoi(request.ServerID)
			if err != nil || serverID < 0 {
				err = errors.New("serverID not correct")
				goto ERROR
			}

			userID, decodeGameID, err := CPUserIDDecode(request.CpUserID)

			if err != nil {
				err = errors.New("userID not correct")
				goto ERROR
			}

			// recheck gameid
			if decodeGameID != gameID {
				err = errors.New("userID and gameId not according")
				goto ERROR
			}

			fee, err := strconv.Atoi(request.Fee)
			if err != nil || fee < 10 {
				err = errors.New("fee not correct")
				goto ERROR
			}

			thirdPayChannelID, err := strconv.Atoi(request.PayID)
			if err != nil || thirdPayChannelID <= 0 {
				err = errors.New("payID not correct")
				goto ERROR
			}
	*/

	// orderID, orderTime := payOrderTable.NewOrderID(userID, gameID)
	orderID, orderTime := EncodePayOrderID(gameID, thirdPayChannelID)

	// according to thirdPayChannelID
	fmt.Println("************")
	thirdPayOrderID, thirdPayUrl, err := PayOrder(thirdPayChannelID, fee, orderID, "商品名称", NotifyUrl)
	fmt.Println(gameID, thirdPayChannelID, orderID, orderTime, thirdPayChannelID, thirdPayUrl)
	if err != nil {
		err = errors.New("new pay order fail")
		goto ERROR
	}

	/* save back order */
	//XXX:
	fmt.Println("************")
	err = payOrderTable.NewBookOrder(gameID, orderID, orderTime,
		accountID, serverID,
		deviceID,
		platformID,
		channelID, fee, request.CpDate, thirdPayOrderID, thirdPayChannelID, thirdPayUrl)

	if err != nil {
		err = errors.New("new pay order datebase fail")
		goto ERROR
	}

	//c.String(200, thirdPayUrl)
	//XXX:
	redirectUrl = thirdPayUrl
	fmt.Println(orderID, thirdPayUrl)
	c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: "", PayOrderID: orderID, PayUrl: thirdPayUrl})
	return

}

var redirectUrl = ""

func CheckOrder(c *gin.Context, date string, min time.Duration) {
	var err error
	goto BEGIN
ERROR:
	LogE(c, err.Error())
	return
BEGIN:
	time.Sleep(min)

	requestPkg := url.Values{"request": {date}}
	response, err := http.PostForm("http://127.0.0.1:9966/checkPay", requestPkg)
	if err != nil {
		goto ERROR
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		goto ERROR
	}
	fmt.Println(string(body))
	defer response.Body.Close()

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

	// database table
	var payOrderTable PayOrderTable
	payOrderTable = PayOrderTable{db}

	/* get gameID */
	gameIDTable := GameIDTable{db}
	accessGameIDMap, err := gameIDTable.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	idGameIDMap, err := gameIDTable.GetAllByID()
	if err != nil {
		fmt.Println(err)
		return
	}

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

	// Op
	router.POST("/order", func(c *gin.Context) {
		Order(c, payOrderTable, accessGameIDMap)
	})
	router.GET("/pay", func(c *gin.Context) {
		fmt.Println(redirectUrl)
		c.Redirect(302, redirectUrl)
	})
	router.POST("/selectPay", func(c *gin.Context) {

		type PayRequest struct {
			PayID         string `json:"payID"`
			CpUserID      string `json:"userID"`
			GameAccessKey string `json:"key"`
			ServerID      string `json:"serverID"`
			ChannelID     string `json:"channelID"`
			PlatformID    string `json:"platformID"`
			VersionID     string `json:"versionID"`
			//DeviceInfo    string `json:"deviceInfo"`
			CpDate   string `json:"cpDate"`
			Fee      string `json:"fee"`
			RoleInfo string `json:"roleInfo"`
			/**/
			//PayOrderID string `json:"payOrderId"`
			//AccountID string `json:"accountID"`
		}

		type RoleInfo struct {
			RoleLevel string `json:"gameUserLevel"`
		}

		type PayRespond struct {
			Result     string `json:"result"`
			ThirdPay   bool   `json:"thirdPay"`
			PayOrderID string `json:"payOrderID"`
			Message    string `json:"message"`
		}

		goto BEGIN
	ERROR:
		fmt.Println(err)
		LogE(c, err.Error())
		c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: "申请订单失败"})
		return
	BEGIN:
		//c.JSON(http.StatusOK, PayRespond{Result: "FAIL", Message: "申请订单失败"})

		var request PayRequest
		requestJson := c.PostForm("request")
		fmt.Println(requestJson)

		err = json.Unmarshal([]byte(requestJson), &request)
		if err != nil {
			err = errors.New("request json parse error")
			goto ERROR
		}

		//errInfo := ""
		//requestJson := c.PostForm("request")

		accountID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo, err := ParseCommonPayRequest(accessGameIDMap, requestJson)
		thirdPayChannelID := platformID
		fmt.Println("1111111")
		if err != nil {
			goto ERROR
		}
		_ = accountID
		_ = thirdPayChannelID
		_ = fee
		tableID := gameID
		_ = tableID
		_ = channelID
		_ = platformID
		_ = serverID
		_ = deviceID
		_ = deviceInfo
		fmt.Println("1111111")

		///*
		//XXX:
		///*
		if platformID != 3 {

			var roleInfo RoleInfo
			err = json.Unmarshal([]byte(request.RoleInfo), &roleInfo)
			if err != nil {
				err = errors.New("role json parse error")
				goto ERROR
			}

			roleLevel, err := strconv.Atoi(roleInfo.RoleLevel)
			if err != nil {
				err = errors.New("role level not correct")
				goto ERROR
			}

			fmt.Println("level :" + roleInfo.RoleLevel)
			if roleLevel >= 70 {
				c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: "", ThirdPay: false, PayOrderID: ""})
				return
			}
			fmt.Println("level :" + roleInfo.RoleLevel)

			fmt.Println("1111111")
			count, err := payOrderTable.GetCount(tableID, accountID, platformID)
			if err != nil {
				goto ERROR
			}

			fmt.Println("1111111")
			if count > 2 {
				c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: "", ThirdPay: false, PayOrderID: ""})
				return
			}
		}
		//*/

		fmt.Println("1111111")
		// orderID, orderTime := payOrderTable.NewOrderID(userID, gameID)
		orderID, orderTime := EncodePayOrderID(gameID, thirdPayChannelID)

		// according to thirdPayChannelID
		//thirdPayOrderID, thirdPayUrl, err := PayOrder(thirdPayChannelID, fee, orderID, "商品名称", NotifyUrl)

		err = payOrderTable.NewBookOrder(gameID, orderID, orderTime,
			accountID, serverID,
			deviceID,
			platformID,
			channelID, fee, request.CpDate, "", thirdPayChannelID, "")

		if err != nil {
			goto ERROR
		}

		fmt.Println(orderID)
		c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: "", ThirdPay: true, PayOrderID: orderID})
		if thirdPayChannelID == 2 {
			go CheckOrder(c, c.PostForm("request"), time.Second*3*60)
			go CheckOrder(c, c.PostForm("request"), time.Second*5*60)
			go CheckOrder(c, c.PostForm("request"), time.Second*10*60)
		}
		//return
		/*
			gameIDDate := accessGameIDMap[request.GameAccessKey]
			if request.GameAccessKey == "" || gameIDDate.AccessKey != request.GameAccessKey {
				err = errors.New("key error")
				goto ERROR
			}
			gameID := gameIDDate.GameID
			tableID := gameID
		*/

		//payOrderTable.GetCount(tableID int, accountID int64, platformID int) (count int, err error) {

		/*
			  fmt.Println(requestJson)
			  i = i + 1
			  if i%2 == 0 {
				c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: errInfo, ThirdPay: true})
			  } else {
				c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: errInfo, ThirdPay: false, PayOrderID: ""})
			  }
		*/
		//c.JSON(http.StatusOK, PayRespond{Result: "SUCCESS", Message: errInfo, PayUrl: "http://www.baidu.com"})
		/*
		  c.JSON(http.StatusOK, CpLoginRespond{Result: "SUCCESS", Message: ""})
		*/
		/* for midashi weixin*/
		//if platformID == 2 {
		/*
			request.PayOrderID = orderID
			time.Sleep(3 * time.Second)

			requestPkg := url.Values{"request": {c.PostForm("request")}}
			response, err := http.PostForm("http://127.0.0.1:9966/checkPay", requestPkg)
			if err != nil {
				goto ERROR
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				goto ERROR
			}
			fmt.Println(string(body))
			defer response.Body.Close()
		*/
		//c.String(http.StatusOK, string(body))
		//}

	})

	router.POST(constShenfutongNotify, func(c *gin.Context) {
		fee, orderID, ThirdPayOrderID, orderTime, err := GetShenfutongPayedNotify(constShenfutongAppId, constShenfutongMerId, constShenfutongPaykey, c)
		//func (cp CpPayNotify) PostToCpNotify(userID string, fee int, orderID string, payTime string, cpDate string) error {
		fmt.Println(fee, ThirdPayOrderID, orderID, orderTime, err)
		//cp := CpPayNotify{"44F3411D4C4A490F91D2D16C893C5145", "http://127.0.0.1:5723", 1}
		err = SaveOrderAndNotifyCP(idGameIDMap, payOrderTable, orderID)
	})

	router.POST("/thirdPay", func(c *gin.Context) {
		type ThirdPayNotifyRequest struct {
			OrderID string `json:"cpOrderId"`
			Fee     string `json:"payFee"`
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

		var request ThirdPayNotifyRequest
		err = json.Unmarshal([]byte(requestJson), &request)
		if err != nil {
			err = errors.New("third pay json parse error")
			goto ERROR
		}

		fee, err := strconv.Atoi(request.Fee)
		if err != nil {
			err = errors.New("fee not correct")
			goto ERROR
		}

		payOrderID := request.OrderID
		gameID, payChannelID, err := DecodePayOrderID(payOrderID)
		_ = payChannelID
		if err != nil {
			err = errors.New("third pay get gameID fail")
			goto ERROR
		}
		tableID := gameID
		row, err := payOrderTable.Get(tableID, payOrderID)

		if err != nil {
			goto ERROR
		}

		if fee == row.Fee {
			if err != nil || fee < 10 {
				err = errors.New("third order fee not correct")
				goto ERROR
			}
		}
		payTime := request.PayTime
		_ = payTime
		err = SaveOrderAndNotifyCP(idGameIDMap, payOrderTable, payOrderID)
		if err != nil {
			goto ERROR
		}
	})

	router.POST("/fixUrl", func(c *gin.Context) {
		/*
			  url := c.PostForm("url")
			  gameIDStr := c.PostForm("gameID")
			  sign := c.PostForm("sign")
			  //fmt.Println(url, gameIDStr, sign)
			  gameID, err := strconv.Atoi(gameIDStr)
			  if err != nil {
				c.String(404, "fail")
				fmt.Println(err)
				return
			  }
			  md5sign := md5.Sum([]byte(gameIDStr + url + accessGameIDMap[gameID].Key))
			  //fmt.Println(cpNotifyMap[gameID].Key)
			  md5str := hex.EncodeToString(md5sign[:])
			  //fmt.Println(md5str)
			  if md5str == sign {
				if url != "" {
				  payNotifyTable.FixUrl(gameID, url)
				  cpNotifyMap[gameID] = PayNotifyRow{cpNotifyMap[gameID].GameID, cpNotifyMap[gameID].Key, url}
				  c.String(200, "success")
				  return
				}
			  }
			  c.String(404, "fail")
		*/
	})
	http.ListenAndServe(NotifyPort, router)
}
