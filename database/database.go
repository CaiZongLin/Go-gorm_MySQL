package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "127.0.0.1"
	database = "vending"
	user     = "search"
	password = "123456"
	root     = "root"
	root_pwd = "s850429s"
)

func SearchProduct(c *gin.Context) {

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}

	name := c.DefaultQuery("name", "")
	if name == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message":  "網址後方沒有輸入商品",
			"correct_link": "http://localhost/product?name={product_name}",
		})
		return
	}
	var price, inventory int
	var product string
	current_product := SearchAll1()
	err = db.QueryRow("select name,price,inventory FROM product_info where name=?", name).Scan(&product, &price, &inventory)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message":     "查詢的商品不存在",
			"current_product": current_product,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product_name":      name,
		"product_price":     price,
		"product_inventory": inventory,
	})

}

func SearchAll(c *gin.Context) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "connect MySQL failed",
			"status":      err,
		})
		return
	}
	rows, err := db.Query("select * from product_info")
	defer rows.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": err,
		})
		return
	}
	for rows.Next() {
		var uid int
		var name string
		var price int
		var inventory int
		var status int
		var update_time []uint8

		err = rows.Scan(&uid, &name, &price, &inventory, &status, &update_time)
		checkErr(err)
		c.IndentedJSON(http.StatusOK, gin.H{
			"product_name":      name,
			"product_price":     price,
			"product_inventory": inventory,
		})
	}

}

func InsertProduct(c *gin.Context) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", root, root_pwd, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "connect MySQL failed",
			"status":      err,
		})
		return
	}
	defer db.Close()
	var status int
	name := c.DefaultPostForm("name", " ")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "沒有輸入名稱參數",
		})
		return
	}
	err = db.QueryRow("select name from product_info where name = ?", name).Scan(&name)
	if err != nil {
		price := c.DefaultPostForm("price", "0")
		if stringToInt(price) < 0 {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "價格不能為負",
			})
			return
		} else if strings.Contains(price, ".") {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "金額不能為小數",
			})
			return
		}
		inventory := c.DefaultPostForm("inventory", "0")
		if stringToInt(inventory) < 0 {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "庫存數量不能為負",
			})
			return
		} else if strings.Contains(inventory, ".") {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "庫存數量不能為小數",
			})
			return
		}

		if inventory == "0" {
			status = 1
		} else {
			status = 0
		}
		_, err1 := db.Exec("insert into product_info(name,price,inventory,status,update_time) values(?,?,?,?,now())", name, price, inventory, status)
		checkErr(err1)
		c.JSON(http.StatusOK, gin.H{
			"message":           "新增成功",
			"product_name":      name,
			"product_price":     price,
			"product_inventory": inventory,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "你輸入的商品已經存在！請重新輸入",
		})
		return
	}
}

func ModifyProduct(c *gin.Context) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", root, root_pwd, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "connect MySQL failed",
			"status":      err,
		})
		return
	}
	defer db.Close()
	var status int
	id := c.DefaultPostForm("id", " ")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "沒有輸入ID",
		})
		return
	}

	current_productID := searchID()
	err = db.QueryRow("select id from product_info where id = ?", id).Scan(&id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message":       "輸入的ID不存在，請重新確認",
			"current_productID": current_productID,
		})
		return
	}
	new_name := c.DefaultPostForm("new_name", " ")
	if new_name == "" {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "沒有輸入新商品名稱",
		})
		return
	}
	new_price := c.DefaultPostForm("new_price", "0")
	if stringToInt(new_price) < 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "價格不能為負",
		})
		return
	} else if strings.Contains(new_price, ".") {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "金額不能為小數",
		})
		return
	}
	new_inventory := c.DefaultPostForm("new_inventory", "0")
	if stringToInt(new_inventory) < 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "庫存數量不能為負",
		})
		return
	} else if strings.Contains(new_inventory, ".") {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "庫存數量不能為小數",
		})
		return
	}

	if new_inventory == "0" {
		status = 1
	} else {
		status = 0
	}
	_, err1 := db.Exec("update product_info set name=?,price=?,inventory=?,status=? where id=?", new_name, new_price, new_inventory, status, id)
	checkErr(err1)
	c.JSON(http.StatusOK, gin.H{
		"message":           "修改成功",
		"product_name":      new_name,
		"product_price":     new_price,
		"product_inventory": new_inventory,
	})

}

