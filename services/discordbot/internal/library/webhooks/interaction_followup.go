package webhooks

// Exemplo: envia mensagens follow-up usando o webhook da interação.
// Teste: execute `/library-followup-demo` e veja a resposta inicial + follow-up editável.
// Inventário: cobre "Interaction responses via webhook".

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const followupCommandName = "library-followup-demo"

// LibraryInteractionFollowup demonstra follow-ups e edição posterior.
type LibraryInteractionFollowup struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryInteractionFollowup cria a feature.
func NewLibraryInteractionFollowup() bot.Feature { return &LibraryInteractionFollowup{} }

// RegisterLibraryInteractionFollowup registra manualmente.
func RegisterLibraryInteractionFollowup() { bot.RegisterFeature(NewLibraryInteractionFollowup()) }

func (f *LibraryInteractionFollowup) Name() string { return "LibraryInteractionFollowup" }

func (f *LibraryInteractionFollowup) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			followupCommandName,
			"Demonstra follow-up e edição via webhook.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryInteractionFollowup) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryInteractionFollowup) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryInteractionFollowup) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryInteractionFollowup) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Resposta inicial enviada."},
	})

	go func() {
		time.Sleep(1 * time.Second)
		msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: "Follow-up via webhook do starterkit."})
		if err != nil {
			return
		}

		time.Sleep(2 * time.Second)
		_, _ = s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{Content: stringPtr("Follow-up editado depois de 2s.")})

		time.Sleep(2 * time.Second)
		_ = s.FollowupMessageDelete(i.Interaction, msg.ID)
	}()
}

func stringPtr(v string) *string { return &v }
