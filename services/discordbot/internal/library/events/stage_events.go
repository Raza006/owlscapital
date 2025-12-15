package events

// Exemplo: observa eventos de Stage Instance.
// Intents necessárias: Guilds.
// Inventário: cobre "Stage instance".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryStageEvents registra eventos de palco.
type LibraryStageEvents struct {
	bot.BaseFeature
}

// NewLibraryStageEvents cria a feature.
func NewLibraryStageEvents() bot.Feature { return &LibraryStageEvents{} }

// RegisterLibraryStageEvents registra manualmente.
func RegisterLibraryStageEvents() { bot.RegisterFeature(NewLibraryStageEvents()) }

func (f *LibraryStageEvents) Name() string { return "LibraryStageEvents" }

func (f *LibraryStageEvents) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryStageEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.StageInstanceCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				stage := evt.(*discordgo.StageInstanceCreate)
				log.Printf("[Library][Stage] Stage criada: %s", stage.StageInstance.Topic)
			},
		},
		reflect.TypeOf(discordgo.StageInstanceUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				stage := evt.(*discordgo.StageInstanceUpdate)
				log.Printf("[Library][Stage] Stage atualizada: %s", stage.StageInstance.Topic)
			},
		},
		reflect.TypeOf(discordgo.StageInstanceDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				stage := evt.(*discordgo.StageInstanceDelete)
				log.Printf("[Library][Stage] Stage removida em guild %s", stage.StageInstance.GuildID)
			},
		},
	}
}
