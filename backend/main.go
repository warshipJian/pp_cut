package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

// 定义Redis客户端
var rdb *redis.Client
var ctx = context.Background()

// 请求结构体
type ClipboardRequest struct {
	Content  string `json:"content"`
	Password string `json:"password,omitempty"`
}

func main() {
	// 初始化Redis连接
	initRedis()

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由
	r := gin.Default()

	// 健康检查接口
	r.GET("/c/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 获取剪切板内容接口
	r.GET("/c/get", func(c *gin.Context) {
		password := c.Query("password")
		clientIP := c.ClientIP()

		// 如果没有提供口令，使用客户端IP
		key := password
		if key == "" {
			key = clientIP
		}

		// 从Redis获取数据
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				// Key不存在，返回空字符串
				c.JSON(200, gin.H{
					"content": "",
				})
				return
			}
			// 其他错误
			log.Printf("Redis获取错误: %v", err)
			c.JSON(500, gin.H{
				"error": "服务器内部错误",
			})
			return
		}

		c.JSON(200, gin.H{
			"content": val,
		})
	})

	// 提交剪切板内容接口
	r.POST("/c/submit", func(c *gin.Context) {
		var req ClipboardRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": "无效的请求数据",
			})
			return
		}

		clientIP := c.ClientIP()

		// 如果没有提供口令，使用客户端IP
		key := req.Password
		if key == "" {
			key = clientIP
		}

		// 处理空字符串提交（清除数据）
		if req.Content == "" {
			_, err := rdb.Del(ctx, key).Result()
			if err != nil {
				log.Printf("Redis删除错误: %v", err)
				c.JSON(500, gin.H{
					"error": "服务器内部错误",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "数据已清除",
			})
			return
		}

		// 设置数据到Redis，过期时间为6小时
		err := rdb.Set(ctx, key, req.Content, 6*time.Hour).Err()
		if err != nil {
			log.Printf("Redis设置错误: %v", err)
			c.JSON(500, gin.H{
				"error": "服务器内部错误",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "数据提交成功",
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	log.Fatal(r.Run(":" + port))
}

// 初始化Redis连接
func initRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0, // 使用默认DB
	})

	// 测试Redis连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
	}

	log.Println("成功连接到Redis")
}
