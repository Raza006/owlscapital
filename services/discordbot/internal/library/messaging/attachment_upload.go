package messaging

// Exemplo: envia um arquivo de texto como attachment na resposta.
// Intents necessárias: Guilds.
// Teste: execute `/library-attachment-upload` e baixe o arquivo retornado.
// Inventário: cobre "Attachments" em mensagens.

import (
	"bytes"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const attachmentCommandName = "library-attachment-upload"

// LibraryAttachmentUpload demonstra o envio de attachments.
type LibraryAttachmentUpload struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryAttachmentUpload cria a feature.
func NewLibraryAttachmentUpload() bot.Feature {
	return &LibraryAttachmentUpload{}
}

// RegisterLibraryAttachmentUpload registra manualmente a feature.
func RegisterLibraryAttachmentUpload() {
	bot.RegisterFeature(NewLibraryAttachmentUpload())
}

func (f *LibraryAttachmentUpload) Name() string { return "LibraryAttachmentUpload" }

func (f *LibraryAttachmentUpload) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			attachmentCommandName,
			"Mostra como anexar um arquivo simples na resposta.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryAttachmentUpload) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryAttachmentUpload) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryAttachmentUpload) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryAttachmentUpload) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := "Starterkit attachment de exemplo\nConsulte docs/library/inventory.md → Mensagens & Embeds"
	buf := bytes.NewBufferString(content)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Segue um arquivo de texto criado via discordgo.",
			Files: []*discordgo.File{
				{
					Name:   "starterkit.txt",
					Reader: buf,
				},
			},
		},
	})
}
