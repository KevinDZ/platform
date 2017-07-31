package main

import (
	"encoding/json"
)

func FixDeviceInfo(jsonDeviceInfo string) string {
	type DevcieInfo struct {
		IMEI string `json:"DeviceId(IMEI)"`
		//`josn:"DeviceSoftwareVersion"`
		MSISDN string `josn:"Line1Number"`
		//`josn:"NetworkCountryIso"`
		//`josn:"NetworkOperator"`
		//`josn:"NetworkOperatorName"`
		//`josn:"NetworkType"`
		//string `josn:"PhoneType"`
		//string `josn:"SimCountryIso"`
		//string `josn:"SimOperator"`
		//string `josn:"SimOperatorName"`
		ICCID string `josn:"SimSerialNumber"`
		//string `josn:"SimState"`
		IMSI string `josn:"SubscriberId(IMSI)"`
		//string `josn:"VoiceMailNumber"`
		Product string `josn:"Product"`
		CPU_ABI string `josn:"CPU_ABI"`
		//TAGS          string `josn:"TAGS"`
		//VERSION_CODES string `josn:"VERSION_CODES.BASE"`
		MODEL        string `josn:"MODEL"`
		SDK          string `josn:"SDK"`
		VERSION      string `josn:"VERSION.RELEASE"`
		Device       string `josn:"DEVICE"`
		Display      string `josn:"DISPLAY"`
		BRAND        string `josn:"BRAND"`
		BOARD        string `josn:"BOARD"`
		FINGERPRINT  string `josn:"FINGERPRINT"`
		ID           string `josn:"ID"`
		MANUFACTURER string `josn:"MANUFACTURER"`
		//string `josn:"USER"`
		OS string `josn:"OS"`
	}
	var deviceInfo DevcieInfo
	err := json.Unmarshal([]byte(jsonDeviceInfo), &deviceInfo)
	_ = err
	var deviceID = []byte(deviceInfo.IMEI)
	if len(deviceID) >= 14 {
		deviceID[8] = deviceInfo.IMEI[13]
		deviceID[9] = deviceInfo.IMEI[11]
		deviceID[10] = deviceInfo.IMEI[10]
		deviceID[11] = deviceInfo.IMEI[9]
		deviceID[12] = deviceInfo.IMEI[12]
		deviceID[13] = deviceInfo.IMEI[8]
		deviceInfo.IMEI = string(deviceID)
	}
	b, _ := json.Marshal(deviceInfo)
	return string(b)

}
func GetDeviceID(jsonDeviceInfo string) (string, error) {
	type DevcieID struct {
		IMEI string `json:"DeviceId(IMEI)"`
	}
	var deviceID DevcieID
	err := json.Unmarshal([]byte(jsonDeviceInfo), &deviceID)
	if err != nil {
		return "", err
	} else {
		return deviceID.IMEI, nil
	}
}
