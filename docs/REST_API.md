# TextGame REST API Specification (v1.0)

**Final Production Build**

## Authentication
Every request must include:
`Authorization: Bearer <API_KEY>`

## Player Endpoints

### GET `/api/v1/player/:id`
Returns player state and their interactive `Scene`.
- **Response**:
    ```json
    {
      "player": { "id": "...", "name": "...", "resources": {"currency": 100}, "level": 1, ... },
      "scene": { "location": {...}, "connections": [...], "quests": [...], "current_step": {...} }
    }
    ```

### POST `/api/v1/player/:id/move`
`{ "location_id": "target" }`

### POST `/api/v1/player/:id/buy`
`{ "property_id": "id" }`

### POST `/api/v1/player/:id/quest/start`
`{ "quest_id": "id" }`

### POST `/api/v1/player/:id/quest/choice`
`{ "quest_id": "id", "choice_index": 0 }`

### GET `/api/v1/leaderboard/:metric`
Returns top 10 players by `level`, `xp`, or a resource (e.g., `currency`).
- **Response**: `[{"player_id": "...", "name": "...", "value": 100}, ...]`

## Admin Endpoints

### POST `/api/v1/admin/reload`
Hot-reload all JSON data.

### POST `/api/v1/admin/add/:type` (`location`, `quest`, `property`, `skill`, `perk`, `resource`)
Upsert an entity. Will overwrite if ID exists.

### DELETE `/api/v1/admin/remove/:type/:id`
Remove an entity.

### GET `/api/v1/admin/get/:type/:id`
Fetch raw entity data.
