package main

import (
	_ "crypto/md5"
	"crypto/md5"
	"database/sql"
	_"encoding/hex"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	_ "time"

	_ "bytes"
	_ "encoding/json"
	"flag"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	_ "log"
	_ "math/rand"
	"net/http"
	_ "net/http/httputil"
	_ "net/url"
	_ "strconv"
	"strings"
	"time"
	"io/ioutil"
)

type ErrorRespond struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type Respond struct {
			Result      string `json:"result"`
			Message     string `json:"message"`
			ChannelID   string `json:"channelID"`
			Hadsend     bool   `json:"hadsend"`
			ReceiveTime string `json:"receiveTime"`
}
func main() {
	flag.Parse()

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	ocpcTable := OcpcTable{db}

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

	router.GET("/ocpc/:name", func(c *gin.Context) {
		var err error
		goto BEGIN
	ERROR:
		LogE(c, err.Error())
		fmt.Println(err)
		c.JSON(http.StatusOK, ErrorRespond{Result: "FAIL", Message: err.Error()})
		return

	BEGIN:
		name := c.Param("name")
		fmt.Println(name)
		strArr := strings.Split(name, "_")
		gameIDStr, platformIDStr, channelIDStr := strArr[0], strArr[1], strArr[2]

		// gameID
		gameID, err := strconv.Atoi(gameIDStr)
		if err != nil {
			err = errors.New("platformID not correct")
			goto ERROR
		}

		// platformID
		platformID, err := strconv.Atoi(platformIDStr)
		if err != nil {
			err = errors.New("platformID not correct")
			goto ERROR
		}

		// channlID
		channelID, err := strconv.Atoi(channelIDStr)
		if err != nil {
			err = errors.New("channelID not correct")
			goto ERROR
		}

		fmt.Println(gameID, platformID, channelID)

		///*
		parm := c.Request.URL.Query()
		LogI(c, fmt.Sprint(parm))
		fmt.Println("c.Request.URL.Query():",parm)
		fmt.Println("callback_url:",parm["callback_url"])
		fmt.Println("adid:",parm["adid"])
		fmt.Println("cid:",parm["cid"])
		fmt.Println("imei:",parm["imei"])
		fmt.Println("ctype:",parm["ctype"])
		fmt.Println("idfa:",parm["idfa"])
		fmt.Println("timestamp:",parm["timestamp"])
		fmt.Println("os:",parm["os"])
		fmt.Println("androidid:",parm["androidid"])
		//android or ios
		if parm["androidid"] != nil {	//如果android存在，则是android手机，存imei
			imeimd5Str := c.Query("imei")
			fmt.Println("imeimd5Str:",imeimd5Str)
			ocpcTable.Insert(gameID, platformID, channelID,
			imeimd5Str,
			"",
			c.Query("callback_url"),
			c.Query("adid"),
			c.Query("cid"),
			c.Query("ctype"))
		c.String(200, "success")
		}else{
			idfa := c.Query("idfa")
			idfamd5 := md5.Sum([]byte(idfa))
			idfamd5Str := hex.EncodeToString(idfamd5[:])
			fmt.Println("idfamd5,idfamd5Str =",idfamd5,idfamd5Str)
			ocpcTable.Insert(gameID, platformID, channelID,
				"",
				idfamd5Str,
				c.Query("callback_url"),
				c.Query("adid"),
				c.Query("cid"),
				c.Query("ctype"))
			c.String(200, "success")
		}		
	})

	router.POST("/ocpcCallback", func(c *gin.Context) {
		var err error
		goto BEGIN
	ERROR:
		LogE(c, err.Error())
		fmt.Println(err)
		c.JSON(http.StatusOK, ErrorRespond{Result: "FAIL", Message: err.Error()})
		return
	BEGIN:
		fmt.Println("c.Request.Form:", c.Request.Form) 
		// platform_id 
		// channel_id 
		// game_id 	
		platformIDStr, isSet := c.GetPostForm("platformID")
			if isSet == false {
				err = errors.New("platformID not correct")
				goto ERROR
			}
			platformID, err := strconv.Atoi(platformIDStr)
			if err != nil {
				err = errors.New("platformID not correct")
				goto ERROR
			}

			gameIDStr, isSet := c.GetPostForm("gameID")
			if isSet == false {
				err = errors.New("gameID not correct")
				goto ERROR
			}
			gameID, err := strconv.Atoi(gameIDStr)
			if err != nil {
				err = errors.New("gameID not correct")
				goto ERROR
			}

			channelIDStr, isSet := c.GetPostForm("channelID")
			if isSet == false {
				err = errors.New("channelID not correct")
				goto ERROR
			}

			channelID, err := strconv.Atoi(channelIDStr)
			if err != nil {
				err = errors.New("channelID not correct")
				goto ERROR
			}

		var imei, idfa , imeimd5Str, idfamd5Str string

		osStr, ok :=  c.GetPostForm("os") 
		if ok == false {
			err = errors.New("ok is false")
			goto ERROR
		}	
		if   osStr != "ios" {
			var isSet bool	
			imei, isSet = c.GetPostForm("imei") //LogI(c, "imei=", imei)
			if isSet == false {
				err = errors.New("isSet is false")
				goto ERROR
			}

			fmt.Println("c.GetPostForm(imei):",imei, gameID, platformID, channelID)
			imeimd5 := md5.Sum([]byte(imei))
			imeimd5Str = hex.EncodeToString(imeimd5[:])
			fmt.Println("imeimd5,imeimd5Str =",imeimd5,imeimd5Str)			
		}else {
			var isSet bool
			idfa, isSet = c.GetPostForm("idfa") //LogI(c, "imei=", imei)
			if isSet == false {
				err = errors.New("isSet is false")
				goto ERROR
			}
			fmt.Println("c.GetPostForm(idfa)",idfa, gameID, platformID, channelID)

			idfamd5 := md5.Sum([]byte(idfa))
			idfamd5Str = hex.EncodeToString(idfamd5[:])
			fmt.Println("idfamd5,idfamd5Str =",idfamd5,idfamd5Str)
		}
		hadsend := false
		callbackUrlArray, channelid , err := ocpcTable.GetCallbackArrayChannelID(imeimd5Str, idfamd5Str, gameID, platformID, channelID)
		if err != nil {
			goto ERROR
			fmt.Println(err)
			return
		}

		isOk := ocpcTable.IsHadsend(gameID, platformID, channelid, imeimd5Str, idfamd5Str)
		if isOk == true {
			hadsend = true
		}

		for _, callbackUrl := range callbackUrlArray {
				fmt.Println("callbackUrl before:",callbackUrl)
				err = ocpcTable.Update(imeimd5Str,idfamd5Str, callbackUrl, gameID)
				if err != nil {
					fmt.Println(err)
					LogE(c, imeimd5Str+" "+idfamd5Str+" "+gameIDStr+" "+channelIDStr+" "+err.Error())
					continue
				}
				fmt.Println("callbackUrl after:",callbackUrl)
				if hadsend == false {
					fmt.Println("send ocpc")
					response, err := http.Get(callbackUrl)
					if err != nil {
						fmt.Println(err)
						//return
						LogE(c, imeimd5Str+" "+idfamd5Str+" "+gameIDStr+" "+channelIDStr+" "+err.Error())
						continue
					}
					fmt.Println("hadsend:",response)
					fmt.Println(" hadsend resp.Body:",response.Body)
					respBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						err = errors.New("respBytes")
						goto ERROR
					}
					respStr := string(respBytes)
					fmt.Println(" respStr:",respStr)
					defer response.Body.Close()
					hadsend = true
					c.JSON(http.StatusOK,Respond{Result:"SUCCESS",Message:""})
					return
				}
		}
		err = errors.New("callbackUrlArray is null")
		goto ERROR		
		//_ = response
	})

	router.POST("/ocpcGetIosChannel",func (c *gin.Context) {
		
		var err error
		goto BEGIN
	ERROR:
		LogE(c, err.Error())
		fmt.Println(err)
		c.JSON(http.StatusOK, ErrorRespond{Result: "FAIL", Message: err.Error()})
		return
	BEGIN:
		
		fmt.Println("c.Request.Form:", c.Request.Form) 
		// platform_id 
		// channel_id 
		// game_id 	
		platformIDStr, isSet := c.GetPostForm("platformID")
		if isSet == false {
			err = errors.New("platformID not correct")
			goto ERROR
		}
		platformID, err := strconv.Atoi(platformIDStr)
		if err != nil {
			err = errors.New("platformID not correct")
			goto ERROR
		}
		fmt.Println("platformID:",platformID)
		gameIDStr, isSet := c.GetPostForm("gameID")
		if isSet == false {
			err = errors.New("gameID not correct")
			goto ERROR
		}
		gameID, err := strconv.Atoi(gameIDStr)
		if err != nil {
			err = errors.New("gameID not correct")
			goto ERROR
		}
		fmt.Println("gameID:",gameID)
		idfaStr, isSet := c.GetPostForm("idfa")
		if isSet == false {
			err = errors.New("idfa not correct")
			goto ERROR
		}

		idfamd5 := md5.Sum([]byte(idfaStr))
		idfamd5Str := hex.EncodeToString(idfamd5[:])
		fmt.Println("idfamd5,idfamd5Str =",idfamd5,idfamd5Str)

		channelId , lasttime , hadsend := ocpcTable.OcpcIOSLastTime(gameID, platformID , idfamd5Str)
		if lasttime == "-1" {
			fmt.Println("时间出错")
			err = errors.New("time")
			goto ERROR
		}
		if channelId == "-1" {
			fmt.Println("渠道出错")
			err = errors.New("channel")
			goto ERROR
		}
		
		fmt.Println("receiveTime:",lasttime)
		fmt.Println("channelId:",channelId)
		
		c.JSON(http.StatusOK,Respond{Result:"SUCCESS",Message:"",ChannelID: channelId,ReceiveTime:lasttime,Hadsend:hadsend})
		
		return
	})
	http.ListenAndServe(":11111", router)
}
