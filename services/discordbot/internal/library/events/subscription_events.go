package events

// Exemplo: acompanha eventos de inscrições (scheduled events, subscriptions).
// Intents necessárias: Guilds.
// Inventário: cobre "Scheduled events / SubscriptionUpdate".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibrarySubscriptionEvents registra handlers de eventos especiais.
type LibrarySubscriptionEvents struct {
	bot.BaseFeature
}

// NewLibrarySubscriptionEvents cria a feature.
func NewLibrarySubscriptionEvents() bot.Feature { return &LibrarySubscriptionEvents{} }

// RegisterLibrarySubscriptionEvents registra manualmente.
func RegisterLibrarySubscriptionEvents() { bot.RegisterFeature(NewLibrarySubscriptionEvents()) }

func (f *LibrarySubscriptionEvents) Name() string { return "LibrarySubscriptionEvents" }

func (f *LibrarySubscriptionEvents) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibrarySubscriptionEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.GuildScheduledEventCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				create := evt.(*discordgo.GuildScheduledEventCreate)
				log.Printf("[Library][ScheduledEvent] Criado evento %s", create.GuildScheduledEvent.Name)
			},
		},
		reflect.TypeOf(discordgo.GuildScheduledEventUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				update := evt.(*discordgo.GuildScheduledEventUpdate)
				log.Printf("[Library][ScheduledEvent] Atualizado evento %s", update.GuildScheduledEvent.Name)
			},
		},
		reflect.TypeOf(discordgo.GuildScheduledEventDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				deleteEvt := evt.(*discordgo.GuildScheduledEventDelete)
				log.Printf("[Library][ScheduledEvent] Removido evento %s", deleteEvt.GuildScheduledEvent.Name)
			},
		},
		reflect.TypeOf(discordgo.SubscriptionUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				sub := evt.(*discordgo.SubscriptionUpdate)
				log.Printf("[Library][Subscription] Atualização de subscription %s estado %s", sub.Subscription.ID, sub.Subscription.Status)
			},
		},
	}
}
