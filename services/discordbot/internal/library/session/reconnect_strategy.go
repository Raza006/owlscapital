package session

// Exemplo: personaliza callbacks de reconnect e heartbeat.
// Inventário: complementa itens de gateway websocket.

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ConfigureReconnectStrategy ajusta handlers de reconnect/heartbeat.
func ConfigureReconnectStrategy(s *discordgo.Session) {
	s.ShouldReconnectOnError = func(err error) bool {
		log.Printf("[Library][Session] Avaliando reconnect após erro: %v", err)
		return true
	}

	s.LastHeartbeatAck = time.Now()
	s.Identify.LargeThreshold = 50
}
