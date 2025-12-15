package events

// Exemplo: acompanha eventos de guilda, membros e cargos.
// Intents necessárias: Guilds, GuildMembers.
// Inventário: cobre "Eventos de guild".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryGuildEvents registra vários handlers de guilda.
type LibraryGuildEvents struct {
	bot.BaseFeature
}

// NewLibraryGuildEvents cria a feature.
func NewLibraryGuildEvents() bot.Feature { return &LibraryGuildEvents{} }

// RegisterLibraryGuildEvents registra manualmente.
func RegisterLibraryGuildEvents() { bot.RegisterFeature(NewLibraryGuildEvents()) }

func (f *LibraryGuildEvents) Name() string { return "LibraryGuildEvents" }

func (f *LibraryGuildEvents) Intents() discordgo.Intent {
	return discordgo.IntentGuilds | discordgo.IntentGuildMembers
}

func (f *LibraryGuildEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.GuildCreate{}): {
			func(s *discordgo.Session, evt interface{}) {
				gc := evt.(*discordgo.GuildCreate)
				log.Printf("[Library][Guild] Guild adicionada: %s (%s)", gc.Guild.Name, gc.Guild.ID)
			},
		},
		reflect.TypeOf(discordgo.GuildUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				gu := evt.(*discordgo.GuildUpdate)
				log.Printf("[Library][Guild] Guild atualizada: %s", gu.Guild.Name)
			},
		},
		reflect.TypeOf(discordgo.GuildDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				gd := evt.(*discordgo.GuildDelete)
				log.Printf("[Library][Guild] Guild removida: %s", gd.Guild.ID)
			},
		},
		reflect.TypeOf(discordgo.GuildMemberAdd{}): {
			func(_ *discordgo.Session, evt interface{}) {
				add := evt.(*discordgo.GuildMemberAdd)
				log.Printf("[Library][Member] %s entrou em %s", add.User.Username, add.GuildID)
			},
		},
		reflect.TypeOf(discordgo.GuildMemberRemove{}): {
			func(_ *discordgo.Session, evt interface{}) {
				rm := evt.(*discordgo.GuildMemberRemove)
				log.Printf("[Library][Member] %s saiu de %s", rm.User.Username, rm.GuildID)
			},
		},
		reflect.TypeOf(discordgo.GuildRoleCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				role := evt.(*discordgo.GuildRoleCreate)
				log.Printf("[Library][Role] cargo criado: %s", role.Role.Name)
			},
		},
		reflect.TypeOf(discordgo.GuildRoleUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				role := evt.(*discordgo.GuildRoleUpdate)
				log.Printf("[Library][Role] cargo atualizado: %s", role.Role.Name)
			},
		},
		reflect.TypeOf(discordgo.GuildRoleDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				role := evt.(*discordgo.GuildRoleDelete)
				log.Printf("[Library][Role] cargo removido: %s", role.RoleID)
			},
		},
	}
}
