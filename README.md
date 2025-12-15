# Owls Capital Discord Bot

A feature-rich Discord bot built with Go and discordgo.

## 🚀 Railway Deployment

### Prerequisites
- GitHub account with this repository
- Railway account ([sign up here](https://railway.app))
- Discord Bot Token
- Discord Guild (Server) ID

### Deployment Steps

1. **Go to Railway**: Visit [railway.app](https://railway.app) and sign in

2. **Create New Project**: 
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose `Raza006/owlscapital`

3. **Configure Environment Variables**:
   Railway will automatically detect the Dockerfile. Add these environment variables in the Railway dashboard:

   ```
   BOT_TOKEN=your_discord_bot_token_here
   GUILD_ID=your_discord_guild_id_here
   ```

   Optional variables (if using PostgreSQL):
   ```
   POSTGRES_USER=your_postgres_user
   POSTGRES_PASSWORD=your_postgres_password
   POSTGRES_DB=your_database_name
   POSTGRES_HOST=your_postgres_host
   POSTGRES_PORT=5432
   ```

4. **Deploy**: 
   - Railway will automatically build and deploy using the root `Dockerfile`
   - Watch the build logs to ensure successful deployment
   - The bot should come online in your Discord server

### Project Structure

```
.
├── Dockerfile                    # Root-level Dockerfile for Railway
├── railway.json                  # Railway configuration
├── .dockerignore                 # Docker ignore rules
├── services/
│   └── discordbot/              # Main bot service
│       ├── cmd/
│       │   ├── bot/             # Bot entry point
│       │   └── healthcheck/     # Health check utility
│       ├── internal/
│       │   ├── features/        # Bot features
│       │   ├── library/         # Reusable components
│       │   └── bot/             # Core bot logic
│       └── assets/              # Bot assets (images, etc.)
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `BOT_TOKEN` | Discord bot token from Discord Developer Portal | ✅ Yes |
| `GUILD_ID` | Your Discord server ID | ✅ Yes |
| `POSTGRES_USER` | PostgreSQL username | ❌ Optional |
| `POSTGRES_PASSWORD` | PostgreSQL password | ❌ Optional |
| `POSTGRES_DB` | PostgreSQL database name | ❌ Optional |
| `POSTGRES_HOST` | PostgreSQL host | ❌ Optional |
| `POSTGRES_PORT` | PostgreSQL port | ❌ Optional |

### Local Development

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Raza006/owlscapital.git
   cd owlscapital
   ```

2. **Set up environment variables**:
   ```bash
   cp ENV_TEMPLATE.txt services/discordbot/.env
   # Edit .env with your credentials
   ```

3. **Run with Docker**:
   ```bash
   docker build -t owls-bot .
   docker run --env-file services/discordbot/.env owls-bot
   ```

4. **Or run with Go**:
   ```bash
   cd services/discordbot
   go run ./cmd/bot
   ```

### Features

- `/ping` - Check bot latency
- Support panel system
- Event logging
- And more...

### Tech Stack

- **Language**: Go 1.25
- **Discord Library**: discordgo
- **Container**: Docker (multi-stage build)
- **Deployment**: Railway
- **Database**: PostgreSQL (optional)

### Support

For issues or questions, please open an issue on GitHub.

---

Built with 🦉 by Owls Capital

