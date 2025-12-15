package voice

// Exemplo: conecta em um canal de voz e envia alguns frames de silêncio.
// Intents necessárias: GuildVoiceStates e permissão de conectar/falar.
// Teste: execute `/library-voice-connect` indicando um canal de voz e veja o bot entrar, enviar áudio silencioso e sair.
// Inventário: cobre "Conexão de voz" e envio básico.

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const voiceConnectCommandName = "library-voice-connect"

// LibraryVoiceConnect demonstra conexão simples ao canal de voz.
type LibraryVoiceConnect struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryVoiceConnect cria a feature.
func NewLibraryVoiceConnect() bot.Feature { return &LibraryVoiceConnect{} }

// RegisterLibraryVoiceConnect registra manualmente.
func RegisterLibraryVoiceConnect() { bot.RegisterFeature(NewLibraryVoiceConnect()) }

func (f *LibraryVoiceConnect) Name() string { return "LibraryVoiceConnect" }

func (f *LibraryVoiceConnect) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		binding := bot.NewSlashCommandBinding(
			voiceConnectCommandName,
			"Entra em um canal de voz e envia áudio silencioso.",
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
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice, discordgo.ChannelTypeGuildStageVoice},
			},
		}
		f.binding = binding
	}
	return f.binding
}

func (f *LibraryVoiceConnect) CommandSpecs() []bot.CommandSpec { return []bot.CommandSpec{f.getBinding().Spec} }

func (f *LibraryVoiceConnect) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryVoiceConnect) Intents() discordgo.Intent {
	return discordgo.IntentGuildVoiceStates
}

func (f *LibraryVoiceConnect) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ApplicationCommandData().Options[0].ChannelValue(nil).ID
	guildID := i.GuildID

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Conectando ao canal de voz..."},
	})

	go func() {
		vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
		if err != nil {
			log.Printf("[Library][Voice] Falha ao conectar: %v", err)
			return
		}
		defer vc.Disconnect()

		vc.Speaking(true)
		defer vc.Speaking(false)

		frame := make([]byte, 960*2) // 20ms de silêncio (PCM)
		for x := 0; x < 50; x++ {
			vc.OpusSend <- frame
			time.Sleep(20 * time.Millisecond)
		}

		log.Println("[Library][Voice] Áudio exemplo finalizado, saindo do canal.")
	}()
}
