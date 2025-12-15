package messaging

// Example: showcases the “embed v2” building blocks (containers, sections, text display, accessories, media gallery, file component).
// Required intents: Guilds (slash command only).
// Test: register `/library-embed-v2`; the response renders three containers instead of the classic embed.
// Inventory: covers “Embeds v2 (containers)”.
//
// Quick checklist to avoid common pitfalls with containers:
// - Every Section **must** define an accessory (Button or Thumbnail). Use a disabled link button for static highlights.
// - Non-link buttons require globally unique CustomID values; follow the `<FeatureName>:<slug>` convention.
// - The total number of components (including nested sections) cannot exceed 40. Keep the layout lean.
// - Always set `MessageFlagsIsComponentsV2`; omit legacy fields such as `Content` and `Embeds` in the same payload.
// - Attachments referenced by a `FileComponent` must also be listed inside `Files`.

import (
	"bytes"
	"strings"

	"github.com/bwmarrin/discordgo"

	"projectdiscord/services/discordbot/internal/bot"
)

const (
	embedV2CommandName  = "library-embed-v2"
	embedV2SelectCustom = "LibraryEmbedV2:topics"
	embedV2Attachment   = "embed-v2-notes.txt"
)

// LibraryEmbedV2Feature demonstrates how to assemble a multi-container response using the new components system.
type LibraryEmbedV2Feature struct {
	bot.BaseFeature
	binding bot.CommandBinding
}

// NewLibraryEmbedV2Feature creates the feature.
func NewLibraryEmbedV2Feature() bot.Feature {
	return &LibraryEmbedV2Feature{}
}

// RegisterLibraryEmbedV2Feature allows manual registration.
func RegisterLibraryEmbedV2Feature() {
	bot.RegisterFeature(NewLibraryEmbedV2Feature())
}

func (f *LibraryEmbedV2Feature) Name() string { return "LibraryEmbedV2" }

func (f *LibraryEmbedV2Feature) getBinding() bot.CommandBinding {
	if f.binding.Spec.Command == nil {
		f.binding = bot.NewSlashCommandBinding(
			embedV2CommandName,
			"Showcases how to compose a full embed v2 message.",
			bot.ScopeGuild,
			f.handleCommand,
			true,
		)
	}
	return f.binding
}

func (f *LibraryEmbedV2Feature) CommandSpecs() []bot.CommandSpec {
	return []bot.CommandSpec{f.getBinding().Spec}
}

func (f *LibraryEmbedV2Feature) ApplicationCommandHandlers() bot.InteractionHandlersMap {
	return f.getBinding().AppCommandHandlers
}

func (f *LibraryEmbedV2Feature) ComponentHandlers() bot.InteractionHandlersMap {
	return bot.InteractionHandlersMap{
		embedV2SelectCustom: f.handleSelect,
	}
}

func (f *LibraryEmbedV2Feature) Intents() discordgo.Intent { return discordgo.IntentGuilds }

func (f *LibraryEmbedV2Feature) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	topicsSelect := buildTopicsSelect()

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
			Components: []discordgo.MessageComponent{
				buildHeroContainer(),
				buildTopicsContainer(topicsSelect),
				buildSpotlightContainer(),
			},
			Files: []*discordgo.File{
				{Name: embedV2Attachment, Reader: buildAttachmentNotes()},
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, message)
}

func (f *LibraryEmbedV2Feature) handleSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	values := i.MessageComponentData().Values
	if len(values) == 0 {
		values = []string{"(nothing selected)"}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "You explored: " + strings.Join(values, ", "),
		},
	})
}

func buildHeroContainer() discordgo.Container {
	accent := 0x5865F2
	return discordgo.Container{
		AccentColor: &accent,
		Components: []discordgo.MessageComponent{
			discordgo.Section{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{Content: "**Embed v2 Starter Pack**"},
					discordgo.TextDisplay{Content: "Containers unlock flexible layouts, modular sections, and richer storytelling than the legacy embed payload."},
				},
				Accessory: discordgo.Button{
					Style: discordgo.LinkButton,
					Label: "Developer docs",
					URL:   "https://discord.com/developers/docs/resources/message#layout-components",
				},
			},
			discordgo.Separator{},
			discordgo.TextDisplay{Content: "_Follow the guide below to tailor layouts, actions, and media for your experience._"},
		},
	}
}

func buildTopicsContainer(selectMenu discordgo.SelectMenu) discordgo.Container {
	accent := 0x57F287

	return discordgo.Container{
		AccentColor: &accent,
		Components: []discordgo.MessageComponent{
			discordgo.Section{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{Content: "**Hands-on practice**"},
					discordgo.TextDisplay{Content: "Pick one or more building blocks to explore templates, gotchas, and recommended patterns."},
				},
				Accessory: discordgo.Button{
					Style:    discordgo.LinkButton,
					Label:    "Component glossary",
					URL:      "https://discord.com/developers/docs/interactions/message-components#component-object-component-types",
					Disabled: false,
				},
			},
			discordgo.Separator{},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					selectMenu,
				},
			},
		},
	}
}

func buildSpotlightContainer() discordgo.Container {
	accent := 0xFEE75C
	imageURL := "https://images.unsplash.com/photo-1523580846011-d3a5bc25702b?auto=format&fit=crop&w=900&q=80"
	description := "Example gallery item"

	return discordgo.Container{
		AccentColor: &accent,
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{Content: "**Media + attachments**"},
			discordgo.TextDisplay{Content: "Combine galleries, files, and follow-up actions to craft richer onboarding or marketing beats."},
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media:       discordgo.UnfurledMediaItem{URL: imageURL},
						Description: stringPtr(description),
					},
				},
			},
			discordgo.FileComponent{
				File:    discordgo.UnfurledMediaItem{URL: "attachment://" + embedV2Attachment},
				Spoiler: false,
			},
		},
	}
}

func buildTopicsSelect() discordgo.SelectMenu {
	min := 1
	return discordgo.SelectMenu{
		CustomID:    embedV2SelectCustom,
		MenuType:    discordgo.StringSelectMenu,
		Placeholder: "Select areas to deep dive",
		MinValues:   &min,
		MaxValues:   3,
		Options: []discordgo.SelectMenuOption{
			{Label: "Layout anatomy", Value: "layout", Description: "Containers, sections, separators"},
			{Label: "Component rules", Value: "components", Description: "Buttons, accessories, IDs"},
			{Label: "Media tricks", Value: "media", Description: "Galleries, thumbnails, files"},
			{Label: "Migration path", Value: "migration", Description: "From embeds 1.0 to containers"},
		},
	}
}

func buildAttachmentNotes() *bytes.Buffer {
	content := "Embed v2 Starter Notes\n- Containers require MessageFlagsIsComponentsV2\n- No legacy Content/Embeds fields allowed in the same payload\n- File components must reference attachments via attachment:// URLs\n"
	return bytes.NewBufferString(content)
}

func stringPtr(v string) *string { return &v }
