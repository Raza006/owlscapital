package events

// Exemplo: monitora eventos de voz (state e server update).
// Intents necessárias: GuildVoiceStates.
// Inventário: cobre "Eventos de voz".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryVoiceEvents registra handlers de voz.
type LibraryVoiceEvents struct {
	bot.BaseFeature
}

// NewLibraryVoiceEvents cria a feature.
func NewLibraryVoiceEvents() bot.Feature { return &LibraryVoiceEvents{} }

// RegisterLibraryVoiceEvents registra manualmente.
func RegisterLibraryVoiceEvents() { bot.RegisterFeature(NewLibraryVoiceEvents()) }

func (f *LibraryVoiceEvents) Name() string { return "LibraryVoiceEvents" }

func (f *LibraryVoiceEvents) Intents() discordgo.Intent {
	return discordgo.IntentGuildVoiceStates
}

func (f *LibraryVoiceEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.VoiceStateUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				state := evt.(*discordgo.VoiceStateUpdate)
				log.Printf("[Library][Voice] %s mudou para canal %s", state.UserID, state.ChannelID)
			},
		},
		reflect.TypeOf(discordgo.VoiceServerUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				update := evt.(*discordgo.VoiceServerUpdate)
				log.Printf("[Library][Voice] Voice server update em guild %s", update.GuildID)
			},
		},
	}
}
