package oauth

// Exemplo: gera uma URL de autorização OAuth2 com escopos básicos.
// Intents: não aplicável.
// Teste: execute `/library-oauth-link` informando client_id e redirect, copie a URL.
// Inventário: cobre "OAuth2 flows".

import (
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const oauthLinkCommandName = "library-oauth-link"

// LibraryOAuthLink demonstra construção de URL OAuth.
type LibraryOAuthLink struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryOAuthLink cria a feature.
func NewLibraryOAuthLink() bot.Feature { return &LibraryOAuthLink{} }

// RegisterLibraryOAuthLink registra manualmente.
func RegisterLibraryOAuthLink() { bot.RegisterFeature(NewLibraryOAuthLink()) }

func (f *LibraryOAuthLink) Name() string { return "LibraryOAuthLink" }

func (f *LibraryOAuthLink) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			oauthLinkCommandName,
			"Gera uma URL OAuth2 com escopos específicos.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		binding.Spec.Command.Options = []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "client_id",
				Description: "ID da aplicação",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "redirect",
				Description: "Redirect URI cadastrado",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "scopes",
				Description: "Escopos separados por espaço (ex.: identify guilds)",
				Required:    true,
			},
		}
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryOAuthLink) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryOAuthLink) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryOAuthLink) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	clientID := data.Options[0].StringValue()
	redirect := url.QueryEscape(data.Options[1].StringValue())
	scopes := url.QueryEscape(strings.Join(strings.Fields(data.Options[2].StringValue()), " "))

	builder := strings.Builder{}
	builder.WriteString(discordgo.EndpointOAuth2Authorize)
	builder.WriteString("?client_id=")
	builder.WriteString(clientID)
	builder.WriteString("&response_type=code&scope=")
	builder.WriteString(scopes)
	builder.WriteString("&redirect_uri=")
	builder.WriteString(redirect)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "URL gerada:\n" + builder.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
