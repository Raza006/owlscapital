package utility

import (
	"fmt"
	"log"

	"projectdiscord/services/discordbot/internal/bot"

	"github.com/bwmarrin/discordgo"
)

const (
	receiveUpdatesCommandName = "receive-updates"
	receiveUpdatesButtonID    = "ReceiveUpdates:Subscribe"
	updatesRoleID             = "1449854473013821470"
)

// ReceiveUpdatesFeature manages the /receive-updates command and role assignment.
type ReceiveUpdatesFeature struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewReceiveUpdatesFeature creates a new instance of the ReceiveUpdatesFeature.
func NewReceiveUpdatesFeature() bot.Feature {
	return &ReceiveUpdatesFeature{}
}

// init registers this feature automatically when the package is imported.
func init() {
	bot.RegisterFeature(NewReceiveUpdatesFeature())
}

// Name returns the name of the feature.
func (f *ReceiveUpdatesFeature) Name() string {
	return "ReceiveUpdates"
}

// getBinding centralizes the command definition and handler binding.
func (f *ReceiveUpdatesFeature) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			receiveUpdatesCommandName,
			"Subscribe to receive updates by clicking the button",
			bot.ScopeGuild,
			f.handleReceiveUpdatesCommand,
			true,
		)
	}
	return f.binding
}

// CommandSpecs returns the command specifications.
func (f *ReceiveUpdatesFeature) CommandSpecs() []bot.CommandSpec {
	b := f.getBinding()
	return []bot.CommandSpec{b.Spec}
}

// ApplicationCommandHandlers returns the command handlers.
func (f *ReceiveUpdatesFeature) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

// ComponentHandlers returns the button click handlers.
func (f *ReceiveUpdatesFeature) ComponentHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		receiveUpdatesButtonID: f.handleSubscribeButton,
	}
}

// Intents returns the required Discord gateway intents.
func (f *ReceiveUpdatesFeature) Intents() discordgo.Intent {
	return discordgo.IntentGuilds | discordgo.IntentGuildMembers
}

// handleReceiveUpdatesCommand handles the /receive-updates slash command.
func (f *ReceiveUpdatesFeature) handleReceiveUpdatesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📢 Receive Updates",
		Description: "Click the button below to get notified about important updates!",
		Color:       0x5865F2, // Discord Blurple
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "What you'll receive:",
				Value:  "• Important announcements\n• Feature updates\n• Community news",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "You can unsubscribe at any time by clicking the button again",
		},
	}

	button := discordgo.Button{
		Label:    "🔔 Subscribe to Updates",
		Style:    discordgo.PrimaryButton,
		CustomID: receiveUpdatesButtonID,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{button},
				},
			},
		},
	})

	if err != nil {
		log.Printf("Error responding to /receive-updates command: %v", err)
	}
}

// handleSubscribeButton handles when a user clicks the subscribe button.
func (f *ReceiveUpdatesFeature) handleSubscribeButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response to avoid timeout
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring interaction: %v", err)
		return
	}

	guildID := i.GuildID
	userID := i.Member.User.ID

	// Check if user already has the role
	hasRole := false
	for _, roleID := range i.Member.Roles {
		if roleID == updatesRoleID {
			hasRole = true
			break
		}
	}

	var responseMessage string

	if hasRole {
		// Remove the role if they already have it (toggle behavior)
		err = s.GuildMemberRoleRemove(guildID, userID, updatesRoleID)
		if err != nil {
			log.Printf("Error removing role from user %s: %v", userID, err)
			responseMessage = "❌ Failed to unsubscribe from updates. Please try again or contact an admin."
		} else {
			responseMessage = "✅ You've been unsubscribed from updates. Click the button again if you change your mind!"
		}
	} else {
		// Add the role
		err = s.GuildMemberRoleAdd(guildID, userID, updatesRoleID)
		if err != nil {
			log.Printf("Error adding role to user %s: %v", userID, err)
			responseMessage = "❌ Failed to subscribe to updates. Please try again or contact an admin."
		} else {
			responseMessage = fmt.Sprintf("✅ Success! You now have the <@&%s> role and will receive updates!", updatesRoleID)
		}
	}

	// Follow-up with the response
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &responseMessage,
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
	}
}

