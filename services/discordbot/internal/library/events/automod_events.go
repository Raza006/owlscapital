package events

// Exemplo: registra ações de Auto Moderation.
// Intents necessárias: Guilds.
// Inventário: cobre "Auto moderation".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryAutoModEvents observa execuções de regras de auto moderação.
type LibraryAutoModEvents struct {
	bot.BaseFeature
}

// NewLibraryAutoModEvents cria a feature.
func NewLibraryAutoModEvents() bot.Feature { return &LibraryAutoModEvents{} }

// RegisterLibraryAutoModEvents registra manualmente.
func RegisterLibraryAutoModEvents() { bot.RegisterFeature(NewLibraryAutoModEvents()) }

func (f *LibraryAutoModEvents) Name() string { return "LibraryAutoModEvents" }

func (f *LibraryAutoModEvents) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryAutoModEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.AutoModerationActionExecution{}): {
			func(_ *discordgo.Session, evt interface{}) {
				action := evt.(*discordgo.AutoModerationActionExecution)
				log.Printf("[Library][AutoMod] Regra %s executada. Ação: %s", action.RuleID, action.Action.Type)
			},
		},
	}
}
