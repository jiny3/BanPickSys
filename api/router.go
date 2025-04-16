package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/BanPickSys/pkg"
	"github.com/jiny3/BanPickSys/service"
)

func SetupRouter(r *gin.Engine) {
	r.LoadHTMLGlob("static/*.html")
	r.Static("/bp/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.GET("/bp", func(c *gin.Context) {
		// 读取 game name
		game := c.Query("game")
		if game == "" {
			game = "豹豹碰碰大作战"
		}
		// 启动 banpick
		bpID, err := service.NewGame(game, GameHandler[game])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/bp/%d", bpID))
	})
	r.GET("/bp/:id", func(c *gin.Context) {
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
	// TODO: 改为通过ws交互
	r.POST("/bp/:id", func(c *gin.Context) {
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
	r.GET("/bp/:id/ws", func(c *gin.Context) {
		id := c.Param("id")
		bpID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
			return
		}
		wsConn, err := pkg.WsUpgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer wsConn.Close()

		service.WsHandler(wsConn, bpID)
	})
}
