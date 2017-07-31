package main

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EncodePayOrderID(gameID int, payChannelID int) (orderID string, orderTime string) {
	now := time.Now()
	orderTime = now.Format("2006-01-02 15:04:05")
	orderID = now.Format("20060102150405") + "_" + fmt.Sprintf("%04d", gameID) + "_" + fmt.Sprintf("%04d", payChannelID) + "_" + string(Krand(6, KC_RAND_KIND_ALL))
	//GetRandomSalt() //h+ fmt.Sprintf("%016d", accountID)

	return orderID, orderTime
}

func DecodePayOrderID(orderID string) (gameID int, payChannelID int, err error) {
	strArray := strings.Split(orderID, "_")
	fmt.Println(strArray)
	//return 0, 0, nil
	///*
	gameID, err = strconv.Atoi(strArray[1])
	if err != nil {
		return -1, -1, err
	}
	payChannelID, err = strconv.Atoi(strArray[2])
	if err != nil {
		return -1, -1, err
	}
	return gameID, payChannelID, nil
	//*/

}

func CPUserIDEncode(accountID int64, gameID int) (encryptedStr string, err error) {

	// md5
	const md5key = "TanCheng-YouXiDanDan"
	accountIDStr := fmt.Sprintf("%010d", accountID)
	gameIDStr := fmt.Sprintf("%04d", gameID)
	md5decode := accountIDStr + gameIDStr + md5key
	//fmt.Println(md5decode)
	md5encode := md5.Sum([]byte(md5decode))
	md5encodeStr := hex.EncodeToString(md5encode[:])
	//fmt.Println(md5encodeStr)
	desdecrypte := accountIDStr + gameIDStr + md5encodeStr[:2]

	// 3des
	const triplekey = "TanCheng" /*8*/ + "YouXiDan" + "Dan" /*11*/ + "Compa" /*ny 7*/
	ciphertext := []byte("TanChengYouXiDanDan")
	iv := ciphertext[:des.BlockSize] // const BlockSize = 8
	block, err := des.NewTripleDESCipher([]byte(triplekey))
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	//fmt.Printf("%s %d\n",desdecrypte,len(desdecrypte))

	encrypted := make([]byte, len(desdecrypte))
	mode.CryptBlocks(encrypted, []byte(desdecrypte))
	//fmt.Printf("%s encrypt to %x \n", desdecrypte, encrypted)
	encryptedStr = hex.EncodeToString(encrypted[:])
	return encryptedStr, nil
}

func CPUserIDDecode(encryptedStr string) (accountID int64, gameID int, err error) {

	// 3des
	const triplekey = "TanCheng" /*8*/ + "YouXiDanDan" /*11*/ + "Compa" /*ny 7*/
	ciphertext := []byte("TanChengYouXiDanDan")
	iv := ciphertext[:des.BlockSize] // const BlockSize = 8
	block, err := des.NewTripleDESCipher([]byte(triplekey))
	if err != nil {
		return -1, -1, err
	}

	encrypted, err := hex.DecodeString(encryptedStr)
	if err != nil {
		return -1, -1, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)
	//fmt.Printf("%x decrypt to %s\n", encrypted, decrypted)
	//fmt.Printf("%s decrypt fo %x\n", decrypted,encrypted )

	decryptedStr := string(decrypted[:])
	accountIDStr := decryptedStr[:10]
	gameIDStr := decryptedStr[10:14]
	//fmt.Println(accountIDStr)
	//fmt.Println(gameIDStr)
	//fmt.Println(decryptedStr[24:])

	// md5:
	const md5key = "TanCheng-YouXiDanDan"
	md5decode := accountIDStr + gameIDStr + md5key
	//fmt.Println(md5decode)
	md5encode := md5.Sum([]byte(md5decode))
	md5encodeStr := hex.EncodeToString(md5encode[:])
	//fmt.Println(md5encodeStr)
	if md5encodeStr[:2] != decryptedStr[14:] {
		err = errors.New("md5 fail")
		return -1, -1, err
	}

	accountID, err = strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		return -1, -1, err
	}
	gameID, err = strconv.Atoi(gameIDStr)
	if err != nil {
		return -1, -1, err
	}

	//accountID = strconv.

	return accountID, gameID, nil
}

