package rest

// Exemplo: lista regras de auto moderação via REST.
// Inventário: cobre "Auto moderation rules".

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restAutomodCommandName = "library-rest-automod"

// LibraryRestAutoModRules demonstra a rota GuildAutoModerationRules.
type LibraryRestAutoModRules struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestAutoModRules cria a feature.
func NewLibraryRestAutoModRules() bot.Feature { return &LibraryRestAutoModRules{} }

// RegisterLibraryRestAutoModRules registra manualmente.
func RegisterLibraryRestAutoModRules() { bot.RegisterFeature(NewLibraryRestAutoModRules()) }

func (f *LibraryRestAutoModRules) Name() string { return "LibraryRestAutoModRules" }

func (f *LibraryRestAutoModRules) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restAutomodCommandName,
			"Lista regras de auto moderação da guild.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestAutoModRules) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryRestAutoModRules) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestAutoModRules) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestAutoModRules) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rules, err := s.GuildAutoModerationRules(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	if len(rules) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Nenhuma regra cadastrada."},
	})
		return
	}

	var lines []string
	for _, rule := range rules {
		lines = append(lines, fmt.Sprintf("%s → Trigger: %s", rule.Name, rule.TriggerMetadata.KeywordFilter))
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: strings.Join(lines, "\n")},
	})
}
