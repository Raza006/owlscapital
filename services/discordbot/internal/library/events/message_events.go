package events

// Exemplo: acompanha eventos de atualização e deleção de mensagens.
// Intents necessárias: GuildMessages.
// Inventário: cobre "Eventos de mensagem (update/delete)".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryMessageEvents registra handlers de atualização/remoção de mensagens.
type LibraryMessageEvents struct {
	bot.BaseFeature
}

// NewLibraryMessageEvents cria a feature.
func NewLibraryMessageEvents() bot.Feature { return &LibraryMessageEvents{} }

// RegisterLibraryMessageEvents registra manualmente.
func RegisterLibraryMessageEvents() { bot.RegisterFeature(NewLibraryMessageEvents()) }

func (f *LibraryMessageEvents) Name() string { return "LibraryMessageEvents" }

func (f *LibraryMessageEvents) Intents() discordgo.Intent { return discordgo.IntentGuildMessages }

func (f *LibraryMessageEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.MessageUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				update := evt.(*discordgo.MessageUpdate)
				log.Printf("[Library][Message] Mensagem %s atualizada. Conteúdo agora: %s", update.ID, update.Content)
			},
		},
		reflect.TypeOf(discordgo.MessageDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				del := evt.(*discordgo.MessageDelete)
				log.Printf("[Library][Message] Mensagem %s removida do canal %s", del.ID, del.ChannelID)
			},
		},
		reflect.TypeOf(discordgo.MessageDeleteBulk{}): {
			func(_ *discordgo.Session, evt interface{}) {
				bulk := evt.(*discordgo.MessageDeleteBulk)
				log.Printf("[Library][Message] Remoção em massa (%d mensagens) no canal %s", len(bulk.Messages), bulk.ChannelID)
			},
		},
	}
}
