package events

// Exemplo: registra eventos de canais (create/update/delete) e pins.
// Intents necessárias: Guilds.
// Inventário: cobre "Eventos de canal".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryChannelEvents cobre eventos básicos de canal.
type LibraryChannelEvents struct {
	bot.BaseFeature
}

// NewLibraryChannelEvents cria a feature.
func NewLibraryChannelEvents() bot.Feature { return &LibraryChannelEvents{} }

// RegisterLibraryChannelEvents registra manualmente.
func RegisterLibraryChannelEvents() { bot.RegisterFeature(NewLibraryChannelEvents()) }

func (f *LibraryChannelEvents) Name() string { return "LibraryChannelEvents" }

func (f *LibraryChannelEvents) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryChannelEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.ChannelCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				ch := evt.(*discordgo.ChannelCreate)
				log.Printf("[Library][Channel] Canal criado: %s (%s)", ch.Channel.Name, ch.Channel.ID)
			},
		},
		reflect.TypeOf(discordgo.ChannelUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				ch := evt.(*discordgo.ChannelUpdate)
				log.Printf("[Library][Channel] Canal atualizado: %s", ch.Channel.Name)
			},
		},
		reflect.TypeOf(discordgo.ChannelDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				ch := evt.(*discordgo.ChannelDelete)
				log.Printf("[Library][Channel] Canal removido: %s", ch.Channel.ID)
			},
		},
		reflect.TypeOf(discordgo.ChannelPinsUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				pins := evt.(*discordgo.ChannelPinsUpdate)
				log.Printf("[Library][Channel] Pins atualizados em %s", pins.ChannelID)
			},
		},
	}
}
