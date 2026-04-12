package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type GameEngine struct {
	sync.RWMutex
	Locations     map[string]Location
	Properties    map[string]Property
	Quests        map[string]Quest
	Skills        map[string]Skill
	Perks         map[string]Perk
	Players       map[string]*Player
	Resources     map[Resource]ResourceDef
	StartingStats map[string]interface{} `json:"-"`
	DataPath      string
}

func NewGameEngine(dataPath string) *GameEngine {
	return &GameEngine{
		Locations:  make(map[string]Location),
		Properties: make(map[string]Property),
		Quests:     make(map[string]Quest),
		Skills:     make(map[string]Skill),
		Perks:      make(map[string]Perk),
		Players:    make(map[string]*Player),
		Resources:  make(map[Resource]ResourceDef),
		DataPath:   dataPath,
	}
}

func (e *GameEngine) LoadData() error {
	e.Lock()
	defer e.Unlock()

	if err := e.loadJSON("locations.json", &e.Locations); err != nil {
		return err
	}
	if err := e.loadJSON("properties.json", &e.Properties); err != nil {
		return err
	}
	if err := e.loadJSON("quests.json", &e.Quests); err != nil {
		return err
	}
	if err := e.loadJSON("skills.json", &e.Skills); err != nil {
		return err
	}
	if err := e.loadJSON("perks.json", &e.Perks); err != nil {
		return err
	}
	if err := e.loadJSON("resources.json", &e.Resources); err != nil {
		return err
	}
	if err := e.loadJSON("starting_stats.json", &e.StartingStats); err != nil {
		return err
	}
	if err := e.loadJSON("players.json", &e.Players); err != nil {
		// If players.json doesn't exist, we just start fresh
		log.Printf("Info: players.json not found, starting with empty player base")
	}

	return nil
}

func (e *GameEngine) loadJSON(filename string, target interface{}) error {
	path := fmt.Sprintf("%s/%s", e.DataPath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, target)
}

func (e *GameEngine) SaveData() error {
	e.RLock()
	defer e.RUnlock()

	if err := e.saveJSON("locations.json", e.Locations); err != nil {
		return err
	}
	if err := e.saveJSON("properties.json", e.Properties); err != nil {
		return err
	}
	if err := e.saveJSON("quests.json", e.Quests); err != nil {
		return err
	}
	if err := e.saveJSON("skills.json", e.Skills); err != nil {
		return err
	}
	if err := e.saveJSON("perks.json", e.Perks); err != nil {
		return err
	}
	if err := e.saveJSON("resources.json", e.Resources); err != nil {
		return err
	}
	if err := e.saveJSON("players.json", e.Players); err != nil {
		return err
	}

	return nil
}

func (e *GameEngine) saveJSON(filename string, data interface{}) error {
	path := fmt.Sprintf("%s/%s", e.DataPath, filename)
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func (e *GameEngine) SavePlayers() error {
	e.RLock()
	defer e.RUnlock()
	return e.saveJSON("players.json", e.Players)
}

func (e *GameEngine) GetPlayer(id string) *Player {
	e.RLock()
	defer e.RUnlock()
	return e.Players[id]
}

func (e *GameEngine) SpawnPlayer(id, name string) *Player {
	e.Lock()
	defer e.Unlock()

	p := &Player{
		ID:             id,
		Name:           name,
		LocationID:     "city_center",
		Resources:      make(map[Resource]int),
		Skills:         make(map[string]int),
		SkillPoints:    0,
		Level:          1,
		Experience:     0,
		Inventory:      make([]string, 0),
		QuestsProgress: make(map[string]*PlayerQuestProgress),
		LastSeen:       time.Now(),
	}

	if e.StartingStats != nil {
		if loc, ok := e.StartingStats["location_id"].(string); ok {
			p.LocationID = loc
		}
		if res, ok := e.StartingStats["resources"].(map[string]interface{}); ok {
			for k, v := range res {
				p.Resources[Resource(k)] = int(v.(float64))
			}
		}
		if exp, ok := e.StartingStats["experience"].(float64); ok {
			p.Experience = int(exp)
		}
		if lvl, ok := e.StartingStats["level"].(float64); ok {
			p.Level = int(lvl)
		}
	}

	e.Players[id] = p
	go e.SavePlayers() // Save immediately
	return p
}

func (e *GameEngine) MovePlayer(playerID, targetLocationID string) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	if _, ok := e.Locations[targetLocationID]; !ok {
		fmt.Printf("DEBUG: Move attempt to non-existent location: %s\n", targetLocationID)
		return fmt.Errorf("target location %s does not exist", targetLocationID)
	}

	currentLoc, ok := e.Locations[player.LocationID]
	if !ok {
		player.LocationID = targetLocationID
		return nil
	}

	connected := false
	for _, conn := range currentLoc.Connections {
		if conn == targetLocationID {
			connected = true
			break
		}
	}

	if !connected {
		return fmt.Errorf("location %s is not connected to %s", targetLocationID, player.LocationID)
	}

	player.LocationID = targetLocationID
	player.LastSeen = time.Now()
	go e.SavePlayers()
	return nil
}

