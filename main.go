package main

import (
	"fmt"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/db_ymt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//	解决查表时自动变成复数的问题，list->lists
			SingularTable: true,
		},
	})
	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(10 * time.Second) //10秒钟

	//结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
	}
	//注意点
	//1、结构体变量大驼峰
	//gorm:指定类型
	//json:json接受时的名称
	//binding:"required" 表示必须传入

	//User = List
	//db.AutoMigrate(&User{})
	db.AutoMigrate(&List{})
	//1、主键问题 首行添加gorm.Model
	//2、lists->list

	//	接口
	r := gin.Default()
	//测试
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	//TODO:增删改查
	//增加
	r.POST("/user/add", func(c *gin.Context) {
		var data List

		err := c.ShouldBindJSON(&data)

		//	判断是否有错误
		if err != nil {
			c.JSON(200, gin.H{
				"message": "add failed",
				"data":    gin.H{},
				"code":    400,
			})
		} else {

			db.Create(&data) //创建数据
			c.JSON(200, gin.H{
				"message": "add successful",
				"data":    data,
				"code":    200,
			})
		}
	})
	//删除
	//1、找到对应的id的条目
	//2、判断id是否存在
	//3、从数据库删除
	//4、返回id已经删除

	r.DELETE("user/delete/:id", func(c *gin.Context) {
		var data []List
		//接受id
		id := c.Param("id")

		//判断id是否存在
		db.Where("id = ?", id).Find(&data)

		//id存在则删除，否则报错
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"message": "id not found",
				"code":    400,
			})
		} else {
			//数据库删除
			db.Where("id = ?", id).Delete(&data)

			c.JSON(200, gin.H{
				"message": "deleted",
				"id":      id,
				"code":    200,
			})
		}
	})
	//修改
	//1、找到对应的id的条目
	//2、判断id是否存在
	//3、从数据库修改
	//4、返回id已经修改
	r.PUT("user/update/:id", func(c *gin.Context) {
		var data List
		//接受id
		id := c.Param("id")

		//判断id是否存在(第二种方法，没找到返回0)
		db.Select("id").Where("id = ?", id).Find(&data)

		//id存在则删除，否则报错
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"message": "id not found",
				"code":    400,
			})
		} else {
			//数据库修改
			err := c.ShouldBindJSON(&data)

			//	判断是否有错误
			if err != nil {
				c.JSON(200, gin.H{
					"message": "change failed",
					"data":    gin.H{},
					"code":    400,
				})
			} else {

				db.Where("id = ?", id).Updates(&data) //修改数据
				c.JSON(200, gin.H{
					"message": "change successful",
					"data":    data,
					"code":    200,
				})
			}
		}
	})
	//查询
	//1、条件查询
	//2、全部查询

	//条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {
		var data []List

		name := c.Param("name")

		db.Where("name = ?", name).Find(&data)

		//判断是否能查询到数据
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"message": "data not found",
				"code":    400,
				"data":    gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"message": "data found",
				"code":    200,
				"data":    data,
			})
		}

	})

	//全部查询
	r.GET("/user/list/", func(c *gin.Context) {
		var data []List

		//获取分成的页数，每页显示上限
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		//返回总数
		var total int64
		//查询数据库
		db.Model(data).Count(&total).Limit(pageSize).Offset(pageNum).Find(&data)

		if len(data) == 0 {
			c.JSON(200, gin.H{
				"message": "data not found",
				"code":    400,
				"data":    gin.H{},
			})

		} else {
			c.JSON(200, gin.H{
				"message": "data found",
				"code":    200,
				"data": gin.H{
					"list":     data,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	//端口号
	PORT := "3001"
	r.Run(":" + PORT)
}
