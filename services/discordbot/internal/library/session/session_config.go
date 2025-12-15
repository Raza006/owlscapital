package session

// Exemplo: constrói uma sessão com intents, shards e logging configurados.
// Inventário: cobre "Construção de Session, intents, shards, logging".

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// NewStarterkitSession cria uma sessão pronta para produção.
func NewStarterkitSession(token string) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Define intents padrão do starterkit.
	s.Identify.Intents = discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildVoiceStates

	// Exemplo de configuração de shards (caso tenha vários).
	s.ShardID = 0
	s.ShardCount = 1

	// Ativa cache interno para permitir consultas via s.State.
	s.StateEnabled = true

	// Customiza logger para integrar com a stack do starterkit.
	s.LogLevel = discordgo.LogInformational
	s.Logger = func(_ int, caller, format string, a ...interface{}) {
		log.Printf("[discordgo][%s] "+format, append([]interface{}{caller}, a...)...)
	}

	// Ajusta limites de reconnect/backoff.
	s.MaxRestRetries = 3
	s.Ratelimiter = discordgo.NewBucket(5, 1*time.Second)

	return s, nil
}
