package utility

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/assets"
	"projectdiscord/services/discordbot/internal/bot"
)

const (
	supportPanelCommandName          = "supportpanel"
	supportPanelButtonID             = "SupportPanel:OpenTicket"
	supportPanelModalID              = "SupportPanel:TicketModal"
	supportPanelSubjectInputID       = "SupportPanel:Subject"
	supportPanelDescriptionInputID   = "SupportPanel:Description"
	supportPanelCloseButtonID        = "SupportPanel:CloseTicket"
	supportPanelCloseConfirmID       = "SupportPanel:ConfirmClose"
	supportPanelCloseCancelID        = "SupportPanel:CancelClose"
	supportPanelChannelID            = "1441902437723017349"
	supportPanelStaffUserID          = "408370072802033664"
	supportPanelTicketPrefix         = "Support-"
	supportPanelThreadArchiveMinutes = 10080
)

// SupportPanelFeature manages the /supportpanel workflow and ticket lifecycle.
type SupportPanelFeature struct {
	bot.BaseFeature
	binding      bot.CommandBinding
	ticketOwners map[string]string
	userTickets  map[string]string
	ticketMu     sync.RWMutex
}

func init() {
	bot.RegisterFeature(NewSupportPanelFeature())
}

// NewSupportPanelFeature wires the feature so it can be registered.
func NewSupportPanelFeature() bot.Feature {
	return &SupportPanelFeature{
		ticketOwners: make(map[string]string),
		userTickets:  make(map[string]string),
	}
}

func (f *SupportPanelFeature) Name() string { return "SupportPanel" }

func (f *SupportPanelFeature) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			supportPanelCommandName,
			"Post the technical support panel in the current channel.",
			bot.ScopeGuild,
			f.handleSupportPanelCommand,
			false,
		)
	}
	return f.binding
}

func (f *SupportPanelFeature) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *SupportPanelFeature) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *SupportPanelFeature) ComponentHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		supportPanelButtonID:       f.handleOpenTicketButton,
		supportPanelCloseButtonID:  f.handleCloseTicketButton,
		supportPanelCloseConfirmID: f.handleCloseTicketConfirm,
		supportPanelCloseCancelID:  f.handleCloseTicketCancel,
	}
}

func (f *SupportPanelFeature) ModalSubmitHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		supportPanelModalID: f.handleTicketModalSubmit,
	}
}

func (f *SupportPanelFeature) Intents() discordgo.Intent {
	return discordgo.IntentGuilds
}

func (f *SupportPanelFeature) handleSupportPanelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !memberHasAdmin(i.Member) {
		respondEphemeral(s, i, "Only administrators can execute this command.")
		return
	}

	panelMessage := buildSupportPanelMessage()
	_, sendErr := s.ChannelMessageSendComplex(i.ChannelID, panelMessage)

	response := "Support panel posted."
	if sendErr != nil {
		log.Printf("[SupportPanel] failed to post panel: %v", sendErr)
		response = "We could not post the support panel. Please try again."
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: response,
		},
	})
	if err != nil {
		log.Printf("[SupportPanel] failed to acknowledge /supportpanel: %v", err)
	}
}

func (f *SupportPanelFeature) handleOpenTicketButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := userIDFromInteraction(i)
	if userID == "" {
		respondEphemeral(s, i, "We couldn't identify your account. Please try again.")
		return
	}

	if existing := f.getTicketForUser(userID); existing != "" {
		respondEphemeral(s, i, fmt.Sprintf("You already have an open ticket: <#%s>. Please close it before opening another.", existing))
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: supportPanelModalID,
			Title:    "Open Technical Support Ticket",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    supportPanelSubjectInputID,
							Label:       "Subject",
							Style:       discordgo.TextInputShort,
							Placeholder: "e.g., Account connection failing",
							MinLength:   4,
							MaxLength:   80,
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    supportPanelDescriptionInputID,
							Label:       "Description",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Share relevant details, logs, or screenshots.",
							MinLength:   10,
							MaxLength:   1000,
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("[SupportPanel] failed to open ticket modal: %v", err)
	}
}

