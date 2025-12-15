package rest

// Exemplo: consulta informações da guild via REST.
// Intents necessárias: Guilds.
// Teste: execute `/library-rest-guild` para ver contagem de membros/canais.
// Inventário: cobre "Guilds (create, modify, widgets, onboarding)" (consulta).

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restGuildCommandName = "library-rest-guild"

// LibraryRestGuildOverview demonstra consulta de guild via REST.
type LibraryRestGuildOverview struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestGuildOverview cria a feature.
func NewLibraryRestGuildOverview() bot.Feature { return &LibraryRestGuildOverview{} }

// RegisterLibraryRestGuildOverview registra manualmente.
func RegisterLibraryRestGuildOverview() { bot.RegisterFeature(NewLibraryRestGuildOverview()) }

func (f *LibraryRestGuildOverview) Name() string { return "LibraryRestGuildOverview" }

func (f *LibraryRestGuildOverview) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restGuildCommandName,
			"Mostra informações básicas da guild via REST.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestGuildOverview) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryRestGuildOverview) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestGuildOverview) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestGuildOverview) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao consultar guild: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	content := fmt.Sprintf("Guild: **%s**\nMembros: %d\nCanais: %d\nBoost level: %d",
		guild.Name,
		guild.ApproximateMemberCount,
		len(guild.Channels),
		guild.PremiumTier,
	)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}
