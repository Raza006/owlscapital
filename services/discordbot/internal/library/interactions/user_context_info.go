package interactions

// Exemplo: comando de contexto (usuário) que mostra informações básicas.
// Intents necessárias: Guilds.
// Teste: após registrar, clique com o botão direito em um usuário e selecione "Library: Ver Perfil".
// Inventário: cobre "Context menu (User)".

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryUserContextInfo demonstra uso de context menu para usuários.
type LibraryUserContextInfo struct {
	bot.BaseFeature
}

// NewLibraryUserContextInfo cria a feature.
func NewLibraryUserContextInfo() bot.Feature {
	return &LibraryUserContextInfo{}
}

// RegisterLibraryUserContextInfo registra manualmente.
func RegisterLibraryUserContextInfo() {
	bot.RegisterFeature(NewLibraryUserContextInfo())
}

func (f *LibraryUserContextInfo) Name() string { return "LibraryUserContextInfo" }

func (f *LibraryUserContextInfo) CommandSpecs() []bot.CommandSpec {
	cmd := &discordgo.ApplicationCommand{
		Name: "Library: Ver Perfil",
		Type: discordgo.UserApplicationCommand,
	}
	return []bot.CommandSpec{{Command: cmd, Scope: bot.ScopeGuild}}
}

func (f *LibraryUserContextInfo) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		"Library: Ver Perfil": f.handleUserContext,
	}
}

func (f *LibraryUserContextInfo) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryUserContextInfo) handleUserContext(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Type != discordgo.UserApplicationCommandDataType {
		return
	}

	target := data.TargetUser()
	content := fmt.Sprintf("Usuário selecionado: %s#%s (ID %s)", target.Username, target.Discriminator, target.ID)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content, Flags: discordgo.MessageFlagsEphemeral},
	})
}
