package main

import (
	_ "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/itsjamie/gin-cors"
	_ "net/http"
	"strconv"
	_ "strings"
	_ "time"
)

func ParseCommonRequest(accessGameIDMap map[string]GameIDRow, jsonStr string) (gameID, channelID, platformID int, deviceID, deviceInfo string, err error) {
	type Request struct {
		GameAccessKey string `json:"key"`
		ChannelID     string `json:"channelID"`
		PlatformID    string `json:"platformID"`
		VersionID     string `json:"versionID"`
		DeviceInfo    string `json:"deviceInfo"`
	}
	var request Request
	goto BEGIN

ERROR:
	return -1, -1, -1, "", "", err

BEGIN:

	err = json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		err = errors.New("json parse error")
		goto ERROR
	}

	// gameID
	gameIDDate := accessGameIDMap[request.GameAccessKey]
	if request.GameAccessKey == "" || gameIDDate.AccessKey != request.GameAccessKey {
		err = errors.New("key error")
		goto ERROR
	}
	gameID = gameIDDate.GameID

	// channlID
	channelID, err = strconv.Atoi(request.ChannelID)
	if err != nil {
		err = errors.New("channelID not correct")
		goto ERROR
	}

	// platformID
	platformID, err = strconv.Atoi(request.PlatformID)
	if err != nil {
		err = errors.New("platformID not correct")
		goto ERROR
	}

	//deviceInfo
	deviceInfo = FixDeviceInfo(request.DeviceInfo)

	//deviceID
	deviceID, err = GetDeviceID(request.DeviceInfo)
	if err != nil {
		goto ERROR
	}

	return gameID, channelID, platformID, deviceID, deviceInfo, nil

}

/*
func ParseThirdAccountRequest(accessGameIDMap map[string]GameIDRow, jsonStr string) (gameID, channelID, platformID int, deviceID, deviceInfo, , err error) {
	type ThirdAccountRequest struct {
		PlatformAccountID string `json:"platformAccountID"`
	}

	goto BEGIN

ERROR:
	fmt.Println(err)
	return -1, -1, -1, -1, -1, -1, "", "", err

BEGIN:
	gameID, channelID, platformID, deviceID, deviceInfo, err := ParseCommonRequest(accessGameIDMap, jsonStr)
	if err != nil {
		goto ERROR
	}

}

func ParseCommonAccountRequest(accessGameIDMap map[string]GameIDRow, jsonStr string) (gameID, channelID, platformID int, deviceID, deviceInfo, account, password string, err error) {
	type AccountRequest struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	goto BEGIN

ERROR:
	fmt.Println(err)
	return -1, -1, -1, -1, -1, -1, "", "", "", err

	gameID, channelID, platformID, deviceID, deviceInfo, err := ParseCommonRequest(accessGameIDMap, jsonStr)
	if err != nil {
		goto ERROR
	}

}
*/

func ParseCommonPayRequest(accessGameIDMap map[string]GameIDRow, jsonStr string) (accountID int64,
	fee, gameID, channelID, platformID, serverID int,
	deviceID, deviceInfo string,
	err error) {
	type Request struct {
		//PayID         string `json:"payID"`
		CpUserID string `json:"userID"`
		ServerID string `json:"serverID"`
		//CpDate   string `json:"cpDate"`
		Fee string `json:"fee"`
	}

	goto BEGIN

ERROR:
	fmt.Println(err)
	return -1, -1, -1, -1, -1, -1, "", "", err

BEGIN:
	var request Request

	err = json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		err = errors.New("json parse error")
		goto ERROR
	}

	// accountID
	accountID, decodeGameID, err := CPUserIDDecode(request.CpUserID)

	if err != nil {
		err = errors.New("userID not correct")
		goto ERROR
	}

	// fee
	fee, err = strconv.Atoi(request.Fee)
	if err != nil || fee < 10 {
		err = errors.New("fee not correct")
		goto ERROR
	}

	// gameID, channelID, platformID, deviceID, deviceInfo
	gameID, channelID, platformID, deviceID, deviceInfo, err = ParseCommonRequest(accessGameIDMap, jsonStr)
	if err != nil {
		return
	}

	// serverID
	serverID, err = strconv.Atoi(request.ServerID)
	if err != nil || serverID < 0 {
		err = errors.New("serverID not correct")
		goto ERROR
	}

	// check gameid
	if decodeGameID != gameID {
		err = errors.New("userID and gameId not  consistence")
		goto ERROR
	}

	return accountID, fee, gameID, channelID, platformID, serverID, deviceID, deviceInfo, nil

}
