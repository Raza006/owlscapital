package rest

// Exemplo: lista canais de texto usando a REST API.
// Intents necessárias: Guilds.
// Inventário: cobre "Channels (text, voice, threads, stage)".

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restChannelCommandName = "library-rest-channels"

// LibraryRestChannelList demonstra a rota GuildChannels.
type LibraryRestChannelList struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestChannelList cria a feature.
func NewLibraryRestChannelList() bot.Feature { return &LibraryRestChannelList{} }

// RegisterLibraryRestChannelList registra manualmente.
func RegisterLibraryRestChannelList() { bot.RegisterFeature(NewLibraryRestChannelList()) }

func (f *LibraryRestChannelList) Name() string { return "LibraryRestChannelList" }

func (f *LibraryRestChannelList) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restChannelCommandName,
			"Lista os canais de texto da guild.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestChannelList) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryRestChannelList) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestChannelList) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestChannelList) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao listar canais: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	var names []string
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildText {
			names = append(names, "#"+ch.Name)
		}
	}
	if len(names) == 0 {
		names = []string{"Nenhum canal de texto encontrado"}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: strings.Join(names, ", ")},
	})
}
