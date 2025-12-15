package interactions

// Exemplo: comando de contexto (mensagem) que cria uma citação formatada.
// Intents necessárias: Guilds.
// Teste: clique com o botão direito em uma mensagem → Apps → "Library: Citar".
// Inventário: cobre "Context menu (Message)".

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryMessageContextQuote mostra como ler dados da mensagem alvo.
type LibraryMessageContextQuote struct {
	bot.BaseFeature
}

// NewLibraryMessageContextQuote cria a feature.
func NewLibraryMessageContextQuote() bot.Feature {
	return &LibraryMessageContextQuote{}
}

// RegisterLibraryMessageContextQuote registra manualmente.
func RegisterLibraryMessageContextQuote() {
	bot.RegisterFeature(NewLibraryMessageContextQuote())
}

func (f *LibraryMessageContextQuote) Name() string { return "LibraryMessageContextQuote" }

func (f *LibraryMessageContextQuote) CommandSpecs() []bot.CommandSpec {
	cmd := &discordgo.ApplicationCommand{
		Name: "Library: Citar",
		Type: discordgo.MessageApplicationCommand,
	}
	return []bot.CommandSpec{{Command: cmd, Scope: bot.ScopeGuild}}
}

func (f *LibraryMessageContextQuote) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		"Library: Citar": f.handleMessageContext,
	}
}

func (f *LibraryMessageContextQuote) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryMessageContextQuote) handleMessageContext(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Type != discordgo.MessageApplicationCommandDataType {
		return
	}

	target := data.TargetMessage()
	content := fmt.Sprintf(
		">>> **%s** disse em <#%s>:\n%s",
		target.Author.Username,
		target.ChannelID,
		target.Content,
	)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}
