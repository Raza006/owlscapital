package voice

// Exemplo: conecta em um canal de voz e loga frames recebidos.
// Intents necessárias: GuildVoiceStates. Ative a flag "ReceiveEnabled" no bot.
// Teste: execute `/library-voice-receive`, fale no canal e veja os logs de pacotes recebidos.
// Inventário: cobre "Recepção de áudio".

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const voiceReceiveCommandName = "library-voice-receive"

// LibraryVoiceReceive demonstra como pegar pacotes de áudio.
type LibraryVoiceReceive struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryVoiceReceive cria a feature.
func NewLibraryVoiceReceive() bot.Feature { return &LibraryVoiceReceive{} }

// RegisterLibraryVoiceReceive registra manualmente.
func RegisterLibraryVoiceReceive() { bot.RegisterFeature(NewLibraryVoiceReceive()) }

func (f *LibraryVoiceReceive) Name() string { return "LibraryVoiceReceive" }

func (f *LibraryVoiceReceive) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			voiceReceiveCommandName,
			"Conecta e loga pacotes de voz recebidos.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
		binding.Spec.Command.Options = []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "canal",
				Description:  "Canal de voz alvo",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
			},
		}
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryVoiceReceive) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryVoiceReceive) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryVoiceReceive) Intents() discordgo.Intent { return discordgo.IntentGuildVoiceStates }

func (f *LibraryVoiceReceive) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ApplicationCommandData().Options[0].ChannelValue(nil).ID
	guildID := i.GuildID

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Conectando em modo recepção..."},
	})

	go func() {
		vc, err := s.ChannelVoiceJoin(guildID, channelID, true, false)
		if err != nil {
			log.Printf("[Library][Voice] Falha ao conectar (receive): %v", err)
			return
		}

		vc.OpusRecv = make(chan *discordgo.Packet, 10)
		vc.RecvTimeout = 2 * time.Second

		for pkt := range vc.OpusRecv {
			log.Printf("[Library][Voice] Pacote de %s (%d bytes)", pkt.UserID, len(pkt.Opus))
		}
	}()
}
