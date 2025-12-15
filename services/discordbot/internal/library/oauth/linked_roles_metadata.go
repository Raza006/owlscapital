package oauth

// Exemplo: registra metadados de linked roles via API.
// Requer token de aplicativo com escopo `role_connections.write`.
// Inventário: cobre "Linked roles".

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

// LibraryLinkedRolesMetadata demonstra chamada para atualizar metadados.
type LibraryLinkedRolesMetadata struct {
	bot.BaseFeature
}

// NewLibraryLinkedRolesMetadata cria a feature.
func NewLibraryLinkedRolesMetadata() bot.Feature { return &LibraryLinkedRolesMetadata{} }

// RegisterLibraryLinkedRolesMetadata registra manualmente.
func RegisterLibraryLinkedRolesMetadata() { bot.RegisterFeature(NewLibraryLinkedRolesMetadata()) }

func (f *LibraryLinkedRolesMetadata) Name() string { return "LibraryLinkedRolesMetadata" }

// Intents não são necessárias.

func (f *LibraryLinkedRolesMetadata) CommandSpecs() []bot.CommandSpec { return nil }

func (f *LibraryLinkedRolesMetadata) ApplicationCommandHandlers() bot.InteractionHandlersMap { return nil }

// UpdateMetadata demonstra como atualizar os metadados usando o client do bot.
func (f *LibraryLinkedRolesMetadata) UpdateMetadata(s *discordgo.Session) {
	metadata := []*discordgo.ApplicationRoleConnectionMetadata{
		{
			Type:        discordgo.ApplicationRoleConnectionMetadataTypeIntegerLessThanOrEqual,
			Key:         "dungeons",
			Name:        "Número de dungeons",
			Description: "Concede o cargo se o número for <= ao configurado",
		},
	}

	if err := s.ApplicationRoleConnectionMetadataUpdate(s.State.User.ID, metadata); err != nil {
		log.Printf("[Library][LinkedRoles] Falha ao atualizar metadados: %v", err)
		return
	}
	log.Println("[Library][LinkedRoles] Metadados atualizados com sucesso.")
}