func (e *GameEngine) BuyProperty(playerID, propertyID string) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	prop, ok := e.Properties[propertyID]
	if !ok {
		return fmt.Errorf("property not found")
	}

	for _, pID := range player.Properties {
		if pID == propertyID {
			return fmt.Errorf("property already owned")
		}
	}

	loc := e.Locations[player.LocationID]
	atLocation := false
	for _, pID := range loc.Properties {
		if pID == propertyID {
			atLocation = true
			break
		}
	}
	if !atLocation {
		return fmt.Errorf("property not available at this location")
	}

	if player.Resources[Currency] < prop.Price {
		return fmt.Errorf("insufficient currency: have %d, need %d", player.Resources[Currency], prop.Price)
	}

	player.Resources[Currency] -= prop.Price
	player.Properties = append(player.Properties, propertyID)
	player.LastSeen = time.Now()

	return nil
}

func (e *GameEngine) StartQuest(playerID, questID string) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	quest, ok := e.Quests[questID]
	if !ok {
		return fmt.Errorf("quest not found")
	}

	if player.LocationID != quest.LocationID {
		return fmt.Errorf("player not at quest location")
	}

	progress, exists := player.QuestsProgress[questID]
	if exists && progress.CurrentStep != "" {
		return fmt.Errorf("quest already in progress")
	}

	if exists && !quest.Repeatable {
		return fmt.Errorf("quest is not repeatable")
	}

	if exists && quest.Repeatable && time.Since(progress.LastComplete) < quest.Cooldown {
		return fmt.Errorf("quest is on cooldown: %v remaining", quest.Cooldown-time.Since(progress.LastComplete))
	}

	if !exists {
		progress = &PlayerQuestProgress{QuestID: questID}
		player.QuestsProgress[questID] = progress
	}

	if len(quest.Steps) == 0 {
		return fmt.Errorf("quest has no steps")
	}

	progress.CurrentStep = quest.Steps[0].ID
	player.LastSeen = time.Now()
	go e.SavePlayers()
	return nil
}

func (e *GameEngine) MakeQuestChoice(playerID, questID string, choiceIndex int) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	quest, ok := e.Quests[questID]
	if !ok {
		return fmt.Errorf("quest not found")
	}

	progress, ok := player.QuestsProgress[questID]
	if !ok || progress.CurrentStep == "" {
		return fmt.Errorf("quest not started")
	}

	var currentStep *QuestStep
	for i := range quest.Steps {
		if quest.Steps[i].ID == progress.CurrentStep {
			currentStep = &quest.Steps[i]
			break
		}
	}

	if currentStep == nil {
		return fmt.Errorf("internal error: current step not found in quest")
	}

	if choiceIndex < 0 || choiceIndex >= len(currentStep.Choices) {
		return fmt.Errorf("invalid choice index")
	}

	choice := currentStep.Choices[choiceIndex]

	// Check Requirements
	for skillID, reqLevel := range choice.Requirements {
		if player.Skills[skillID] < reqLevel {
			return fmt.Errorf("requirement %s level %d not met", skillID, reqLevel)
		}
	}

	// Check Perks
	for _, perkID := range choice.PerkReq {
		hasPerk := false
		for _, pID := range player.Perks {
			if pID == perkID {
				hasPerk = true
				break
			}
		}
		if !hasPerk {
			return fmt.Errorf("perk %s required", perkID)
		}
	}

	// Check Items
	for _, item := range choice.ItemReq {
		hasItem := false
		for _, invItem := range player.Inventory {
			if invItem == item {
				hasItem = true
				break
			}
		}
		if !hasItem {
			return fmt.Errorf("item %s required", item)
		}
	}

	// Apply Rewards
	for res, amount := range choice.Rewards {
		player.Resources[res] += amount
	}
	if choice.ExpReward > 0 {
		e.AddExperience(player, choice.ExpReward)
	}

	// Transition
	if choice.NextStepID == "" {
		progress.CurrentStep = ""
		progress.LastComplete = time.Now()
		if quest.Experience > 0 {
			e.AddExperience(player, quest.Experience)
		}
	} else {
		progress.CurrentStep = choice.NextStepID
	}

	player.LastSeen = time.Now()
	go e.SavePlayers()
	return nil
}