func (f *SupportPanelFeature) handleTicketModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := userIDFromInteraction(i)
	if userID == "" {
		respondEphemeral(s, i, "We couldn't identify your account. Please try again.")
		return
	}
	if existing := f.getTicketForUser(userID); existing != "" {
		respondEphemeral(s, i, fmt.Sprintf("You already have an open ticket: <#%s>. Please close it before opening another.", existing))
		return
	}

	subject := strings.TrimSpace(extractModalValue(i, supportPanelSubjectInputID))
	description := strings.TrimSpace(extractModalValue(i, supportPanelDescriptionInputID))

	if subject == "" || description == "" {
		respondEphemeral(s, i, "Both the subject and description are required.")
		return
	}

	threadName := fmt.Sprintf("%s | %s%s", memberDisplayName(i.Member), supportPanelTicketPrefix, randomTicketDigits())
	thread, err := s.ThreadStartComplex(supportPanelChannelID, &discordgo.ThreadStart{
		Name:                threadName,
		AutoArchiveDuration: supportPanelThreadArchiveMinutes,
		Type:                discordgo.ChannelTypeGuildPrivateThread,
		Invitable:           false,
	})
	if err != nil {
		log.Printf("[SupportPanel] failed to create ticket thread: %v", err)
		respondEphemeral(s, i, "We could not create your ticket. Please try again or alert the staff.")
		return
	}

	if err := s.ThreadMemberAdd(thread.ID, userID); err != nil {
		log.Printf("[SupportPanel] failed to add member %s to thread %s: %v", userID, thread.ID, err)
	}

	f.setTicket(thread.ID, userID)

	_, err = s.ChannelMessageSendComplex(thread.ID, buildTicketMessage(subject, description, userID))
	if err != nil {
		log.Printf("[SupportPanel] failed to send ticket details for thread %s: %v", thread.ID, err)
	}

	respondEphemeral(s, i, fmt.Sprintf("Ticket created: <#%s>. Our admins will follow up shortly.", thread.ID))
}

func (f *SupportPanelFeature) handleCloseTicketButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !memberHasAdmin(i.Member) {
		respondEphemeral(s, i, "Only administrators can close tickets.")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Are you sure you want to delete this ticket? This will remove the entire thread for everyone.",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style:    discordgo.DangerButton,
							Label:    "Yes",
							CustomID: supportPanelCloseConfirmID,
						},
						discordgo.Button{
							Style:    discordgo.SecondaryButton,
							Label:    "Cancel",
							CustomID: supportPanelCloseCancelID,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("[SupportPanel] failed to send close confirmation: %v", err)
	}
}

func (f *SupportPanelFeature) handleCloseTicketCancel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Ticket deletion cancelled.",
			Components: []discordgo.MessageComponent{},
		},
	}); err != nil {
		log.Printf("[SupportPanel] failed to acknowledge close cancel: %v", err)
	}
}

func (f *SupportPanelFeature) handleCloseTicketConfirm(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !memberHasAdmin(i.Member) {
		respondEphemeral(s, i, "Only administrators can close tickets.")
		return
	}

	requesterID := f.getTicketOwner(i.ChannelID)
	if _, err := s.ChannelDelete(i.ChannelID); err != nil {
		log.Printf("[SupportPanel] failed to delete ticket channel %s: %v", i.ChannelID, err)
		respondEphemeral(s, i, "Something went wrong while deleting the ticket. Please try again.")
		return
	}
	f.clearTicketByThread(i.ChannelID)

	if requesterID != "" {
		dmChannel, err := s.UserChannelCreate(requesterID)
		if err != nil {
			log.Printf("[SupportPanel] failed to open DM with %s: %v", requesterID, err)
		} else {
			if _, err := s.ChannelMessageSend(dmChannel.ID, "Your support ticket has been closed. If you have any other inquiries then feel free to open another support ticket if needed."); err != nil {
				log.Printf("[SupportPanel] failed to DM %s about ticket closure: %v", requesterID, err)
			}
		}
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Ticket deleted and the requester has been notified.",
			Components: []discordgo.MessageComponent{},
		},
	}); err != nil {
		log.Printf("[SupportPanel] failed to acknowledge ticket deletion: %v", err)
	}
}

func buildSupportPanelContainer() discordgo.Container {
	return discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media: discordgo.UnfurledMediaItem{
							URL: "attachment://" + assets.SupportBannerFilename,
						},
					},
				},
			},
			discordgo.Section{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{Content: "**Need Technical Help?**"},
					discordgo.TextDisplay{Content: "Tap the button on the right to open a private **technical support** ticket with our team."},
				},
				Accessory: discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Support Ticket",
					CustomID: supportPanelButtonID,
					Emoji: &discordgo.ComponentEmoji{
						Name: "discotoolsxyzicon2",
						ID:   "1439322459902971944",
					},
				},
			},
			discordgo.Separator{},
			discordgo.TextDisplay{Content: "### Support Rules"},
			discordgo.TextDisplay{
				Content: strings.Join([]string{
					"- Summarize your **technical issue** clearly in the subject (e.g., login problems, account connection, payment or access issues).",
					"- Include any relevant details or screenshots in the description.",
					"- Tickets should be opened **only for technical issues**. General questions should not be submitted here.",
					"- Be respectful and allow staff time to respond.",
				}, "\n"),
			},
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media: discordgo.UnfurledMediaItem{
							URL: "attachment://" + assets.OwlsFooterFilename,
						},
					},
				},
			},
		},
	}
}

func buildSupportPanelMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Flags:      discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{buildSupportPanelContainer()},
		Files:      buildBrandingFiles(),
	}
}

func buildTicketMessage(subject, description, requesterID string) *discordgo.MessageSend {
	container := discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media: discordgo.UnfurledMediaItem{
							URL: "attachment://" + assets.SupportBannerFilename,
						},
					},
				},
			},
			discordgo.TextDisplay{
				Content: fmt.Sprintf("<@%s> will help with this request. Thanks %s for reaching out.", supportPanelStaffUserID, mentionOrName(requesterID)),
			},
			discordgo.Separator{},
			discordgo.TextDisplay{
				Content: fmt.Sprintf("**Subject:** %s", subject),
			},
			discordgo.TextDisplay{
				Content: fmt.Sprintf("**Description:**\n%s", description),
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.DangerButton,
						Label:    "Close Ticket",
						CustomID: supportPanelCloseButtonID,
					},
				},
			},
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media: discordgo.UnfurledMediaItem{
							URL: "attachment://" + assets.OwlsFooterFilename,
						},
					},
				},
			},
		},
	}

	mentions := []string{supportPanelStaffUserID}
	if requesterID != "" {
		mentions = append(mentions, requesterID)
	}

	return &discordgo.MessageSend{
		Flags:      discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{container},
		Files:      buildBrandingFiles(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: mentions,
		},
	}
}

func buildBrandingFiles() []*discordgo.File {
	return []*discordgo.File{
		{
			Name:   assets.SupportBannerFilename,
			Reader: bytes.NewReader(assets.SupportBanner),
		},
		{
			Name:   assets.OwlsFooterFilename,
			Reader: bytes.NewReader(assets.OwlsFooter),
		},
	}
}

func extractModalValue(i *discordgo.InteractionCreate, customID string) string {
	data := i.ModalSubmitData()

	for _, component := range data.Components {
		row, ok := component.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, inner := range row.Components {
			input, ok := inner.(*discordgo.TextInput)
			if !ok {
				continue
			}
			if input.CustomID == customID {
				return input.Value
			}
		}
	}
	return ""
}

func memberHasAdmin(m *discordgo.Member) bool {
	if m == nil {
		return false
	}
	return m.Permissions&discordgo.PermissionAdministrator != 0
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if i == nil || s == nil {
		return
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
	if err != nil {
		log.Printf("[SupportPanel] failed to respond ephemerally: %v", err)
	}
}

func memberDisplayName(m *discordgo.Member) string {
	if m == nil {
		return "Member"
	}
	if m.Nick != "" {
		return m.Nick
	}
	if m.User != nil && m.User.GlobalName != "" {
		return m.User.GlobalName
	}
	if m.User != nil && m.User.Username != "" {
		return m.User.Username
	}
	return "Member"
}

func mentionOrName(userID string) string {
	if userID == "" {
		return "the member"
	}
	return fmt.Sprintf("<@%s>", userID)
}

func randomTicketDigits() string {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		log.Printf("[SupportPanel] failed to generate ticket digits: %v", err)
		return "0000"
	}
	return fmt.Sprintf("%04d", n.Int64())
}

func boolPtr(v bool) *bool {
	return &v
}

func (f *SupportPanelFeature) getTicketOwner(threadID string) string {
	f.ticketMu.RLock()
	defer f.ticketMu.RUnlock()
	return f.ticketOwners[threadID]
}

func (f *SupportPanelFeature) getTicketForUser(userID string) string {
	if userID == "" {
		return ""
	}
	f.ticketMu.RLock()
	defer f.ticketMu.RUnlock()
	return f.userTickets[userID]
}

func (f *SupportPanelFeature) setTicket(threadID, userID string) {
	if threadID == "" || userID == "" {
		return
	}
	f.ticketMu.Lock()
	defer f.ticketMu.Unlock()
	f.ticketOwners[threadID] = userID
	f.userTickets[userID] = threadID
}

func (f *SupportPanelFeature) clearTicketByThread(threadID string) string {
	f.ticketMu.Lock()
	defer f.ticketMu.Unlock()
	userID := f.ticketOwners[threadID]
	delete(f.ticketOwners, threadID)
	if userID != "" {
		delete(f.userTickets, userID)
	}
	return userID
}

func userIDFromInteraction(i *discordgo.InteractionCreate) string {
	if i == nil || i.Member == nil || i.Member.User == nil {
		return ""
	}
	return i.Member.User.ID
}
