package rest

// Exemplo: faz upsert de um comando simples via REST.
// Inventário: cobre "Aplicação (commands, guild commands)".

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restCommandSyncName = "library-rest-sync-command"

// LibraryApplicationCommandSync demonstra ApplicationCommandCreate.
type LibraryApplicationCommandSync struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryApplicationCommandSync cria a feature.
func NewLibraryApplicationCommandSync() bot.Feature { return &LibraryApplicationCommandSync{} }

// RegisterLibraryApplicationCommandSync registra manualmente.
func RegisterLibraryApplicationCommandSync() { bot.RegisterFeature(NewLibraryApplicationCommandSync()) }

func (f *LibraryApplicationCommandSync) Name() string { return "LibraryApplicationCommandSync" }

func (f *LibraryApplicationCommandSync) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			restCommandSyncName,
			"Cria ou atualiza um comando de teste via REST.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryApplicationCommandSync) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryApplicationCommandSync) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryApplicationCommandSync) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "starterkit-rest-demo",
		Description: "Comando criado via REST demonst.",
	}

	if _, err := s.ApplicationCommandCreate(s.State.User.ID, i.GuildID, cmd); err != nil {
		log.Printf("[Library][REST] Falha ao criar comando: %v", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao criar comando: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Comando starterkit-rest-demo sincronizado!", Flags: discordgo.MessageFlagsEphemeral},
	})
}
