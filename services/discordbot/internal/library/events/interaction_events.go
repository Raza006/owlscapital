package events

// Exemplo: observa o evento InteractionCreate diretamente.
// Intents necessárias: Guilds.
// Inventário: cobre "Eventos de interação".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryInteractionEvents registra handler para InteractionCreate.
type LibraryInteractionEvents struct {
	bot.BaseFeature
}

// NewLibraryInteractionEvents cria a feature.
func NewLibraryInteractionEvents() bot.Feature { return &LibraryInteractionEvents{} }

// RegisterLibraryInteractionEvents registra manualmente.
func RegisterLibraryInteractionEvents() { bot.RegisterFeature(NewLibraryInteractionEvents()) }

func (f *LibraryInteractionEvents) Name() string { return "LibraryInteractionEvents" }

func (f *LibraryInteractionEvents) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryInteractionEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.InteractionCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				interaction := evt.(*discordgo.InteractionCreate)
				log.Printf("[Library][Interaction] Tipo: %d, ID: %s", interaction.Type, interaction.ID)
			},
		},
	}
}
