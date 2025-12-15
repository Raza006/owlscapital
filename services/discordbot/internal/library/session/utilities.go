package session

// Exemplo: utilidades diversas do pacote `discordgo` (conversão de IDs, timestamps).
// Inventário: cobre "Utilidades gerais (`util.go`)".

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// SnowflakeInfo retorna timestamp e worker extraídos de um ID do Discord.
func SnowflakeInfo(id string) (time.Time, int64) {
	return discordgo.SnowflakeTimestamp(id), discordgo.SnowflakeWorkerID(id)
}

// LinkifyChannel constrói a string de menção do canal.
func LinkifyChannel(channelID string) string {
	return discordgo.StrID(channelID)
}
