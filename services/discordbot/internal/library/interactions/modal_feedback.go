package interactions

// Exemplo: abre um modal e trata a submissão.
// Intents necessárias: Guilds.
// Teste: execute `/library-modal-feedback`, preencha o formulário e veja a confirmação.
// Inventário: cobre "Modals".

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const (
	modalFeedbackCommandName = "library-modal-feedback"
	modalFeedbackCustomID    = "LibraryModalFeedback:submit"
)

// LibraryModalFeedback demonstra o fluxo completo de modals.
type LibraryModalFeedback struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryModalFeedback cria a feature.
func NewLibraryModalFeedback() bot.Feature {
	return &LibraryModalFeedback{}
}

// RegisterLibraryModalFeedback registra manualmente.
func RegisterLibraryModalFeedback() {
	bot.RegisterFeature(NewLibraryModalFeedback())
}

func (f *LibraryModalFeedback) Name() string { return "LibraryModalFeedback" }

func (f *LibraryModalFeedback) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			modalFeedbackCommandName,
			"Abre um modal de feedback simples.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryModalFeedback) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryModalFeedback) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryModalFeedback) ModalSubmitHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		modalFeedbackCustomID: f.handleModalSubmit,
	}
}

func (f *LibraryModalFeedback) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryModalFeedback) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: modalFeedbackCustomID,
			Title:    "Feedback Starterkit",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "feedback",
							Label:       "Conte o que achou",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Adoraria saber sua opinião!",
							Required:    true,
						},
					},
				},
			},
		},
	})
}

func (f *LibraryModalFeedback) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	feedback := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Recebido! Seu feedback foi: %q", feedback),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
