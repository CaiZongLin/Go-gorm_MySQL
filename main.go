package main

import (
	"Lex/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/product_list", database.SearchAll) //搜尋所有產品
	router.GET("/product", database.SearchProduct)  //Query Search 搜尋單一商品-> product?name=
	router.GET("/product_name", database.SearchProductName)
	router.POST("/create", database.InsertProduct)   //新增商品
	router.POST("/modify", database.ModifyProduct)   // 修改商品資訊
	router.POST("/buy", database.Buy)                //購買商品
	router.GET("/performance", database.Performance) // 查看營業額
	router.Run(":80")

}

func returnJson(c *gin.Context) {
	m := map[string]string{"status": "ok"}
	j, _ := json.Marshal(m)
	c.Data(http.StatusOK, "application/json", j)
}
