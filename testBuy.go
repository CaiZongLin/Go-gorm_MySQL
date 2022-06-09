package main

import (
	"Lex/database"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	process := time.Now() //計算時間
	wg := new(sync.WaitGroup)
	//列出所有商品
	allProduct := database.SearchAll1()
	var total_price float64 //計算當下客人購買金額
	var rwlock sync.RWMutex
	customer := [5]string{"Lexus", "Karen", "Jimmy", "Oscar", "Rayne"} //五個客人
	for i := 0; i < 5; i++ {
		go func(name string) {
			wg.Add(1)
			for i := 0; i < 20; i++ { //購買次數
				number := randomInt(len(allProduct))        //產生隨機數字，範圍為所有商品的數量
				_, price := toBuy(allProduct[number], name) //購買
				rwlock.Lock()
				total_price += price.(float64) //購買成功吐回購買品項金額，加入到total_price
				rwlock.Unlock()
				time.Sleep(1 * time.Second) //延遲一秒
			}
			wg.Done()
		}(customer[i])
	}
	time.Sleep(1 * time.Second)
	wg.Wait()
	dbPrice := performance() //到DB撈資料
	if total_price == 0 {    //比對
		fmt.Println("沒有客人購買")
		return
	} else if total_price == dbPrice.(float64) {
		fmt.Println("與資料庫比對金額正確")
	} else {
		fmt.Println("與資料庫比對金額錯誤")
	}

	fmt.Printf("耗費時間:%v", time.Since(process))
}

func toBuy(name, customer string) (interface{}, interface{}) { //模擬客人購買
	fmt.Println(customer)
	response, err := http.Post(
		"http://localhost/buy",
		"application/x-www-form-urlencoded",
		strings.NewReader("name="+name+"&customer="+customer),
	)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var code map[string]interface{}
	json.Unmarshal([]byte(body), &code)

	if code["code"] == "-3" {
		fmt.Println(customer + "購買失敗,原因:庫存不足 ")
		return nil, 0.0
	} else {
		fmt.Println(customer + "購買成功")
	}
	return code["product_name"], code["product_price"]
}
func performance() interface{} { //撈DB的帳務
	url := "http://localhost/performance"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var code map[string]interface{}
	json.Unmarshal([]byte(body), &code)
	return code["total_price"]
}

func randomInt(number int) int { //隨機產生數字
	number2 := int64(number) //轉成int64才能放到下面big.NewInt
	result, _ := rand.Int(rand.Reader, big.NewInt(number2))
	num := result.String()         //big.Int轉String
	num2, err := strconv.Atoi(num) //String轉int
	if err != nil {
		fmt.Println(err)
	}
	return num2
}
