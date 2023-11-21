package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// 创建留言表
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT, time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, is_active INTEGER)")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	// 添加留言
	r.POST("/messages", func(c *gin.Context) {
		var json struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec("INSERT INTO messages (content,is_active) VALUES (?,1)", json.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		c.JSON(http.StatusOK, gin.H{"id": id, "content": json.Content})
	})

	// 获取单个留言
	r.GET("/messages/:id", func(c *gin.Context) {
		id := c.Param("id")
		var content string
		var time string
		var active int
		err := db.QueryRow("SELECT content,time,is_active FROM messages WHERE id = ?", id).Scan(&content, &time, &active)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id, "content": content, "time": time, "is_active": active})
	})

	r.GET("/messages", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM messages")
		if err != nil {
			// 处理查询错误
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer func(rows *sql.Rows) {
			_ = rows.Close()
		}(rows)
		// 遍历结果集并获取留言内容
		var messages []gin.H
		for rows.Next() {
			var id int
			var content string
			var time string
			var active int
			err = rows.Scan(&id, &content, &time, &active)
			if err != nil {
				// 处理扫描错误
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			message := gin.H{"id": id, "content": content, "time": time, "is_active": active}
			messages = append(messages, message)
		}
		// 返回获取到的留言内容
		c.JSON(200, gin.H{"messages": messages})
	})

	// 更新留言
	r.PUT("/messages/:id", func(c *gin.Context) {
		id := c.Param("id")
		var json struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE messages SET content = ?,time = CURRENT_TIMESTAMP,is_active = 1 WHERE id = ?", json.Content, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id, "content": json.Content, "is_active": 1})
	})

	// 删除留言
	r.DELETE("/messages/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("UPDATE messages SET is_active = 0 WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Message deleted successfully"})
	})

	_ = r.Run(":10070")

}
