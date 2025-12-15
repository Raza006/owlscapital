package interactions

// Exemplo: envia um botão simples e responde ao clique.
// Intents necessárias: Guilds.
// Teste: execute `/library-button-basic`, clique no botão e observe a resposta.
// Inventário: cobre "Componentes (botões)" (docs/library/inventory.md).

import (
	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const (
	buttonBasicCommandName = "library-button-basic"
	buttonBasicCustomID    = "LibraryButtonBasic:primary"
)

// LibraryButtonBasic demonstra um botão Discord simples.
type LibraryButtonBasic struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryButtonBasic cria a feature de exemplo.
func NewLibraryButtonBasic() bot.Feature {
	return &LibraryButtonBasic{}
}

// RegisterLibraryButtonBasic registra manualmente a feature de exemplo.
func RegisterLibraryButtonBasic() {
	bot.RegisterFeature(NewLibraryButtonBasic())
}

// Name identifica a feature.
func (f *LibraryButtonBasic) Name() string { return "LibraryButtonBasic" }

func (f *LibraryButtonBasic) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			buttonBasicCommandName,
			"Gera uma mensagem com um botão clicável.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

// CommandSpecs expõe o slash command.
func (f *LibraryButtonBasic) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

// ApplicationCommandHandlers associa o comando ao handler.
func (f *LibraryButtonBasic) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

// ComponentHandlers registra o handler do botão.
func (f *LibraryButtonBasic) ComponentHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		buttonBasicCustomID: f.handleButton,
	}
}

// Intents necessárias para slash commands.
func (f *LibraryButtonBasic) Intents() discordgo.Intent {
	return discordgo.IntentGuilds
}

func (f *LibraryButtonBasic) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := &discordgo.InteractionResponseData{
		Content: "Clique no botão para ver a resposta.",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.PrimaryButton,
						CustomID: buttonBasicCustomID,
						Label:    "Pressione aqui",
					},
				},
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: message,
	})
}

func (f *LibraryButtonBasic) handleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Botão pressionado com sucesso!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
