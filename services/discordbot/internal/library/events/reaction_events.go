package events

// Exemplo: observa eventos de reação em mensagens.
// Intents necessárias: GuildMessageReactions.
// Inventário: cobre "Eventos de reação".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryReactionEvents registra handlers para reações.
type LibraryReactionEvents struct {
	bot.BaseFeature
}

// NewLibraryReactionEvents cria a feature.
func NewLibraryReactionEvents() bot.Feature { return &LibraryReactionEvents{} }

// RegisterLibraryReactionEvents registra manualmente.
func RegisterLibraryReactionEvents() { bot.RegisterFeature(NewLibraryReactionEvents()) }

func (f *LibraryReactionEvents) Name() string { return "LibraryReactionEvents" }

func (f *LibraryReactionEvents) Intents() discordgo.Intent {
	return discordgo.IntentGuildMessageReactions
}

func (f *LibraryReactionEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.MessageReactionAdd{}): {
			func(_ *discordgo.Session, evt interface{}) {
				add := evt.(*discordgo.MessageReactionAdd)
				log.Printf("[Library][Reaction] %s adicionou %s na mensagem %s", add.UserID, add.Emoji.Name, add.MessageID)
			},
		},
		reflect.TypeOf(discordgo.MessageReactionRemove{}): {
			func(_ *discordgo.Session, evt interface{}) {
				rm := evt.(*discordgo.MessageReactionRemove)
				log.Printf("[Library][Reaction] %s removeu %s da mensagem %s", rm.UserID, rm.Emoji.Name, rm.MessageID)
			},
		},
		reflect.TypeOf(discordgo.MessageReactionRemoveAll{}): {
			func(_ *discordgo.Session, evt interface{}) {
				rm := evt.(*discordgo.MessageReactionRemoveAll)
				log.Printf("[Library][Reaction] Todas reações removidas da mensagem %s", rm.MessageID)
			},
		},
	}
}