func Buy(c *gin.Context) {
	var inventory, price, status int
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", root, root_pwd, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}
	defer db.Close()

	now_sale_product := searchOnSale()
	buy := c.DefaultPostForm("name", " ")
	if buy == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-1",
			"err_message": "沒有輸入購買商品",
		})
		return
	}
	current_product := SearchAll1()
	err = db.QueryRow("select inventory,price,status from product_info where name = ?", buy).Scan(&inventory, &price, &status)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":            "-2",
			"err_message":     "輸入的商品不存在，請點連結重新確認所有商品名稱",
			"current_product": current_product,
		})
		return
	}
	if status == 1 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-3",
			"err_message": "輸入的商品庫存為0，請看下方目前可販售商品",
			"now_on_sale": now_sale_product,
		})
		return
	}
	customer := c.DefaultPostForm("customer", " ")
	if customer == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-1",
			"err_message": "沒有輸入購買者姓名",
		})
		return
	}
	//新增訂單資訊
	_, err2 := db.Exec("insert into sale_info(customer,production,price,update_time) values(?,?,?,now())", customer, buy, price)
	checkErr(err2)

	//修改商品表庫存
	if inventory == 1 {
		status = 1
	} else {
		status = 0
	}
	_, err1 := db.Exec("update product_info set inventory=?,status=? where name=?", inventory-1, status, buy)
	checkErr(err1)
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":          "0",
		"message":       "購買成功",
		"product_name":  buy,
		"product_price": price,
	})

}
func Performance(c *gin.Context) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", root, root_pwd, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "connect MySQL failed",
			"status":      err,
		})
		return
	}
	defer db.Close()
	timestring := time.Now().Format("2006-01-02")
	date := c.DefaultQuery("date", timestring)
	if len(date) != 10 {
		c.JSON(http.StatusOK, gin.H{
			"err_message":    "日期格式錯誤",
			"correct_format": "YYYY-MM-DD",
		})
		return
	}
	today_sale := "select production,count(*),SUM(price) from sale_info where update_time BETWEEN '" + date + " 00:00:59' AND '" + date + " 23:59:59' GROUP BY production;"
	//today_sale := "select production,count(*),SUM(price) from sale_info GROUP BY production;"
	rows, err := db.Query(today_sale)
	checkErr(err)
	var total_price int
	mapInstances := make(map[string]int)
	for rows.Next() {
		var production string
		var sale_count int
		var price int
		err = rows.Scan(&production, &sale_count, &price)
		checkErr(err)
		mapInstances[production] = sale_count
		total_price += price
	}
	if total_price == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "輸入的日期沒有營業額",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"product_info": mapInstances,
		"total_price":  total_price,
	})

}

//下方為api內用的func
func checkErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}
func buy2(customer string) {

	var buy string
	var inventory, price, status int
	now_sale_product := searchOnSale()
	fmt.Println("目前有販售的商品有:", now_sale_product)
	fmt.Println("請輸入要購買的品項: ")
	fmt.Scanln(&buy)

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", root, root_pwd, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("connect MySQL failed", err)
		return
	}
	defer db.Close()
	err = db.QueryRow("select inventory,price,status from product_info where name = ?", buy).Scan(&inventory, &price, &status)
	if err != nil {
		fmt.Println("您要購買的品項不存在，請重新輸入", err)
		return
	}

	//新增訂單資訊
	_, err2 := db.Exec("insert into sale_info(customer,production,price,update_time) values(?,?,?,now())", customer, buy, price)
	checkErr(err2)

	//修改商品表庫存
	if inventory == 1 {
		status = 1
	} else {
		status = 0
	}

	_, err1 := db.Exec("update product_info set inventory=?,status=? where name=?", inventory-1, status, buy)
	checkErr(err1)

	fmt.Println("購買成功")

}
func searchOnSale() []string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("connect MySQL failed", err)
		return nil
	}
	var product []string
	rows, err := db.Query("select * from product_info")
	defer db.Close()
	for rows.Next() {
		var uid int
		var name string
		var price int
		var inventory int
		var status int
		var update_time []uint8

		err = rows.Scan(&uid, &name, &price, &inventory, &status, &update_time)
		checkErr(err)

		if status == 0 { //status =1為下架
			product = append(product, name)
		}
	}
	return product
}

func SearchAll1() []string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("connect MySQL failed", err)
		return nil
	}
	rows, err := db.Query("select * from product_info")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var product []string
	for rows.Next() {
		var uid int
		var name string
		var price int
		var inventory int
		var status int
		var update_time []uint8

		err = rows.Scan(&uid, &name, &price, &inventory, &status, &update_time)
		checkErr(err)
		product = append(product, name)

	}
	return product
}
func searchID() []int {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("connect MySQL failed", err)
		return nil
	}
	rows, err := db.Query("select * from product_info")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var product []int
	for rows.Next() {
		var uid int
		var name string
		var price int
		var inventory int
		var status int
		var update_time []uint8

		err = rows.Scan(&uid, &name, &price, &inventory, &status, &update_time)
		checkErr(err)
		product = append(product, uid)

	}
	return product
}
func stringToInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	return 0
}
