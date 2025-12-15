package interactions

// Exemplo: slash command com autocomplete de cargos fictícios.
// Intents necessárias: Guilds.
// Teste: digite `/library-autocomplete-role` e comece a escrever, veja sugestões aparecerem.
// Inventário: cobre "Autocomplete".

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const autocompleteCommandName = "library-autocomplete-role"

var autocompleteRoles = []string{
	"Healer",
	"Tank",
	"Damage",
	"Crafter",
	"Gatherer",
}

// LibraryAutocompleteRole demonstra o fluxo de autocomplete.
type LibraryAutocompleteRole struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryAutocompleteRole cria a feature.
func NewLibraryAutocompleteRole() bot.Feature {
	return &LibraryAutocompleteRole{}
}

// RegisterLibraryAutocompleteRole registra manualmente.
func RegisterLibraryAutocompleteRole() {
	bot.RegisterFeature(NewLibraryAutocompleteRole())
}

func (f *LibraryAutocompleteRole) Name() string { return "LibraryAutocompleteRole" }

func (f *LibraryAutocompleteRole) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		cmd := bot.NewSlashCommandBinding(
			autocompleteCommandName,
			"Demonstra autocomplete em uma opção de texto.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		cmd.Spec.Command.Options = []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "role",
				Description: "Selecione uma função",
				Required:    true,
				Autocomplete: true,
			},
		}
		f.binding = cmd
	}
	return f.binding
}

func (f *LibraryAutocompleteRole) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryAutocompleteRole) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryAutocompleteRole) AutocompleteHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		autocompleteCommandName: f.handleAutocomplete,
	}
}

func (f *LibraryAutocompleteRole) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryAutocompleteRole) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	choice := i.ApplicationCommandData().Options[0].StringValue()
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Você escolheu: " + choice},
	})
}

func (f *LibraryAutocompleteRole) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	focused := i.ApplicationCommandData().Options[0].StringValue()

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, role := range autocompleteRoles {
		if strings.Contains(strings.ToLower(role), strings.ToLower(focused)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  role,
				Value: role,
			})
		}
	}
	if len(choices) == 0 {
		choices = []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Nenhuma opção",
				Value: "none",
			},
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{Choices: choices},
	})
}
