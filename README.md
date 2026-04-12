# TextGame: The Neon Underworld Engine

Welcome to the core engine of **TextGame**, a high-performance, data-driven RPG backend designed for immersive terminal or bot-based experiences. This is not just a game; it is a simulation of a cyber-noir reality where skills, perks, and resources are the only currency that matters.

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- [Optional] Docker

### Interactive Development
1. **Server**: `go run cmd/api/main.go`
2. **Player CLI**: `go run cmd/cli/main.go`
3. **Admin CLI**: `go run cmd/admin-cli/main.go`

### Production Deployment (Docker)
```bash
docker build -t textgame .
docker run -p 8080:8080 -v $(pwd)/data:/app/data textgame
```

## 📖 Documentation
- [REST API Specification](file:///home/unix/src/TextGame/docs/REST_API.md): Detailed endpoints for players and admins.
- [JSON Data Format](file:///home/unix/src/TextGame/docs/JSON_FORMAT.md): Schema for world-building.

## ⚖️ Ethics & Conduct
We take our community standards very seriously. Please refer to our [CODE_OF_CONDUCT.md](file:///home/unix/src/TextGame/CODE_OF_CONDUCT.md) for our strict guidelines on interaction and existence within the machine.

## 📜 License
This project is licensed under the [LICENSE](file:///home/unix/src/TextGame/LICENSE) provided in this repository. All users, contributors, and operators of the grid MUST comply with the terms specified. Failure to adhere to the license is a breach of the integrity of the system and will be prosecuted accordingly.