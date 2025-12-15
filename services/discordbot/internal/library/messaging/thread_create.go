package messaging

// Exemplo: cria uma thread pública em um canal de texto e envia uma mensagem.
// Intents necessárias: Guilds + GuildMessages (para enviar na thread).
// Teste: execute `/library-thread-create` em um canal de texto. Veja a thread "Starterkit Thread" criada.
// Inventário: cobre "Threads".

import (
	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const threadCreateCommandName = "library-thread-create"

// LibraryThreadCreate demonstra criação de threads.
type LibraryThreadCreate struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryThreadCreate cria a feature.
func NewLibraryThreadCreate() bot.Feature {
	return &LibraryThreadCreate{}
}

// RegisterLibraryThreadCreate registra manualmente.
func RegisterLibraryThreadCreate() {
	bot.RegisterFeature(NewLibraryThreadCreate())
}

func (f *LibraryThreadCreate) Name() string { return "LibraryThreadCreate" }

func (f *LibraryThreadCreate) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			threadCreateCommandName,
			"Cria uma thread pública e envia uma mensagem de boas-vindas.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryThreadCreate) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryThreadCreate) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryThreadCreate) Intents() discordgo.Intent {
	return discordgo.IntentGuilds | discordgo.IntentGuildMessages
}

func (f *LibraryThreadCreate) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	name := "Starterkit Thread"

	thread, err := s.StartThread(channelID, name, 60, discordgo.ChannelTypeGuildPublicThread)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Falha ao criar thread: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	_, _ = s.ChannelMessageSend(thread.ID, "Thread criada pelo starterkit! Consulte docs/library/inventory.md → Threads")

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thread criada em: <#" + thread.ID + ">",
		},
	})
}
