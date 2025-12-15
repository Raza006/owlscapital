package events

// Exemplo: acompanha eventos de threads (criação, atualização, membros).
// Intents necessárias: Guilds, GuildMessages, GuildMembers.
// Inventário: cobre "Eventos de thread".

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryThreadEvents registra handlers de threads.
type LibraryThreadEvents struct {
	bot.BaseFeature
}

// NewLibraryThreadEvents cria a feature.
func NewLibraryThreadEvents() bot.Feature { return &LibraryThreadEvents{} }

// RegisterLibraryThreadEvents registra manualmente.
func RegisterLibraryThreadEvents() { bot.RegisterFeature(NewLibraryThreadEvents()) }

func (f *LibraryThreadEvents) Name() string { return "LibraryThreadEvents" }

func (f *LibraryThreadEvents) Intents() discordgo.Intent {
	return discordgo.IntentGuilds | discordgo.IntentGuildMessages | discordgo.IntentGuildMembers
}

func (f *LibraryThreadEvents) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf(discordgo.ThreadCreate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				thread := evt.(*discordgo.ThreadCreate)
				log.Printf("[Library][Thread] Criada thread %s em %s", thread.Thread.Name, thread.Thread.ParentID)
			},
		},
		reflect.TypeOf(discordgo.ThreadDelete{}): {
			func(_ *discordgo.Session, evt interface{}) {
				thread := evt.(*discordgo.ThreadDelete)
				log.Printf("[Library][Thread] Removida thread %s", thread.ID)
			},
		},
		reflect.TypeOf(discordgo.ThreadMemberUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				update := evt.(*discordgo.ThreadMemberUpdate)
				log.Printf("[Library][Thread] Membro %s alterou atividade na thread %s", update.ThreadMember.UserID, update.ID)
			},
		},
		reflect.TypeOf(discordgo.ThreadMembersUpdate{}): {
			func(_ *discordgo.Session, evt interface{}) {
				update := evt.(*discordgo.ThreadMembersUpdate)
				log.Printf("[Library][Thread] Thread %s agora tem %d membros", update.ID, update.MemberCount)
			},
		},
	}
}