func (e *GameEngine) AddExperience(player *Player, amount int) {
	player.Experience += amount
	for player.Experience >= e.ExpForLevel(player.Level+1) {
		player.Level++
		player.SkillPoints += 5
		fmt.Printf("DEBUG: Player %s leveled up to %d!\n", player.ID, player.Level)
	}
}

func (e *GameEngine) ExpForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	return (level - 1) * 1000
}

func (e *GameEngine) SpendSkillPoint(playerID, skillID string) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	if player.SkillPoints <= 0 {
		return fmt.Errorf("no skill points available")
	}

	if _, ok := e.Skills[skillID]; !ok {
		return fmt.Errorf("skill %s not found", skillID)
	}

	player.Skills[skillID]++
	player.SkillPoints--
	player.LastSeen = time.Now()
	go e.SavePlayers()
	return nil
}

func (e *GameEngine) AcquirePerk(playerID, perkID string) error {
	e.Lock()
	defer e.Unlock()

	player, ok := e.Players[playerID]
	if !ok {
		return fmt.Errorf("player not found")
	}

	perk, ok := e.Perks[perkID]
	if !ok {
		return fmt.Errorf("perk not found")
	}

	for _, pID := range player.Perks {
		if pID == perkID {
			return fmt.Errorf("perk already owned")
		}
	}

	for skillID, reqLevel := range perk.Requirements {
		if player.Skills[skillID] < reqLevel {
			return fmt.Errorf("requirement %s level %d not met", skillID, reqLevel)
		}
	}

	player.Perks = append(player.Perks, perkID)
	player.LastSeen = time.Now()
	go e.SavePlayers()
	return nil
}

func (e *GameEngine) GetScene(playerID string) (*Scene, error) {
	e.RLock()
	defer e.RUnlock()

	player, ok := e.Players[playerID]
	if !ok {
		return nil, fmt.Errorf("player not found")
	}

	loc, ok := e.Locations[player.LocationID]
	if !ok {
		return nil, fmt.Errorf("location %s not found", player.LocationID)
	}

	scene := &Scene{
		Location:    loc,
		Connections: make([]Location, 0),
		Properties:  make([]Property, 0),
		Quests:      make([]Quest, 0),
	}

	for _, connID := range loc.Connections {
		if c, ok := e.Locations[connID]; ok {
			scene.Connections = append(scene.Connections, c)
		}
	}

	for _, propID := range loc.Properties {
		if p, ok := e.Properties[propID]; ok {
			scene.Properties = append(scene.Properties, p)
		}
	}

	for _, questID := range loc.Quests {
		if q, ok := e.Quests[questID]; ok {
			// Filter out completed unrepeatable quests
			if prog, exists := player.QuestsProgress[questID]; exists {
				if !q.Repeatable && prog.CurrentStep == "" && !prog.LastComplete.IsZero() {
					continue
				}
			}
			scene.Quests = append(scene.Quests, q)
		}
	}

	for qID, prog := range player.QuestsProgress {
		if prog.CurrentStep != "" {
			if q, ok := e.Quests[qID]; ok {
				for _, step := range q.Steps {
					if step.ID == prog.CurrentStep {
						scene.CurrentStep = &step
						break
					}
				}
			}
			break
		}
	}

	return scene, nil
}

func (e *GameEngine) AddLocation(loc Location) error {
	e.Lock()
	e.Locations[loc.ID] = loc
	e.Unlock()
	return e.saveJSON("locations.json", e.Locations)
}

func (e *GameEngine) AddQuest(q Quest) error {
	e.Lock()
	e.Quests[q.ID] = q
	e.Unlock()
	return e.saveJSON("quests.json", e.Quests)
}

