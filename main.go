package main

import (
        "fmt"
        "time"
        "net/http"
        "context"
        "github.com/gin-gonic/gin"
        "github.com/go-redis/redis/v8"
)

type Data struct {
        Content string `json:"content"`
}

var (
        rds *redis.Client
        err error
)

func GetIpKey(c *gin.Context) string {
        ip := c.ClientIP()
        ip = "cut:" + ip
        return ip
}

func RedisInit() *redis.Client {
        rdb := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })
        result := rdb.Ping(context.Background())
        fmt.Println("redis ping:", result.Val())
        if result.Val()!="PONG"{
                // 连接有问题
                return nil
        }
        return rdb
}

func GetData(c *gin.Context) {

        ip := GetIpKey(c)
        content, err := rds.Get(c, ip).Result()
        if err != nil {
		content = ""
        }
        c.JSON(http.StatusOK, gin.H{
                "data": content,
        })

}

func SetData(c *gin.Context) {

        var data Data
        if err := c.ShouldBindJSON(&data); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
                return
        }
        if len(data.Content) > 4096 {
                c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "data too long, max: 4096"})
                return
        }
        content := data.Content
        ip := GetIpKey(c)

        // 将content参数存储在Redis中，以客户端IP地址作为键
        err := rds.Set(c, ip, content, 12 * time.Hour).Err()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                        "error": fmt.Sprintf("Failed to store data: %v", err),
                })
                return
        }
        // 返回成功响应
        c.JSON(http.StatusOK, gin.H{"m":fmt.Sprintf("")})

}


func main() {

        // 初始化redis
        rds = RedisInit()

        // 初始化Gin路由
        router := gin.Default()

        // 获取数据
        router.GET("/data", GetData)

        // 设置数据
        router.POST("/data", SetData)

        // 404
        router.NoRoute(func(c *gin.Context) {
            c.String(http.StatusNotFound, "404")
        })

        // 启动服务器
        if err := router.Run("127.0.0.1:8080"); err != nil {
                panic(err)
        }
}
