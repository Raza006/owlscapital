package utility

import (
	"fmt"
	"log"
	"reflect"

	"projectdiscord/services/discordbot/internal/bot"

	"github.com/bwmarrin/discordgo"
)

const (
	receiveUpdatesCommandName = "receive-updates"
	updatesRoleID             = "1449854473013821470"
	updatesEmojiName          = "owlsnoti"
	updatesEmojiID            = "1449215977068691576"
)

// ReceiveUpdatesFeature manages the /receive-updates command and role assignment via reactions.
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
			"Post the updates notification panel",
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

// TypedEventHandlers returns the event handlers for reaction events.
func (f *ReceiveUpdatesFeature) TypedEventHandlers() bot.TypedEventHandlersMap {
	return bot.TypedEventHandlersMap{
		reflect.TypeOf((*discordgo.MessageReactionAdd)(nil)).Elem():    {f.wrapReactionAdd},
		reflect.TypeOf((*discordgo.MessageReactionRemove)(nil)).Elem(): {f.wrapReactionRemove},
	}
}

// Intents returns the required Discord gateway intents.
func (f *ReceiveUpdatesFeature) Intents() discordgo.Intent {
	return discordgo.IntentGuilds | discordgo.IntentGuildMembers | discordgo.IntentGuildMessageReactions
}

// wrapReactionAdd wraps the typed handler for the event dispatcher.
func (f *ReceiveUpdatesFeature) wrapReactionAdd(s *discordgo.Session, v interface{}) {
	if r, ok := v.(*discordgo.MessageReactionAdd); ok {
		f.handleReactionAdd(s, r)
	}
}

// wrapReactionRemove wraps the typed handler for the event dispatcher.
func (f *ReceiveUpdatesFeature) wrapReactionRemove(s *discordgo.Session, v interface{}) {
	if r, ok := v.(*discordgo.MessageReactionRemove); ok {
		f.handleReactionRemove(s, r)
	}
}

// handleReceiveUpdatesCommand handles the /receive-updates slash command.
func (f *ReceiveUpdatesFeature) handleReceiveUpdatesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send ephemeral response to the command executor
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Updates panel posted!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to /receive-updates command: %v", err)
		return
	}

	// Create the embed with no color
	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("React with <:owlsnoti:1449215977068691576> to get access to https://discord.com/channels/718624848812834903/1449226651064991806 so you can be notified for all alerts in the Observatory."),
	}

	// Send the message as the bot (not showing who executed the command)
	msg, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		log.Printf("Error sending updates panel message: %v", err)
		return
	}

	// Add the custom emoji reaction to the message
	err = s.MessageReactionAdd(i.ChannelID, msg.ID, fmt.Sprintf("%s:%s", updatesEmojiName, updatesEmojiID))
	if err != nil {
		log.Printf("Error adding reaction to message: %v", err)
	}
}

// handleReactionAdd handles when a user adds a reaction to subscribe.
func (f *ReceiveUpdatesFeature) handleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore bot reactions
	if r.UserID == s.State.User.ID {
		return
	}

	// Check if it's the correct emoji
	if r.Emoji.Name != updatesEmojiName || r.Emoji.ID != updatesEmojiID {
		return
	}

	// Add the role to the user
	err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, updatesRoleID)
	if err != nil {
		log.Printf("Error adding role to user %s: %v", r.UserID, err)
		
		// Try to DM the user about the error
		channel, dmErr := s.UserChannelCreate(r.UserID)
		if dmErr == nil {
			s.ChannelMessageSend(channel.ID, "❌ Failed to add the updates role. Please contact an administrator.")
		}
		return
	}

	log.Printf("Added updates role to user %s", r.UserID)

	// Send DM confirmation
	channel, err := s.UserChannelCreate(r.UserID)
	if err != nil {
		log.Printf("Error creating DM channel for user %s: %v", r.UserID, err)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("✅ You now have the <@&%s> role! You'll receive notifications for all alerts in the Observatory.", updatesRoleID))
	if err != nil {
		log.Printf("Error sending DM to user %s: %v", r.UserID, err)
	}
}

// handleReactionRemove handles when a user removes their reaction to unsubscribe.
func (f *ReceiveUpdatesFeature) handleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	// Ignore bot reactions
	if r.UserID == s.State.User.ID {
		return
	}

	// Check if it's the correct emoji
	if r.Emoji.Name != updatesEmojiName || r.Emoji.ID != updatesEmojiID {
		return
	}

	// Remove the role from the user
	err := s.GuildMemberRoleRemove(r.GuildID, r.UserID, updatesRoleID)
	if err != nil {
		log.Printf("Error removing role from user %s: %v", r.UserID, err)
		
		// Try to DM the user about the error
		channel, dmErr := s.UserChannelCreate(r.UserID)
		if dmErr == nil {
			s.ChannelMessageSend(channel.ID, "❌ Failed to remove the updates role. Please contact an administrator.")
		}
		return
	}

	log.Printf("Removed updates role from user %s", r.UserID)

	// Send DM confirmation
	channel, err := s.UserChannelCreate(r.UserID)
	if err != nil {
		log.Printf("Error creating DM channel for user %s: %v", r.UserID, err)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "✅ You've been unsubscribed from Observatory updates. React again if you change your mind!")
	if err != nil {
		log.Printf("Error sending DM to user %s: %v", r.UserID, err)
	}
}

