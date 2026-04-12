package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://localhost:8080/api/v1"
const apiKey = "dev-key"

func main() {
	reader := bufio.NewReader(os.Stdin)
	playerID := "player1"

	fmt.Println("Welcome to TextGame CLI!")
	fmt.Print("Enter Player ID (default player1): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		playerID = input
	}

	for {
		resp := getPlayer(playerID)
		if resp == nil {
			fmt.Println("Failed to get player state. Is the server running?")
			return
		}

		player := resp["player"].(map[string]interface{})
		scene := resp["scene"].(map[string]interface{})

		fmt.Printf("\n--- %s | Level: %v (%v XP) ---\n", player["name"], player["level"], player["experience"])
		fmt.Printf("Skill Points: %v\n", player["skill_points"])

		resources, ok := player["resources"].(map[string]interface{})
		if ok {
			fmt.Print("Resources: ")
			for k, v := range resources {
				fmt.Printf("%s: %v ", k, v)
			}
			fmt.Println()
		}

		// Quests and current step (checked early to decide on location info)
		quests, _ := scene["quests"].([]interface{})
		currentStep, ok := scene["current_step"].(map[string]interface{})
		activeQuestID := ""
		if ok && currentStep != nil {
			progress, _ := player["quests_progress"].(map[string]interface{})
			for qID, p := range progress {
				prog := p.(map[string]interface{})
				if step, ok := prog["current_step"].(string); ok && step != "" {
					activeQuestID = qID
					break
				}
			}
		}

		// Only show Scene Info if NOT in an active quest
		if activeQuestID == "" {
			loc := scene["location"].(map[string]interface{})
			fmt.Printf("\n[LOCATION: %s]\n%s\n", loc["name"], loc["description"])

			conns, _ := scene["connections"].([]interface{})
			if len(conns) > 0 {
				fmt.Print("Available paths: ")
				for _, c := range conns {
					m := c.(map[string]interface{})
					fmt.Printf("[%s] ", m["id"])
				}
				fmt.Println()
			}

			props, _ := scene["properties"].([]interface{})
			if len(props) > 0 {
				fmt.Print("Properties here: ")
				for _, p := range props {
					m := p.(map[string]interface{})
					fmt.Printf("[%s: %v Credits] ", m["id"], m["price"])
				}
				fmt.Println()
			}
		} else {
			// In a quest: Just show quest info
			fmt.Printf("\n[ACTIVE QUEST: %s]\n%s\n", activeQuestID, currentStep["description"])
		}

		if activeQuestID != "" {
			fmt.Println("\nQuest Actions: [c]hoice, [s]tatus, [x]exit")
		} else {
			fmt.Println("\nActions: [m]ove, [b]uy, [q]uest start, [l]evel skill, [p]erk acquire, [s]tatus, [x]exit")
		}

		fmt.Print("> ")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "m":
			fmt.Print("Enter target location ID: ")
			target, _ := reader.ReadString('\n')
			movePlayer(playerID, strings.TrimSpace(target))
		case "b":
			fmt.Print("Enter property ID: ")
			prop, _ := reader.ReadString('\n')
			buyProperty(playerID, strings.TrimSpace(prop))
		case "q":
			if len(quests) > 0 {
				fmt.Print("Available Quests here: ")
				for _, q := range quests {
					m := q.(map[string]interface{})
					fmt.Printf("[%s] ", m["id"])
				}
				fmt.Println()
			} else {
				fmt.Println("No quests available at this location.")
			}
			fmt.Print("Enter quest ID to start: ")
			quest, _ := reader.ReadString('\n')
			startQuest(playerID, strings.TrimSpace(quest))
		case "c":
			if activeQuestID == "" {
				fmt.Println("No active quest")
				continue
			}
			// Show choices here before asking for index
			if currentStep != nil {
				choices, _ := currentStep["choices"].([]interface{})
				fmt.Println("Available Choices:")
				for i, ch := range choices {
					m := ch.(map[string]interface{})
					fmt.Printf("  [%d] %s\n", i, m["text"])
				}
			}
			fmt.Print("Enter choice index (0, 1, ...): ")
			idxStr, _ := reader.ReadString('\n')
			var idx int
			fmt.Sscanf(idxStr, "%d", &idx)
			makeChoice(playerID, activeQuestID, idx)
		case "l":
			fmt.Print("Enter skill ID to level: ")
			skill, _ := reader.ReadString('\n')
			spendSkillPoint(playerID, strings.TrimSpace(skill))
		case "p":
			fmt.Print("Enter perk ID to acquire: ")
			perk, _ := reader.ReadString('\n')
			acquirePerk(playerID, strings.TrimSpace(perk))
		case "s":
			continue
		case "x":
			return
		default:
			fmt.Println("Unknown action")
		}
	}
}

func getPlayer(id string) map[string]interface{} {
	req, _ := http.NewRequest("GET", baseURL+"/player/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func movePlayer(id, loc string) {
	body, _ := json.Marshal(map[string]string{"location_id": loc})
	post(id+"/move", body)
}

func buyProperty(id, prop string) {
	body, _ := json.Marshal(map[string]string{"property_id": prop})
	post(id+"/buy", body)
}

func startQuest(id, quest string) {
	body, _ := json.Marshal(map[string]string{"quest_id": quest})
	post(id+"/quest/start", body)
}

func makeChoice(id, quest string, index int) {
	body, _ := json.Marshal(map[string]interface{}{"quest_id": quest, "choice_index": index})
	post(id+"/quest/choice", body)
}

func spendSkillPoint(id, skill string) {
	body, _ := json.Marshal(map[string]string{"skill_id": skill})
	post(id+"/skill/spend", body)
}

func acquirePerk(id, perk string) {
	body, _ := json.Marshal(map[string]string{"perk_id": perk})
	post(id+"/perk/acquire", body)
}

func post(path string, body []byte) {
	req, _ := http.NewRequest("POST", baseURL+"/player/"+path, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Printf(">> Failure: %s\n", errResp.Error)
	} else {
		fmt.Println(">> Action successful!")
	}
}
