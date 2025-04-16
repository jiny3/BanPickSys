package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/BanPickSys/game"
	"github.com/jiny3/BanPickSys/pkg"
	"github.com/jiny3/BanPickSys/service"
)

func SetupRouter(r *gin.Engine) {
	r.LoadHTMLGlob("static/*.html")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	bpGroup := r.Group("/bp")
	{
		bpGroup.Static("/static", "./static")
		bpGroup.GET("/", func(c *gin.Context) {
			// 读取 gameName name
			gameName := c.Query("game")
			if gameName == "" {
				gameName = "豹豹碰碰大作战"
			}
			// 启动 banpick
			bpID, err := service.NewBP(gameName, game.Handlers[gameName])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Redirect(http.StatusFound, fmt.Sprintf("/bp/%d", bpID))
		})

		bpIdGroup := bpGroup.Group("/:id")
		{
			bpIdGroup.GET("/", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				_, err = service.GetBP(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.HTML(http.StatusOK, "index.html", gin.H{
					"id": fmt.Sprintf("%d", bpID),
				})
			})
			bpIdGroup.POST("/submit", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				_, err = service.GetBP(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var req struct {
					EntryID  string `json:"entry_id"`
					PlayerID string `json:"player_id"`
				}
				err = c.ShouldBindJSON(&req)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
					return
				}
				eid, err := strconv.ParseInt(req.EntryID, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
					return
				}
				pid, err := strconv.ParseInt(req.PlayerID, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
					return
				}
				err = service.SendEvent(bpID, pid, eid)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"msg": "ok"})
			})
			bpIdGroup.POST("/join", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				_, err = service.GetBP(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				var req struct {
					PlayerID string `json:"id"`
					Role     string `json:"role"`
				}
				err = c.ShouldBindJSON(&req)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
					return
				}
				pid, err := strconv.ParseInt(req.PlayerID, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
					return
				}
				err = service.Join(bpID, pid, req.Role)
				if err != nil {
					c.JSON(http.StatusAlreadyReported, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"msg": "ok"})
			})
			bpIdGroup.POST("/leave", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				_, err = service.GetBP(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				var req struct {
					PlayerID string `json:"id"`
				}
				err = c.ShouldBindJSON(&req)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
					return
				}
				pid, err := strconv.ParseInt(req.PlayerID, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
					return
				}
				err = service.Leave(bpID, pid)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"msg": "ok"})
			})
			bpIdGroup.GET("/status", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				bp, err := service.GetBP(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"res": bp,
				})
			})
			bpIdGroup.GET("/entries", func(c *gin.Context) {
				id := c.Param("id")
				bpID, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bp ID"})
					return
				}
				res, err := service.GetEntries(bpID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"res": res})
			})
			bpIdGroup.GET("/result", func(c *gin.Context) {
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
			bpIdGroup.GET("/ws", func(c *gin.Context) {
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
	}
}
