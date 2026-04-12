package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"unix-supremacist.github.io/textgame/internal/game"
)

type Server struct {
	Engine *game.GameEngine
	Key    string
}

func NewServer(engine *game.GameEngine, apiKey string) *Server {
	return &Server{
		Engine: engine,
		Key:    apiKey,
	}
}

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKey != s.Key {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Server) SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/v1")
	auth.Use(s.AuthMiddleware())
	{
		auth.GET("/player/:id", s.HandleGetPlayer)
		auth.POST("/player/:id/move", s.HandleMove)
		auth.POST("/player/:id/buy", s.HandleBuyProperty)
		auth.POST("/player/:id/quest/start", s.HandleStartQuest)
		auth.POST("/player/:id/quest/choice", s.HandleQuestChoice)
		auth.POST("/player/:id/skill/spend", s.HandleSpendSkillPoint)
		auth.POST("/player/:id/perk/acquire", s.HandleAcquirePerk)

		admin := auth.Group("/admin")
		{
			admin.POST("/reload", s.HandleReload)
			admin.POST("/add/location", s.AddLocation)
			admin.POST("/add/property", s.AddProperty)
			admin.POST("/add/quest", s.AddQuest)
			admin.POST("/add/skill", s.AddSkill)
			admin.POST("/add/perk", s.AddPerk)
			admin.POST("/add/resource", s.AddResource)

			admin.DELETE("/remove/location/:id", s.RemoveLocation)
			admin.DELETE("/remove/property/:id", s.RemoveProperty)
			admin.DELETE("/remove/quest/:id", s.RemoveQuest)
			admin.DELETE("/remove/skill/:id", s.RemoveSkill)
			admin.DELETE("/remove/perk/:id", s.RemovePerk)
			admin.DELETE("/remove/resource/:id", s.RemoveResource)

			admin.GET("/get/location/:id", s.GetLocation)
			admin.GET("/get/property/:id", s.GetProperty)
			admin.GET("/get/quest/:id", s.GetQuest)
			admin.GET("/get/skill/:id", s.GetSkill)
			admin.GET("/get/perk/:id", s.GetPerk)
			admin.GET("/get/resource/:id", s.GetResource)
		}

		// Public/Player info
		auth.GET("/leaderboard/:metric", s.HandleGetLeaderboard)
	}

	return r
}
