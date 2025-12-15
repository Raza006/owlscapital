package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"projectdiscord/services/discordbot/internal/bot"

	"github.com/bwmarrin/discordgo"
)

// As features ativas são registradas automaticamente via init() em cada pacote de feature,
// e entram no binário por meio de imports com underscore no arquivo features_enabled.go
// (ex.: `import _ ".../internal/features/utility"`). Para adicionar uma nova feature,
// crie o pacote, chame bot.RegisterFeature() no init() da feature e adicione seu import com `_` em features_enabled.go.
var activeFeatures = bot.RegisteredFeatures()

func main() {
	log.Println("Starting Starter Kit bot...")

	// Load configuration from environment variables
	cfg, err := bot.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal error loading configuration: %v", err)
	}
	if cfg.GuildID == "" {
		log.Println("Warning: GUILD_ID is empty. Commands will be registered globally and may take up to 1 hour to propagate.")
	}

	// Load all active features and build the registry of handlers and commands
	registry := bot.LoadFeatures(activeFeatures)

	// Create a new Discord session using the bot token
	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Set the intents for the Discord session based on the features' needs
	dg.Identify.Intents = registry.Intents
	if dg.Identify.Intents == 0 {
		if cfg.EnableDefaultIntents {
			// Fallback opcional para facilitar o início rápido (guildas + DMs)
			dg.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentDirectMessages
		} else {
			log.Printf("Nenhuma feature definiu intents. Defina intents nas features ou habilite o fallback com DEFAULT_INTENTS=true.")
		}
	}

	// Add the universal dispatcher from the features package as the single event handler
	dg.AddHandler(bot.CreateDispatcher(registry))

	// Register all application commands after the Ready event to ensure state is initialized
	dg.AddHandlerOnce(func(s *discordgo.Session, _ *discordgo.Ready) {
		bot.RegisterCommands(s, cfg.GuildID, registry)
	})

	// Open a websocket connection to Discord and begin listening
	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}

	// Wait here until CTRL-C or other term signal is received
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	log.Println("Bot is shutting down.")
	dg.Close()
}