/*
func Get( accountID int64, gameID int ) {
  const (
	md5key = "TanCheng-YouXiDanDan"
	triplekey = "TanCheng"+"YouXiDanDan"+"Compa"//ny 7
  )
  ciphertext := []byte("TanChengYouXiDanDan")
  iv := ciphertext[:des.BlockSize] // const BlockSize = 8
  accountIDStr := fmt.Sprintf( "%016d", accountID)
  gameIDStr  := fmt.Sprintf("%08d", gameID )
  //md5decode := accountIDStr + gameIDStr + md5key
  //md5encode := md5.Sum( []byte(md5decode) )
  //md5encodeStr := hex.EncodeToString( md5encode[:] )
  desdecrypte := accountIDStr + gameIDStr //+ md5encodeStr
  block,err := des.NewTripleDESCipher([]byte(triplekey))
  if err != nil {
	fmt.Printf("%s \n", err.Error())
  }
  mode := cipher.NewCBCEncrypter(block, iv)
  fmt.Printf("%s %d\n",desdecrypte,len(desdecrypte))

  encrypted := make([]byte, len(desdecrypte))
  mode.CryptBlocks(encrypted, []byte(desdecrypte))
  fmt.Printf("%s encrypt to %x \n", desdecrypte, encrypted)
  _=err
  _=mode


  decrypter := cipher.NewCBCDecrypter(block, iv)
  decrypted := make([]byte, len(encrypted))
  decrypter.CryptBlocks(decrypted, encrypted)
  //fmt.Printf("%x decrypt to %s\n", encrypted, decrypted)
  fmt.Printf("%s decrypt fo %x\n", decrypted,encrypted )

}
func main() {
  //Get( 1, 1 );
  str,err := Encrypt(1213123990,123)
  if err != nil {
	fmt.Println(err)
  }
  a,b,err:=Decrypt(str)
  if err != nil {
	fmt.Println(err)
  }
  fmt.Println(a,b)

  const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02 15:04 MST"
  )

  //v := "2017-03-12 13:57 UTC"
  v := "2017-03-12 13:23 CST"
  then, err := time.Parse(timeFormat, v)
  if err != nil {
	fmt.Println(err)
	return
  }
  duration := time.Since(then)
  fmt.Println(duration.Hours())

  //g := time.Now()
  ////q, err := time.Parse(timeFormat, g)
  //d:= time.Since(g)
  //fmt.Println(d.Hours())

  fmt.Println(time.Now())

  fmt.Println(time.Now().Sub(then).Hours())
  fmt.Println(time.Now().Sub(then).Minutes())

}
*/

/*
func main() {
  // because we are going to use TripleDES... therefore we Triple it!
  triplekey := "12345678" + "12345678" + "12345678"
  // you can use append as well if you want



  // plaintext will cause panic: crypto/cipher: input not full blocks
  // IF it is not the correct BlockSize. ( des.BlockSize = 8 bytes )
  // to fix this issue, plaintext may need to be padded to the whole block
  // ( 8 bytes ) for the simplicity of this tutorial, we will just keep
  // the plaintext input to 8 bytes
  plaintext := []byte("Hello Wo") // Hello Wo = 8 bytes.



  block,err := des.NewTripleDESCipher([]byte(triplekey))

  if err != nil {
	fmt.Printf("%s \n", err.Error())
	os.Exit(1)
  }

  fmt.Printf("%d bytes NewTripleDESCipher key with block size of %d bytes\n", len(triplekey), block.BlockSize)


  ciphertext := []byte("ABCDEF1234567890")
  iv := ciphertext[:des.BlockSize] // const BlockSize = 8

  // encrypt

  mode := cipher.NewCBCEncrypter(block, iv)
  encrypted := make([]byte, len(plaintext))
  //mode.CryptBlocks(ciphertext[des.BlockSize:], plaintext)
  mode.CryptBlocks(encrypted, plaintext)
  fmt.Printf("%s encrypt to %x \n", plaintext, encrypted)


  //decrypt
  decrypter := cipher.NewCBCDecrypter(block, iv)
  decrypted := make([]byte, len(plaintext))
  decrypter.CryptBlocks(decrypted, encrypted)
  fmt.Printf("%x decrypt to %s\n", encrypted, decrypted)

}
*/
