package rest

// Exemplo: consulta um membro específico e verifica se está banido.
// Intents necessárias: GuildMembers (para membros); requer permissão BanMembers para checar ban.
// Inventário: cobre "Members e bans".

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const restMemberCommandName = "library-rest-member"

// LibraryRestMemberLookup demonstra rotas GuildMember e GuildBan.
type LibraryRestMemberLookup struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryRestMemberLookup cria a feature.
func NewLibraryRestMemberLookup() bot.Feature { return &LibraryRestMemberLookup{} }

// RegisterLibraryRestMemberLookup registra manualmente.
func RegisterLibraryRestMemberLookup() { bot.RegisterFeature(NewLibraryRestMemberLookup()) }

func (f *LibraryRestMemberLookup) Name() string { return "LibraryRestMemberLookup" }

func (f *LibraryRestMemberLookup) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			restMemberCommandName,
			"Consulta um membro e verifica banimento.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		binding.Spec.Command.Options = []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user_id",
				Description: "ID do usuário a verificar",
				Required:    true,
			},
		}
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryRestMemberLookup) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryRestMemberLookup) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryRestMemberLookup) Intents() discordgo.Intent { return discordgo.IntentGuildMembers }

func (f *LibraryRestMemberLookup) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].StringValue()

	member, err := s.GuildMember(i.GuildID, userID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Erro ao buscar membro: " + err.Error(), Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	content := fmt.Sprintf("Membro: %s#%s\nCargo(s): %d", member.User.Username, member.User.Discriminator, len(member.Roles))

	if ban, err := s.GuildBan(i.GuildID, userID); err == nil {
		content += "\nBanido: sim (motivo: " + ban.Reason + ")"
	} else {
		content += "\nBanido: não"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content, Flags: discordgo.MessageFlagsEphemeral},
	})
}
