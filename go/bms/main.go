package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	//程序启动连接数据库
	err := initDB()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	//解析模板
	r.LoadHTMLGlob("template/**/*") //模板解析
	//各个路由
	r.GET("/user/login", Login) //登录界面
	r.POST("user/login", UserLogin)
	r.GET("/book/list", book1ListHandle) //查询书籍
	r.GET("/book/new", newBookHandle)    //增加书籍 , 第一次get返回html模板,给用户填写
	r.POST("book/new", createBookHandle)
	r.GET("book/delete", deleteHandle) //删除书籍
	r.GET("/book/update", newHandle)   //更新书籍信息,价格 或书名
	r.POST("book/update", updateHandle)
	r.Run()
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login.html", nil)
}

// 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user, err := GetUserInfo(username)
	if password == user.Password {
		c.Redirect(http.StatusMovedPermanently, "/book/list")
	}
	if err != nil {
		c.String(http.StatusOK, "您还未注册账户！请先注册")
		return
	}
}

// 查询书籍信息
func book1ListHandle(c *gin.Context) {
	//选择数据库
	//查数据
	//返回浏览器
	bookList, err := queryAllBook()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
		return
	}
	//以模板形式返回 , 因此需要再建一个 book_list.html 模板
	c.HTML(http.StatusOK, "book/book_list.html", gin.H{
		"code": 0,
		"data": bookList,
	})
}

// 查插入书籍数据
func newBookHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "book/new_book.html", nil)
}

// 增加书籍
func createBookHandle(c *gin.Context) {
	//增加新书,从form表单中提取数据
	titleVal := c.PostForm("title")
	priceVal := c.PostForm("price")
	//上面接受到的是string类型 , 存储前还需类型转换一下
	price, err := strconv.ParseFloat(priceVal, 64)
	if err != nil {
		fmt.Println("转换失败")
		return
	}
	//将提取的数据写入数据库 , 调用写好的insertAllBook()
	err = insertAllBook(titleVal, price)
	if err != nil {
		c.String(http.StatusOK, "插入数据失败")
		return
	}
	//到此数据插入成功
	//为了友好的交互跳转到书籍显示界面
	//使用重定向进行跳转
	c.Redirect(http.StatusMovedPermanently, "/book/list")
}

// 删除书籍
func deleteHandle(c *gin.Context) {
	//拿去query-string数据 , 然后根据不通数据删除指定编号的书籍
	idVal := c.Query("id")
	//将id 转换为整型
	id, err := strconv.ParseInt(idVal, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
	}
	//删除数据
	err = deleteBook(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
	}
	//重定向到书籍展示页面
	c.Redirect(http.StatusMovedPermanently, "/book/list")
}

// 显示书籍更新页面
func newHandle(c *gin.Context) {
	//拿去query-string数据 , 然后根据不通数据更新指定编号的书籍
	idVal := c.Query("id")
	//将id 转换为整型
	id, err := strconv.ParseInt(idVal, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
		return
	}
	book, err := querySingalBook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
		return
	}
	//向指定路径返送html模板
	c.HTML(http.StatusOK, "book/updatebook.html", book)
}

// 修改
func updateHandle(c *gin.Context) {
	//拿到from表单里面的update信息
	//增加新书 , 从from表单中提取数据
	titleVal := c.PostForm("title")
	priceVal := c.PostForm("price")
	idVal := c.PostForm("id")
	//上面接收到的是string类型 , 存储前还需要转换一下
	price, err := strconv.ParseFloat(priceVal, 64)
	if err != nil {
		fmt.Println("转换失败")
		return
	}
	id, err := strconv.ParseInt(idVal, 10, 64)
	if err != nil {
		fmt.Println("转换失败")
		return
	}
	err = updateBook(titleVal, price, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":  err.Error(),
			"code": 1,
		})
		return
	}
	//重定向到书籍展示界面
	c.Redirect(http.StatusMovedPermanently, "/book/list")
}