func (e *GameEngine) AddProperty(p Property) error {
	e.Lock()
	e.Properties[p.ID] = p
	e.Unlock()
	return e.saveJSON("properties.json", e.Properties)
}

func (e *GameEngine) AddSkill(sk Skill) error {
	e.Lock()
	e.Skills[sk.ID] = sk
	e.Unlock()
	return e.saveJSON("skills.json", e.Skills)
}

func (e *GameEngine) AddPerk(p Perk) error {
	e.Lock()
	e.Perks[p.ID] = p
	e.Unlock()
	return e.saveJSON("perks.json", e.Perks)
}

func (e *GameEngine) AddResource(r ResourceDef) error {
	e.Lock()
	e.Resources[r.ID] = r
	e.Unlock()
	return e.saveJSON("resources.json", e.Resources)
}

func (e *GameEngine) RemoveLocation(id string) error {
	e.Lock()
	delete(e.Locations, id)
	e.Unlock()
	return e.saveJSON("locations.json", e.Locations)
}

func (e *GameEngine) RemoveQuest(id string) error {
	e.Lock()
	delete(e.Quests, id)
	e.Unlock()
	return e.saveJSON("quests.json", e.Quests)
}

func (e *GameEngine) RemoveProperty(id string) error {
	e.Lock()
	delete(e.Properties, id)
	e.Unlock()
	return e.saveJSON("properties.json", e.Properties)
}

func (e *GameEngine) RemoveSkill(id string) error {
	e.Lock()
	delete(e.Skills, id)
	e.Unlock()
	return e.saveJSON("skills.json", e.Skills)
}

func (e *GameEngine) RemovePerk(id string) error {
	e.Lock()
	delete(e.Perks, id)
	e.Unlock()
	return e.saveJSON("perks.json", e.Perks)
}

func (e *GameEngine) RemoveResource(id Resource) error {
	e.Lock()
	delete(e.Resources, id)
	e.Unlock()
	return e.saveJSON("resources.json", e.Resources)
}

func (e *GameEngine) AdminGetLocation(id string) (Location, bool) {
	e.RLock()
	defer e.RUnlock()
	loc, ok := e.Locations[id]
	return loc, ok
}

func (e *GameEngine) AdminGetQuest(id string) (Quest, bool) {
	e.RLock()
	defer e.RUnlock()
	q, ok := e.Quests[id]
	return q, ok
}

func (e *GameEngine) AdminGetProperty(id string) (Property, bool) {
	e.RLock()
	defer e.RUnlock()
	p, ok := e.Properties[id]
	return p, ok
}

func (e *GameEngine) AdminGetSkill(id string) (Skill, bool) {
	e.RLock()
	defer e.RUnlock()
	s, ok := e.Skills[id]
	return s, ok
}

func (e *GameEngine) AdminGetPerk(id string) (Perk, bool) {
	e.RLock()
	defer e.RUnlock()
	p, ok := e.Perks[id]
	return p, ok
}

func (e *GameEngine) AdminGetResource(id Resource) (ResourceDef, bool) {
	e.RLock()
	defer e.RUnlock()
	r, ok := e.Resources[id]
	return r, ok
}

type LeaderboardEntry struct {
	PlayerID string `json:"player_id"`
	Name     string `json:"name"`
	Value    int    `json:"value"`
}

func (e *GameEngine) GetLeaderboard(metric string) []LeaderboardEntry {
	e.RLock()
	defer e.RUnlock()

	entries := make([]LeaderboardEntry, 0, len(e.Players))
	for id, p := range e.Players {
		val := 0
		switch metric {
		case "level":
			val = p.Level
		case "xp":
			val = p.Experience
		default:
			// Check if it's a resource
			if rVal, ok := p.Resources[Resource(metric)]; ok {
				val = rVal
			}
		}
		entries = append(entries, LeaderboardEntry{
			PlayerID: id,
			Name:     p.Name,
			Value:    val,
		})
	}

	// Simple sort (descending)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Value > entries[i].Value {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Limit to top 10
	if len(entries) > 10 {
		return entries[:10]
	}
	return entries
}

func (e *GameEngine) RunProductionTicker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			e.Lock()
			for _, player := range e.Players {
				for _, propID := range player.Properties {
					if prop, ok := e.Properties[propID]; ok {
						for res, amount := range prop.ResourceProduction {
							player.Resources[res] += amount
						}
					}
				}
			}
			e.Unlock()
		}
	}()
}
