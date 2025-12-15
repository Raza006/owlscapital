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

3. **Add PostgreSQL Database**:
   - In your Railway project, click `+ New`
   - Select `Database` → `PostgreSQL`
   - Railway will automatically create the database and set environment variables

4. **Configure Environment Variables**:
   Railway will automatically detect the Dockerfile and set most variables. You only need to add:

   ```
   BOT_TOKEN=your_discord_bot_token_here
   GUILD_ID=your_discord_guild_id_here
   ```

   The bot will automatically detect these PostgreSQL variables from Railway:
   ```
   POSTGRES_USER (auto-set by Railway)
   POSTGRES_PASSWORD (auto-set by Railway)
   POSTGRES_DB (auto-set by Railway)
   POSTGRES_HOST (auto-set by Railway)
   POSTGRES_PORT (auto-set by Railway)
   ```

5. **Deploy**: 
   - Railway will automatically build and deploy using the root `Dockerfile`
   - Watch the build logs to ensure successful deployment
   - **The bot will automatically create all database tables on first startup!** ✨
   - The bot should come online in your Discord server

### 🔄 Automatic Database Migrations

The bot uses **GORM AutoMigrate** to automatically create and update database tables. When the bot starts, it will:

1. ✅ Connect to your PostgreSQL database
2. ✅ Create all required tables if they don't exist
3. ✅ Update existing tables if the schema changed
4. ✅ Insert default settings for the Freemium system

**No manual SQL scripts needed!** The bot handles everything automatically.

#### Database Tables Created:

**Ambassador System:**
- `ambassador_applicants` - Pending ambassador applications
- `ambassadors` - Active ambassadors with stats
- `claims` - User claim tracking
- `conversions` - Conversion history
- `payout_requests` - Payout management
- `balance_transactions` - Balance audit log
- `weekly_leaderboard` - Weekly performance tracking

**Freemium System:**
- `freemium_modules` - Educational modules
- `freemium_lessons` - Lesson content
- `freemium_access_log` - Access analytics
- `freemium_settings` - System configuration

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
| `POSTGRES_USER` | PostgreSQL username | ✅ Yes (auto-set by Railway) |
| `POSTGRES_PASSWORD` | PostgreSQL password | ✅ Yes (auto-set by Railway) |
| `POSTGRES_DB` | PostgreSQL database name | ✅ Yes (auto-set by Railway) |
| `POSTGRES_HOST` | PostgreSQL host | ✅ Yes (auto-set by Railway) |
| `POSTGRES_PORT` | PostgreSQL port (default: 5432) | ✅ Yes (auto-set by Railway) |

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
- `/receive-updates` - Subscribe to updates with a click of a button
- Support panel system with ticket management
- Event logging
- Auto-migration database system
- And more...

### Tech Stack

- **Language**: Go 1.25
- **Discord Library**: discordgo v0.29.0
- **ORM**: GORM v1.25 with Auto-Migration
- **Database**: PostgreSQL 18
- **Container**: Docker (multi-stage build)
- **Deployment**: Railway
- **Features**: Auto-migration, Ambassador System, Freemium Education System

### Support

For issues or questions, please open an issue on GitHub.

---

Built with 🦉 by Owls Capital

