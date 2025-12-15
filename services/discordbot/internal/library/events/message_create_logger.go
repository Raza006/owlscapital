package events

// Exemplo: loga mensagens recebidas via evento MessageCreate.
// Intents necessárias: GuildMessages.
// Teste: habilite `discordgo.IntentGuildMessages`, envie mensagem em um canal e observe o log.
// Inventário: cobre "Eventos de mensagem" (docs/library/inventory.md).

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryMessageCreateLogger demonstra como observar eventos de mensagem.
type LibraryMessageCreateLogger struct {
	bot.BaseFeature
}

// NewLibraryMessageCreateLogger cria a feature de exemplo.
func NewLibraryMessageCreateLogger() bot.Feature {
	return &LibraryMessageCreateLogger{}
}

// RegisterLibraryMessageCreateLogger registra manualmente a feature.
func RegisterLibraryMessageCreateLogger() {
	bot.RegisterFeature(NewLibraryMessageCreateLogger())
}

// Name identifica a feature.
func (f *LibraryMessageCreateLogger) Name() string { return "LibraryMessageCreateLogger" }

// Intents necessárias para receber o evento.
func (f *LibraryMessageCreateLogger) Intents() discordgo.Intent {
	return discordgo.IntentGuildMessages
}

// TypedEventHandlers associa o handler ao tipo concreto do evento.
func (f *LibraryMessageCreateLogger) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.MessageCreate{}): {
			func(s *discordgo.Session, evt interface{}) {
				mc, ok := evt.(*discordgo.MessageCreate)
				if !ok {
					return
				}
				log.Printf("[Library] %s#%s → canal %s: %s",
					mc.Author.Username,
					mc.Author.Discriminator,
					mc.ChannelID,
					mc.Content,
				)
			},
		},
	}
}
