package interactions

// Exemplo: comando com nome e descrição traduzidos usando locales.
// Intents necessárias: Guilds.
// Teste: altere o idioma do cliente Discord e veja os textos localizados.
// Inventário: cobre "localização / idiomas".

import (
	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const localizationCommandName = "library-localized"

// LibraryLocalization demonstra campos de localização.
type LibraryLocalization struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryLocalization cria a feature.
func NewLibraryLocalization() bot.Feature {
	return &LibraryLocalization{}
}

// RegisterLibraryLocalization registra manualmente.
func RegisterLibraryLocalization() {
	bot.RegisterFeature(NewLibraryLocalization())
}

func (f *LibraryLocalization) Name() string { return "LibraryLocalization" }

func (f *LibraryLocalization) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			localizationCommandName,
			"Shows localized text",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		cmd := binding.Spec.Command
		cmd.NameLocalizations = map[discordgo.Locale]string{
			discordgo.PortugueseBR: "biblioteca-localizada",
			discordgo.SpanishES:    "biblioteca-localizada",
		}
		cmd.DescriptionLocalizations = map[discordgo.Locale]string{
			discordgo.PortugueseBR: "Exibe texto no idioma do cliente.",
			discordgo.SpanishES:    "Muestra texto en el idioma del cliente.",
		}
		binding.Spec.Command = cmd
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryLocalization) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryLocalization) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryLocalization) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryLocalization) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	locale := i.Locale
	content := "Idioma atual: " + string(locale)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}
