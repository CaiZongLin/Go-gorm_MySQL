package database

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	UserName   string = "root"
	Password   string = "s850429s"
	Addr       string = "127.0.0.1"
	Port       int    = 3306
	searchPort int    = 3307
	Database   string = "vending"
	root       string = "root"
	root_pwd   string = "s850429s"
)

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Inventory   string `json:"inventory"`
	Status      int64  `json:"status"`
	Update_time time.Time
}

var product Product
var products []Product

type Sale struct {
	ID          int64     `json:"id"`
	Customer    string    `json:"customer"`
	Production  string    `json:"production"`
	Price       string    `json:"price"`
	Update_time time.Time `json:"update_time"`
	Comment     string    `json:"comment"`
}

var sale_info Sale

func (Product) TableName() string {
	return "product_info"
}

func (Sale) TableName() string {
	return "sale_info"
}

func stringToInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	return 0
}

func SearchProduct(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", root, root_pwd, Addr, searchPort, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}

	id := c.DefaultQuery("id", "")
	if id == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message":  "網址後方沒有輸入商品ID",
			"correct_link": "http://localhost/product?id={product_id}",
		})
		return
	}
	if err2 := conn.Where("id=?", id).Find(&product).Error; err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": err2,
		})
		return
	}
	if product.Name == "" { // product.Name is empty means no data
		c.JSON(http.StatusOK, gin.H{
			"err_message": "查詢商品不存在",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"product_name":      product.Name,
			"product_price":     product.Price,
			"product_inventory": product.Inventory,
		})
	}
}

func SearchProductName(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", root, root_pwd, Addr, searchPort, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
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
	if err2 := conn.Where("name=?", name).Find(&product).Error; err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": err2,
		})
		return
	}
	if product.Name == "" { // product.Name is empty means no data
		c.JSON(http.StatusOK, gin.H{
			"err_message": "查詢商品不存在",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"product_id":        product.ID,
			"product_name":      product.Name,
			"product_price":     product.Price,
			"product_inventory": product.Inventory,
		})
	}
}

func SearchAll(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", root, root_pwd, Addr, searchPort, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}

	conn.Find(&products)
	for _, p := range products {
		c.IndentedJSON(http.StatusOK, gin.H{
			"product_name":      p.Name,
			"product_price":     p.Price,
			"product_inventory": p.Inventory,
		})
	}
}

func InsertProduct(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}
	var status int64

	name := c.DefaultPostForm("name", " ")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "沒有輸入名稱參數",
		})
		return
	}

	if errSearch := conn.Where("name=?", name).Find(&product).Error; errSearch != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": errSearch,
		})
		return
	}

	if product.ID == 0 { //0代表商品不存在 ，可以新增
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
		} else if price == "" {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "金額不能為空",
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
		} else if inventory == "" {
			c.IndentedJSON(http.StatusOK, gin.H{
				"err_message": "庫存不能為空",
			})
			return
		}

		if inventory == "0" {
			status = 1
		} else {
			status = 0
		}
		new_product := &Product{Name: name, Price: price, Inventory: inventory, Status: status, Update_time: time.Now()}
		if errCreate := conn.Debug().Create(&new_product).Error; errCreate == nil {
			c.JSON(http.StatusOK, gin.H{
				"message":           "新增成功",
				"product_name":      name,
				"product_price":     price,
				"product_inventory": inventory,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"err_message": errCreate,
			})
			return
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "你輸入的商品已經存在！請重新輸入",
		})
		return
	}
}

func ModifyProduct(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}
	var status int64
	id := c.DefaultPostForm("id", " ")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "沒有輸入ID",
		})
		return
	}

	if errSearch := conn.Where("id=?", id).Find(&product).Error; errSearch != nil {
		c.JSON(http.StatusOK, gin.H{
			"err_message": errSearch,
		})
		return
	}

	if product.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"err_message": "輸入的ID不存在，請重新確認",
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
	} else if new_price == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "庫存不能為空",
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
	} else if new_inventory == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": "庫存不能為空",
		})
		return
	}

	if new_inventory == "0" {
		status = 1
	} else {
		status = 0
	}

	Modify := Product{Name: new_name, Price: new_price, Inventory: new_inventory, Status: status, Update_time: time.Now()}
	if errModify := conn.Debug().Model(&Product{}).Where("id = ?", id).Updates(Modify).Error; errModify == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":           "修改成功",
			"product_name":      new_name,
			"product_price":     new_price,
			"product_inventory": new_inventory,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"err_message": errModify,
		})
		return
	}
}

func Buy(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}

	buy := c.DefaultPostForm("name", " ")
	if buy == "" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-1",
			"err_message": "沒有輸入購買商品",
		})
		return
	}
	conn.Where("name=?", buy).Find(&product)
	if product.ID == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-2",
			"err_message": "輸入的商品不存在，請點連結重新確認所有商品名稱",
		})
		return
	} else if product.Inventory == "0" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":        "-3",
			"err_message": "商品庫存為0,無法購買",
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

	sale_info := &Sale{Customer: customer, Production: buy, Price: product.Price, Update_time: time.Now()}
	conn.Create(&sale_info)

	//修改商品表庫存
	var status int64
	if stringToInt(product.Inventory)-1 == 0 {
		status = 1
	} else {
		status = 0
	}
	if errbuy := conn.Debug().Model(&Product{}).Where("name = ?", buy).Updates(Product{Inventory: fmt.Sprint(stringToInt(product.Inventory) - 1), Status: status}).Error; errbuy == nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":          "0",
			"message":       "購買成功",
			"product_name":  buy,
			"product_price": product.Price,
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": errbuy,
		})
	}
}

func Performance(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "connect MySQL failed",
			"status":  err,
		})
		return
	}

	timestring := time.Now().Format("2006-01-02 15:04:05")
	startdate := c.DefaultQuery("startdate", timestring)
	enddate := c.DefaultQuery("enddate", timestring)
	// if len(date) != 10 {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"err_message":    "日期格式錯誤",
	// 		"correct_format": "YYYY-MM-DD",
	// 	})
	// 	return
	// }

	rows, err := conn.Debug().Table("sale_info").Select("production,count(*) as number ,SUM(price) as money").Where("update_time between ?  and ? ", startdate, enddate).Group("production").Rows()
	//today_sale := "select production,count(*),SUM(price) from sale_info where update_time BETWEEN '" + date + " 00:00:59' AND '" + date + " 23:59:59' GROUP BY production;"
	//today_sale := "select production,count(*),SUM(price) from sale_info GROUP BY production;"
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"err_message": err,
		})
		return
	}
	defer rows.Close()
	mapInstances := make(map[string]int64)
	var total_price int64
	for rows.Next() {
		var production string
		var price int64
		var money int64
		err = rows.Scan(&production, &price, &money)
		if err != nil {
			return
		}
		mapInstances[production] = price
		total_price += money
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
