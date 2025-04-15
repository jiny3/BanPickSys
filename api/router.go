package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/BanPickSys/model"
	"github.com/jiny3/BanPickSys/service"
)

func SetupRouter(r *gin.Engine) {
	r.LoadHTMLGlob("static/*.html")
	r.Static("/bp/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.GET("/bp", func(c *gin.Context) {
		// 启动一次banpick
		playNum := 2
		players := make([]model.Player, playNum)
		for i := range playNum {
			players[i] = model.Player{
				ID:     int64(i + 1),
				Name:   fmt.Sprintf("玩家%d", i+1),
				Banned: []model.Entry{},
				Picked: []model.Entry{},
			}
		}
		bpID, err := service.NewGame("豹豹碰碰大作战", players)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/bp/%d", bpID))
	})
	r.GET("/bp/:id", func(c *gin.Context) {
		// 获取游戏状态
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		_, err = service.GetGame(bpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"bp": fmt.Sprintf("%d", bpID),
		})
	})
	r.POST("/bp/:id", func(c *gin.Context) {
		// 玩家操作
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		_, err = service.GetGame(bpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var req struct {
			EntryID  int64 `json:"entry_id"`
			PlayerID int64 `json:"player_id"`
		}

		// Example JSON:
		// {
		//   "entry_id": 123,
		//   "player_id": 456
		// }
		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
			return
		}
		err = service.SendEvent(bpID, req.PlayerID, req.EntryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})
	r.GET("/bp/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		bp, err := service.GetGame(bpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"bp": bp,
		})
	})
	r.GET("/bp/:id/entries", func(c *gin.Context) {
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		bp, err := service.GetGame(bpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"entries": bp.Entries})
	})
	r.GET("/bp/:id/result", func(c *gin.Context) {
		// 获取游戏结果
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		res, err := service.GetResult(bpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"res": res})
	})
}
