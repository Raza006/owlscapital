package utility

import (
	"fmt"
	"log"
	"time"

	"projectdiscord/services/discordbot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

// PingFeature is the implementation for the /ping command.
type PingFeature struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewPingFeature creates a new instance of the PingFeature.
func NewPingFeature() bot.Feature {
	return &PingFeature{}
}

// init registers this feature automatically when the package is imported.
func init() {
	bot.RegisterFeature(NewPingFeature())
}

// Name returns the name of the feature.
func (p *PingFeature) Name() string {
	return "Ping"
}

// getBinding centraliza a definição do comando e o vínculo com o handler,
// evitando duplicação entre CommandSpec e ApplicationCommandHandlers.
// Usa lazy init para garantir uma única fonte de verdade.
func (p *PingFeature) getBinding() bot.CommandBinding {
    if p.binding.Spec.Command == nil {
        p.binding = bot.NewSlashCommandBinding(
            "ping",
            "Responds with Pong! and the API latency.",
            bot.ScopeGuild,
            p.handlePing,
            true,
        )
    }
    return p.binding
}

// CommandSpecs retorna a especificação do comando usando o helper, sem repetir o name.
func (p *PingFeature) CommandSpecs() []bot.CommandSpec {
	b := p.getBinding()
	return []bot.CommandSpec{b.Spec}
}

// ApplicationCommandHandlers returns the specific handlers for the commands of this feature.
func (p *PingFeature) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return p.getBinding().AppCommandHandlers
}

// handlePing is the handler for the /ping command.
func (p *PingFeature) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Round(time.Millisecond)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🏓 Pong! API latency is `%s`.", latency),
		},
	})

	if err != nil {
		log.Printf("Error responding to /ping command: %v", err)
	}
}


