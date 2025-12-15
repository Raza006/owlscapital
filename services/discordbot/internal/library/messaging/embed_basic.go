package messaging

// Exemplo: Envia um embed simples com estilo padronizado.
// Intents necessárias: Guilds (para slash command).
// Teste: registre o comando `/library-embed-basic`, execute e veja o embed formatado.
// Inventário: cobre "Mensagens & Embeds" (docs/library/inventory.md).

import (
	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const embedBasicCommandName = "library-embed-basic"

// LibraryEmbedBasicFeature demonstra como construir e responder com um embed.
type LibraryEmbedBasicFeature struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryEmbedBasicFeature cria a feature exemplo.
func NewLibraryEmbedBasicFeature() bot.Feature {
	return &LibraryEmbedBasicFeature{}
}

// RegisterLibraryEmbedBasicFeature permite registrar a feature manualmente.
func RegisterLibraryEmbedBasicFeature() {
	bot.RegisterFeature(NewLibraryEmbedBasicFeature())
}

// Name identifica a feature no registrador.
func (f *LibraryEmbedBasicFeature) Name() string { return "LibraryEmbedBasic" }

func (f *LibraryEmbedBasicFeature) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			embedBasicCommandName,
			"Exibe um embed básico com título, descrição e rodapé.",
			bot.ScopeGuild,
			f.handleEmbedBasic,
			true,
		)
	}
	return f.binding
}

// CommandSpecs expõe o comando slash.
func (f *LibraryEmbedBasicFeature) CommandSpecs() []bot.CommandSpec {
	binding := f.getBinding()
	binding.Spec.Command.DescriptionLocalizations = map[string]string{
		discordgo.BrazilianPortuguese: "Mostra um embed de exemplo com estilo starter kit.",
	}
	return []bot.CommandSpec{binding.Spec}
}

// ApplicationCommandHandlers associa o comando ao handler.
func (f *LibraryEmbedBasicFeature) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

// Intents retorna as intents necessárias.
func (f *LibraryEmbedBasicFeature) Intents() discordgo.Intent {
	return discordgo.IntentGuilds
}

func (f *LibraryEmbedBasicFeature) handleEmbedBasic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Starter Kit Embed",
		Description: "Este embed demonstra **título**, _descrição_ e campos padronizados.",
		Color:       0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Como testar",
				Value:  "Execute o comando e observe o layout. Personalize para sua guild.",
				Inline: false,
			},
			{
				Name:   "Referência",
				Value:  "docs/library/inventory.md → Mensagens & Embeds",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Starterkit • Exemplo de embed",
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
