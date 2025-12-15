package events

// Exemplo: monitora o ciclo de vida do gateway (Ready, Resumed, Disconnect).
// Intents necessárias: nenhuma específica além das usadas por outras features.
// Inventário: cobre "Gateway reconnect/ready/resumed heartbeat".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryGatewayLifecycle registra eventos de ciclo de vida.
type LibraryGatewayLifecycle struct {
	bot.BaseFeature
}

// NewLibraryGatewayLifecycle cria a feature.
func NewLibraryGatewayLifecycle() bot.Feature { return &LibraryGatewayLifecycle{} }

// RegisterLibraryGatewayLifecycle registra manualmente.
func RegisterLibraryGatewayLifecycle() { bot.RegisterFeature(NewLibraryGatewayLifecycle()) }

func (f *LibraryGatewayLifecycle) Name() string { return "LibraryGatewayLifecycle" }

func (f *LibraryGatewayLifecycle) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.Ready{}): {
			func(_ *discordgo.Session, evt interface{}) {
				ready := evt.(*discordgo.Ready)
				log.Printf("[Library][Gateway] Ready recebido. Usuario: %s", ready.User.Username)
			},
		},
		reflect.TypeOf(discordgo.Resumed{}): {
			func(_ *discordgo.Session, evt interface{}) {
				log.Println("[Library][Gateway] Conexão resumida com sucesso")
			},
		},
		reflect.TypeOf(discordgo.Disconnect{}): {
			func(_ *discordgo.Session, evt interface{}) {
				disc := evt.(*discordgo.Disconnect)
				log.Printf("[Library][Gateway] Desconectado. Código: %d, Clean: %t", disc.Code, disc.Clean)
			},
		},
	}
}
