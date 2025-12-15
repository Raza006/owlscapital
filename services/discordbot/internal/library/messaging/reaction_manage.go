package messaging

// Exemplo: adiciona e remove reações em uma mensagem enviada pelo bot.
// Intents necessárias: Guilds.
// Teste: execute `/library-reaction-manage`, veja os emojis adicionados e removidos automaticamente.
// Inventário: cobre "Reações" (add/remove/remove-all).

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const reactionCommandName = "library-reaction-manage"

// LibraryReactionManage demonstra APIs de reação.
type LibraryReactionManage struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryReactionManage cria a feature.
func NewLibraryReactionManage() bot.Feature {
	return &LibraryReactionManage{}
}

// RegisterLibraryReactionManage registra manualmente.
func RegisterLibraryReactionManage() {
	bot.RegisterFeature(NewLibraryReactionManage())
}

func (f *LibraryReactionManage) Name() string { return "LibraryReactionManage" }

func (f *LibraryReactionManage) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			reactionCommandName,
			"Demonstra como adicionar e remover reações.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryReactionManage) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryReactionManage) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryReactionManage) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryReactionManage) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := &discordgo.InteractionResponseData{Content: "Gerenciando reações neste exemplo."}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})

	go func() {
		time.Sleep(1 * time.Second)
		msg, err := s.InteractionResponse(i.Interaction)
		if err != nil {
			return
		}

		channelID := msg.ChannelID
		messageID := msg.ID

		emojis := []string{"👍", "🔥"}
		for _, emoji := range emojis {
			_ = s.MessageReactionAdd(channelID, messageID, emoji)
		}

		time.Sleep(2 * time.Second)
		_ = s.MessageReactionRemove(channelID, messageID, "👍", "@me")

		time.Sleep(2 * time.Second)
		_ = s.MessageReactionsRemoveAll(channelID, messageID)
	}()
}
