package webhooks

// Exemplo: executa um webhook com ID/token fornecidos.
// Teste: crie um webhook no canal, cole ID e token em `/library-webhook-exec`.
// Inventário: cobre "Webhook execution e gerenciamento".

import (
	"strconv"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const webhookExecCommandName = "library-webhook-exec"

// LibraryWebhookExecute demonstra envio de mensagem via webhook.
type LibraryWebhookExecute struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryWebhookExecute cria a feature.
func NewLibraryWebhookExecute() bot.Feature { return &LibraryWebhookExecute{} }

// RegisterLibraryWebhookExecute registra manualmente.
func RegisterLibraryWebhookExecute() { bot.RegisterFeature(NewLibraryWebhookExecute()) }

func (f *LibraryWebhookExecute) Name() string { return "LibraryWebhookExecute" }

func (f *LibraryWebhookExecute) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			webhookExecCommandName,
			"Executa um webhook com conteúdo simples.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		binding.Spec.Command.Options = []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "id", Description: "ID do webhook", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "token", Description: "Token do webhook", Required: true},
		}
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryWebhookExecute) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryWebhookExecute) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryWebhookExecute) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	id := data.Options[0].StringValue()
	token := data.Options[1].StringValue()

	wid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "ID inválido", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	_, err = s.WebhookExecute(wid, token, true, &discordgo.WebhookParams{Content: "Mensagem via webhook executada pelo starterkit."})
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao executar webhook: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Webhook executado!"},
	})
}
