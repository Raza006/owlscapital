package interactions

// Exemplo: menu select multi escolha mostrando rows aninhados.
// Intents necessárias: Guilds.
// Teste: execute `/library-select-multi`, selecione opções e observe resposta.
// Inventário: cobre "Componentes avançados (multi-select, rows)".

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const (
	selectMultiCommandName = "library-select-multi"
	selectMultiCustomID   = "LibrarySelectMulti:roles"
)

// LibrarySelectMulti demonstra menus múltiplos.
type LibrarySelectMulti struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibrarySelectMulti cria a feature.
func NewLibrarySelectMulti() bot.Feature { return &LibrarySelectMulti{} }

// RegisterLibrarySelectMulti registra manualmente.
func RegisterLibrarySelectMulti() { bot.RegisterFeature(NewLibrarySelectMulti()) }

func (f *LibrarySelectMulti) Name() string { return "LibrarySelectMulti" }

func (f *LibrarySelectMulti) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			selectMultiCommandName,
			"Exibe um menu multi seleção.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibrarySelectMulti) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibrarySelectMulti) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibrarySelectMulti) ComponentHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		selectMultiCustomID:         f.handleSelect,
		"LibrarySelectMulti:cancel": f.handleCancel,
	}
}

func (f *LibrarySelectMulti) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibrarySelectMulti) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rows := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    selectMultiCustomID,
					MinValues:   1,
					MaxValues:   3,
					Placeholder: "Escolha até 3 funções",
					Options: []discordgo.SelectMenuOption{
						{Label: "Healer", Value: "healer", Description: "Cuida da equipe"},
						{Label: "Tank", Value: "tank", Description: "Protege aliados"},
						{Label: "Damage", Value: "damage", Description: "Causa dano"},
						{Label: "Crafter", Value: "crafter"},
						{Label: "Gatherer", Value: "gatherer"},
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Style: discordgo.SecondaryButton, CustomID: "LibrarySelectMulti:cancel", Label: "Cancelar"},
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Selecione as funções desejadas:",
			Components: rows,
		},
	})
}

func (f *LibrarySelectMulti) handleSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	values := i.MessageComponentData().Values
	content := "Você escolheu: " + strings.Join(values, ", ")

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (f *LibrarySelectMulti) handleCancel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Seleção cancelada.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
