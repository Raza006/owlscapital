package state

// Exemplo: consulta informações do cache `Session.State`.
// Inventário: cobre "Estado interno e cache".

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// GuildSummary retorna nome e contagem de membros da guild a partir do cache.
func GuildSummary(s *discordgo.Session, guildID string) (string, int, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return "", 0, err
	}
	return guild.Name, guild.MemberCount, nil
}

// EnsureMemberCached demonstra como sincronizar membros caso não estejam no cache.
func EnsureMemberCached(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)
	if err == nil {
		return member, nil
	}

	// Solicita chunk manual para trazer o membro ao cache.
	err = s.RequestGuildMembers(guildID, userID, 1, false)
	if err != nil {
		return nil, err
	}

	member, err = s.State.Member(guildID, userID)
	if err != nil {
		return nil, errors.New("membro não encontrado após sync")
	}
	return member, nil
}
