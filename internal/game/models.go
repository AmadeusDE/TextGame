package game

import "time"

type Resource string

const (
	Currency Resource = "currency"
)

type ResourceDef struct {
	ID          Resource `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
}

type Location struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Connections []string `json:"connections"`
	Properties  []string `json:"properties"`
	Quests      []string `json:"quests"`
}

type Property struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	Price              int              `json:"price"`
	ResourceProduction map[Resource]int `json:"resource_production"`
	ProductionInterval time.Duration    `json:"production_interval"`
}

type Quest struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	LocationID  string        `json:"location_id"`
	Repeatable  bool          `json:"repeatable"`
	Cooldown    time.Duration `json:"cooldown"`
	Steps       []QuestStep   `json:"steps"`
	Experience  int           `json:"experience"`
}

type QuestStep struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Choices     []Choice `json:"choices"`
}

type Choice struct {
	Text         string           `json:"text"`
	Requirements map[string]int   `json:"requirements"` // skill_id -> min_level
	ItemReq      []string         `json:"item_req"`     // item names
	PerkReq      []string         `json:"perk_req"`     // perk ids
	Rewards      map[Resource]int `json:"rewards"`
	NextStepID   string           `json:"next_step_id"` // "" means end of quest
	ExpReward    int              `json:"exp_reward"`
}

type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Perk struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Requirements map[string]int `json:"requirements"` // skill_id -> min_level
}

type PlayerQuestProgress struct {
	QuestID      string    `json:"quest_id"`
	CurrentStep  string    `json:"current_step"` // "" if not started or complete
	LastComplete time.Time `json:"last_complete"`
}

type Player struct {
	ID             string                          `json:"id"`
	Name           string                          `json:"name"`
	LocationID     string                          `json:"location_id"`
	Resources      map[Resource]int                `json:"resources"`
	Skills         map[string]int                  `json:"skills"` // skill_id -> level
	SkillPoints    int                             `json:"skill_points"`
	Level          int                             `json:"level"`
	Experience     int                             `json:"experience"`
	Perks          []string                        `json:"perks"`
	Inventory      []string                        `json:"inventory"`
	Properties     []string                        `json:"properties"`
	QuestsProgress map[string]*PlayerQuestProgress `json:"quests_progress"`
	LastSeen       time.Time                       `json:"last_seen"`
}

type Scene struct {
	Location    Location   `json:"location"`
	Connections []Location `json:"connections"`
	Properties  []Property `json:"properties"`
	Quests      []Quest    `json:"quests"`
	CurrentStep *QuestStep `json:"current_step,omitempty"`
}
