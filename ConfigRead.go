package main
import (
"fmt"
"encoding/json"
"errors"
"os"
"strings"
"bufio"
)

func ReadKey(projectId,distributorId ,path string)(string){
	relaPath :=  projectId + "_" + distributorId
	 f, err := os.Open(path)
	defer f.Close()
	if nil == err {
		buf := bufio.NewReader(f)
		for {
		        str, err := buf.ReadString('\n')
			if err != nil  {
			            if !(err.Error() == "EOF") {
			                fmt.Println(err)
			                os.Exit(1)
			            }
			}
			str = strings.Replace(str, "\n", "", -1)
		 	fmt.Println("str:",str) 
			//var respond Respond
			//判断relapath和文件的一致		
	        		//除去外围变化的key=projectId+distirbutorId
			m := make(map[string]interface{},1024)
	        		err = json.Unmarshal([]byte(str), &m)
	        		if err != nil {
	        			err = errors.New("str Unmarshal error")
	        			return ""
	        		}
	        		fmt.Println("m:",m)
	        		var readJson []byte
	        		for k,v := range m {  
	        			fmt.Println("relaPath:",relaPath)      			
				if k == relaPath {
					fmt.Println("k",k)
					fmt.Println("v:",v)
					readJson, err = json.Marshal(v)
					if err != nil {
						fmt.Println(err)    
				            		return ""
					}
					fmt.Println("readJson:",string(readJson))
					return string(readJson)
				}
				if err != nil && err.Error() == "EOF"{
				            fmt.Println(err)    
				            return string(readJson)
				}
	        		}
		}
	        	return ""			
	}
	return  ""
}

func Pdp(projectId,distributorId ,path string) string{	
	key := ReadKey(projectId,distributorId ,path)	
	return  key
}