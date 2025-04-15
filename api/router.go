package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/BanPickSys/model"
	"github.com/jiny3/BanPickSys/service"
)

func SetupRouter(r *gin.Engine) {
	r.LoadHTMLGlob("static/*.html")
	r.Static("/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/game", func(c *gin.Context) {
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
		gameID, err := service.NewGame("豹豹碰碰大作战", players)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// 让前端定向到game/gameID
		c.Redirect(302, fmt.Sprintf("/game/%d", gameID))
	})
	r.GET("/game/:id", func(c *gin.Context) {
		// 获取游戏状态
		id := c.Param("id")
		gameID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid game ID"})
			return
		}
		game, err := service.GetGame(gameID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"game": game.ID,
		})
	})
	r.POST("/game/:id", func(c *gin.Context) {
		// 玩家操作
		id := c.Param("id")
		gameID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid game ID"})
			return
		}
		_, err = service.GetGame(gameID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
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
			c.JSON(400, gin.H{"error": "invalid entry"})
			return
		}
		err = service.SendEvent(gameID, req.PlayerID, req.EntryID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"msg": "ok"})
	})
	r.GET("/game/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		gameID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid game ID"})
			return
		}
		game, err := service.GetGame(gameID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"game": game,
		})
	})
	r.GET("/game/:id/entries", func(c *gin.Context) {
		id := c.Param("id")
		gameID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid game ID"})
			return
		}
		game, err := service.GetGame(gameID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"entries": game.Entries})
	})
	r.GET("/game/:id/result", func(c *gin.Context) {
		// 获取游戏结果
		id := c.Param("id")
		gameID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid game ID"})
			return
		}
		res, err := service.GetResult(gameID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"res": res})
	})
}
