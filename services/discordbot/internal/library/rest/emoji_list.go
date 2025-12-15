package rest

// Exemplo: lista emojis personalizados da guild.
// Inventário: cobre "Emojis, stickers".

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restEmojiCommandName = "library-rest-emojis"

// LibraryRestEmojiList demonstra uso de GuildEmojis.
type LibraryRestEmojiList struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestEmojiList cria a feature.
func NewLibraryRestEmojiList() bot.Feature { return &LibraryRestEmojiList{} }

// RegisterLibraryRestEmojiList registra manualmente.
func RegisterLibraryRestEmojiList() { bot.RegisterFeature(NewLibraryRestEmojiList()) }

func (f *LibraryRestEmojiList) Name() string { return "LibraryRestEmojiList" }

func (f *LibraryRestEmojiList) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restEmojiCommandName,
			"Lista emojis personalizados da guild.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryRestEmojiList) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryRestEmojiList) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestEmojiList) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryRestEmojiList) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	emojis, err := s.GuildEmojis(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao listar emojis: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	if len(emojis) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Sem emojis personalizados."},
		})
		return
	}

	var lines []string
	for _, emoji := range emojis {
		lines = append(lines, fmt.Sprintf("%s → :%s:", emoji.ID, emoji.Name))
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: strings.Join(lines, "\n")},
	})
}
