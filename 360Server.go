package main

const (
	//360测试demo服务器
	//qihoourl = "http://sdbxapp.msdk.mobilem.360.cn" 
)
var (
	urls string
	buf string
)

func func main() {
	router := gin.Default()
	router.POST("/qihooLogin", func (c *gin.Context) {
		QihooLogin(c)
	})	
}

func QihooLogin(c *gin.Context) {
	type ThirdLoginRequest struct {
		type ThirdLoginRequest struct{	 
	 	ThirdLogin	struct{
	 		AccessToken string `json:"accessToken"`
	 		Fields string `json:"fields"`
		} 	`json:"thirdLogin"`	
	}
	}
	var err error
	goto BEGIN
ERROR:
	fmt.Println(err)
	LogE(c, err.Error())
	return

BEGIN:
	var thirdrequest ThirdLoginRequest
	requestJson := c.PostForm("request")
	fmt.Println("360Login requestJson :",requestJson)
	
	err = json.Unmarshal(requestJson, &thirdrequest)
	if err != nil {
		err = errors.New("360LoginRequest json parse error")
		goto ERROR
	}
	fmt.Println("thirdrequest :",thirdrequest)
	//发送到360服务器
	type QihooLonginRequest struct {
		AccessToken string `json:"accessToken"`
	 	Fields string `json:"fields"`
	}
	var request QihooLonginRequest
	request.AccessToken = thirdrequest.ThirdLoginRequest.AccessToken
	fmt.Println("AccessToken :",request.AccessToken)
	jsons, err := json.Marshal(request)
	if err != nil {
		goto ERROR
	}
	fmt.Println("360 Login request jsons:",string(jsons))

	if thirdrequest.Fields == "" {
		urls = "http://sdbxapp.msdk.mobilem.360.cn?access_token="  + request.AccessToken
		/*reg := regexp.MustCompile(`,`)
		result := reg.FindAllStringSubmatch(thirdrequest.Fields,-1)
		fmt.Printf("result : %q",result)*/
		/*for k , v := range result {			
			fmt.Println("result:",result[k])
		}*/
	}else {
		urls = "http://sdbxapp.msdk.mobilem.360.cn?access_token="  + request.AccessToken + "&fields=" + request.Fields
	}
	/*reader := bytes.NewBuffer(jsons)
	qihoourl = "http://sdbxapp.msdk.mobilem.360.cn?access_token="  + request.AccessToken + "&fields="
	qihoorequest, err := http.NewRequest("GET", qihoourl, reader)
    if err != nil {
        goto ERROR
    }
    fmt.Println("Qihoo request:",qihoorequest)
    client := &http.Client{}
    resp, err := client.Do(qihoorequest)
    if err != nil {
        goto ERROR
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        goto ERROR
    }
    respStr := string(respBytes)
   fmt.Println("qihoo respStr:",respStr)
    defer resp.Body.Close()
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("qihoo  str:",*str)*/
	resp, err := http.Get(urls)
	if err != nil {
		err = errors.New("Get response json parse error")
		goto ERROR
	}
	fmt.Println("HTTP GET:",resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)		//[]byte
	fmt.Println(string(body))
	type QihooLoginRespond struct {
		Id string `json:"id"`
		Name string `json:"name"`
		Avatar string `json:"avatar"`
		Sex string `json:"sex"`	//360用户性别，仅在fields中包含时候才返回,返回值为：男，女或者未知
		Area string `json:"area"`	//360用户地区，仅在fields中包含时候才返回
		Nick string `json:"nick"`	//用户昵称，无值时候返回空
	}
	var respond QihooLoginRespond
	//qihoo返回GET的respond结果
	err := json.Unmarshal(body,&respond)
	if err != nil {
		err = errors.New("QihooLoginRespond json parse error")
		goto ERROR
	}
	fmt.Println("QihooLoginRespond:",respond)

	buf = respond.Id +"@360"
}