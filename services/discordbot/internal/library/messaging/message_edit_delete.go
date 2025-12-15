package messaging

// Exemplo: demonstra como editar e deletar a resposta inicial de um slash command.
// Intents necessárias: Guilds.
// Teste: registre `/library-message-edit`, execute e observe o conteúdo ser atualizado após 2s, depois removido em seguida.
// Inventário: cobre "Mensagens avançadas (edição/deleção)".

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const messageEditCommandName = "library-message-edit"

// LibraryMessageEdit demonstra edição e deleção da resposta de interação.
type LibraryMessageEdit struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryMessageEdit cria a feature.
func NewLibraryMessageEdit() bot.Feature {
	return &LibraryMessageEdit{}
}

// RegisterLibraryMessageEdit registra manualmente.
func RegisterLibraryMessageEdit() {
	bot.RegisterFeature(NewLibraryMessageEdit())
}

func (f *LibraryMessageEdit) Name() string { return "LibraryMessageEdit" }

func (f *LibraryMessageEdit) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			messageEditCommandName,
			"Mostra como editar e deletar a resposta de uma interação.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryMessageEdit) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryMessageEdit) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryMessageEdit) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryMessageEdit) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Responde em dois passos para conseguir editar a mensagem original.
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	time.AfterFunc(2*time.Second, func() {
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPtr("Mensagem atualizada após 2 segundos!"),
		})

		time.AfterFunc(2*time.Second, func() {
			_ = s.InteractionResponseDelete(i.Interaction)
		})
	})
}

func stringPtr(v string) *string { return &v }
