package rest

// Exemplo: lista eventos agendados da guild.
// Inventário: cobre "Scheduled events" na REST.

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restScheduledCommandName = "library-rest-scheduled"

// LibraryRestScheduledList demonstra uso de GuildScheduledEvents.
type LibraryRestScheduledList struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestScheduledList cria a feature.
func NewLibraryRestScheduledList() bot.Feature { return &LibraryRestScheduledList{} }

// RegisterLibraryRestScheduledList registra manualmente.
func RegisterLibraryRestScheduledList() { bot.RegisterFeature(NewLibraryRestScheduledList()) }

func (f *LibraryRestScheduledList) Name() string { return "LibraryRestScheduledList" }

func (f *LibraryRestScheduledList) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restScheduledCommandName,
			"Lista eventos agendados.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestScheduledList) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryRestScheduledList) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestScheduledList) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestScheduledList) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	events, err := s.GuildScheduledEvents(i.GuildID, true)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	if len(events) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Nenhum evento agendado."},
	})
		return
	}

	var lines []string
	for _, evt := range events {
		lines = append(lines, fmt.Sprintf("%s → %s", evt.Name, evt.Status.String()))
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: strings.Join(lines, "\n")},
	})
}
