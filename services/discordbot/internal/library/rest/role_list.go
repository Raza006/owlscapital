package rest

// Exemplo: lista cargos da guild e destaca cargos gerenciados.
// Inventário: cobre "Roles e permissions".

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restRoleCommandName = "library-rest-roles"

// LibraryRestRoleList demonstra uso de GuildRoles.
type LibraryRestRoleList struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestRoleList cria a feature.
func NewLibraryRestRoleList() bot.Feature { return &LibraryRestRoleList{} }

// RegisterLibraryRestRoleList registra manualmente.
func RegisterLibraryRestRoleList() { bot.RegisterFeature(NewLibraryRestRoleList()) }

func (f *LibraryRestRoleList) Name() string { return "LibraryRestRoleList" }

func (f *LibraryRestRoleList) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restRoleCommandName,
			"Mostra os cargos disponíveis.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestRoleList) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryRestRoleList) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestRoleList) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestRoleList) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao listar cargos: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	var names []string
	for _, role := range roles {
		label := role.Name
		if role.Managed {
			label += " (gerenciado)"
		}
		names = append(names, label)
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: strings.Join(names, ", ")},
	})
}
