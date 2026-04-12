package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unix-supremacist.github.io/textgame/internal/game"
)

func (s *Server) HandleGetPlayer(c *gin.Context) {
	id := c.Param("id")
	player := s.Engine.GetPlayer(id)
	if player == nil {
		player = s.Engine.SpawnPlayer(id, "Adventurer")
	}

	scene, _ := s.Engine.GetScene(id)

	c.JSON(http.StatusOK, gin.H{
		"player": player,
		"scene":  scene,
	})
}

func (s *Server) HandleMove(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		LocationID string `json:"location_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.MovePlayer(id, req.LocationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleBuyProperty(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PropertyID string `json:"property_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.BuyProperty(id, req.PropertyID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleStartQuest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		QuestID string `json:"quest_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.StartQuest(id, req.QuestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleQuestChoice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		QuestID     string `json:"quest_id"`
		ChoiceIndex int    `json:"choice_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.MakeQuestChoice(id, req.QuestID, req.ChoiceIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleSpendSkillPoint(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		SkillID string `json:"skill_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.SpendSkillPoint(id, req.SkillID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleAcquirePerk(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PerkID string `json:"perk_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Engine.AcquirePerk(id, req.PerkID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s.Engine.GetPlayer(id))
}

func (s *Server) HandleReload(c *gin.Context) {
	if err := s.Engine.LoadData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "data reloaded"})
}

func (s *Server) AddLocation(c *gin.Context) {
	var loc game.Location
	if err := c.ShouldBindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddLocation(loc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (s *Server) AddQuest(c *gin.Context) {
	var q game.Quest
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddQuest(q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (s *Server) AddProperty(c *gin.Context) {
	var p game.Property
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddProperty(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) AddSkill(c *gin.Context) {
	var sk game.Skill
	if err := c.ShouldBindJSON(&sk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddSkill(sk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sk)
}

func (s *Server) AddPerk(c *gin.Context) {
	var p game.Perk
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddPerk(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) AddResource(c *gin.Context) {
	var r game.ResourceDef
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Engine.AddResource(r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (s *Server) RemoveLocation(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemoveLocation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "location removed", "id": id})
}

func (s *Server) RemoveQuest(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemoveQuest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "quest removed", "id": id})
}

func (s *Server) RemoveProperty(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemoveProperty(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "property removed", "id": id})
}

func (s *Server) RemoveSkill(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemoveSkill(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "skill removed", "id": id})
}

func (s *Server) RemovePerk(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemovePerk(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "perk removed", "id": id})
}

func (s *Server) RemoveResource(c *gin.Context) {
	id := c.Param("id")
	if err := s.Engine.RemoveResource(game.Resource(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "resource removed", "id": id})
}

func (s *Server) GetLocation(c *gin.Context) {
	id := c.Param("id")
	loc, ok := s.Engine.AdminGetLocation(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (s *Server) GetQuest(c *gin.Context) {
	id := c.Param("id")
	q, ok := s.Engine.AdminGetQuest(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "quest not found"})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (s *Server) GetProperty(c *gin.Context) {
	id := c.Param("id")
	p, ok := s.Engine.AdminGetProperty(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "property not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) GetSkill(c *gin.Context) {
	id := c.Param("id")
	sk, ok := s.Engine.AdminGetSkill(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
		return
	}
	c.JSON(http.StatusOK, sk)
}

func (s *Server) GetPerk(c *gin.Context) {
	id := c.Param("id")
	p, ok := s.Engine.AdminGetPerk(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "perk not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) GetResource(c *gin.Context) {
	id := c.Param("id")
	r, ok := s.Engine.AdminGetResource(game.Resource(id))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (s *Server) HandleGetLeaderboard(c *gin.Context) {
	metric := c.Param("metric")
	leaderboard := s.Engine.GetLeaderboard(metric)
	c.JSON(http.StatusOK, leaderboard)
}
