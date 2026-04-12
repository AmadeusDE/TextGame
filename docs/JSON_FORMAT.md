# TextGame JSON Data Format

This document describes the structure of the JSON files in the `data/` directory.

## `locations.json`
A map of `location_id` to Location object.
```json
{
  "city_center": {
    "id": "city_center",
    "name": "City Center",
    "description": "Heart of the city.",
    "connections": ["slums"],
    "properties": ["nightclub"],
    "quests": ["neon_heist"]
  }
}
```

## `quests.json`
A map of `quest_id` to Quest object.
```json
{
  "neon_heist": {
    "id": "neon_heist",
    "name": "The Neon Heist",
    "description": "A high-stakes robbery.",
    "location_id": "city_center",
    "repeatable": false,
    "cooldown": 0,
    "experience": 500,
    "steps": [
      {
        "id": "entry",
        "description": "How do you enter?",
        "choices": [
          {
            "text": "Hack in.",
            "requirements": { "hacking": 2 },
            "next_step_id": "terminal",
            "exp_reward": 100
          }
        ]
      }
    ]
  }
}
```

### Choice Object fields:
- `text`: Display text for the choice.
- `requirements`: Map of `skill_id` to minimum level required.
- `item_req`: List of item names required in inventory.
- `perk_req`: List of perk IDs required.
- `rewards`: Map of `resource_id` to amount gained/lost.
- `next_step_id`: ID of the next step. Empty string ends the quest.
- `exp_reward`: Experience gained immediately upon choosing.

## `properties.json`
A map of `property_id` to Property object.
- `production_interval`: Nanoseconds (e.g., `60000000000` for 1 minute).

## `skills.json` & `perks.json`
Simple maps of IDs to definition objects. Perks include `requirements` (skill levels).

## `starting_stats.json`
Defines the initial state for new players.
```json
{
  "location_id": "city_center",
  "resources": { "currency": 10000 },
  "skills": { "hacking": 1 },
  "inventory": ["pickaxe"]
}
```
