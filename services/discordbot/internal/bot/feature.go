package bot

import (
    "reflect"

    "github.com/bwmarrin/discordgo"
)

// CommandScope defines where a command should be registered.
// ScopeGuild: registro por guild (propagação rápida). ScopeGlobal: registro global (propagação lenta).
type CommandScope int

const (
	ScopeGuild CommandScope = iota
	ScopeGlobal
)

// CommandSpec bundles a command definition with its intended registration scope.
type CommandSpec struct {
	Command *discordgo.ApplicationCommand
	Scope   CommandScope
}

// EventHandler is a generic function that can handle any type of event from discordgo.
type EventHandler func(s *discordgo.Session, v interface{})

// InteractionHandler is a specific handler for all types of InteractionCreate events.
type InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// TypedEventHandlersMap maps a concrete gateway event type to one or more handlers.
// Example key: reflect.TypeOf((*discordgo.MessageCreate)(nil)).Elem()
type TypedEventHandlersMap map[reflect.Type][]EventHandler

// InteractionHandlersMap maps a CustomID (for components/modals) or a command name (for commands/autocomplete) to a handler.
type InteractionHandlersMap map[string]InteractionHandler

// Feature is the definitive interface for any bot module.
// It allows a module to register commands and handlers for any Discord event or interaction.
type Feature interface {
	Name() string
	// CommandSpecs returns the commands and where each should be registered (guild/global).
	CommandSpecs() []CommandSpec
	Intents() discordgo.Intent

	// ApplicationCommandHandlers returns a map of command names to their handlers.
	// This single method handles all application command types: slash, user context, and message context.
	// The key is the command name.
	ApplicationCommandHandlers() InteractionHandlersMap

	// ComponentHandlers returns a map of component CustomIDs to their handlers (for buttons, select menus).
	ComponentHandlers() InteractionHandlersMap

	// ModalSubmitHandlers returns a map of modal CustomIDs to their submission handlers.
	ModalSubmitHandlers() InteractionHandlersMap

	// AutocompleteHandlers returns a map of command names to their autocomplete handlers.
	AutocompleteHandlers() InteractionHandlersMap

    // TypedEventHandlers returns handlers keyed by the concrete gateway event type.
    // This removes reliance on string event names and avoids typo-prone keys.
    TypedEventHandlers() TypedEventHandlersMap
}

// BaseFeature provides a default, empty implementation of the Feature interface.
// Embed this in feature structs to avoid boilerplate for unused methods.
type BaseFeature struct{}

func (b *BaseFeature) Name() string { return "BaseFeature" }

func (b *BaseFeature) CommandSpecs() []CommandSpec { return nil }

func (b *BaseFeature) Intents() discordgo.Intent { return 0 }

func (b *BaseFeature) ApplicationCommandHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) ComponentHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) ModalSubmitHandlers() InteractionHandlersMap { return nil }

func (b *BaseFeature) AutocompleteHandlers() InteractionHandlersMap { return nil }
func (b *BaseFeature) TypedEventHandlers() TypedEventHandlersMap { return nil }

// --- Global feature registry for auto-registration via init() ---

var (
	registeredFeatures []Feature
)

// RegisterFeature allows a feature package to register itself during init().
// This enables plug-and-play addition of new features without editing a central list.
func RegisterFeature(f Feature) {
	registeredFeatures = append(registeredFeatures, f)
}

// RegisteredFeatures returns the slice of all features that called RegisterFeature.
func RegisteredFeatures() []Feature { return registeredFeatures }

// --- Helpers to avoid drift between command specs and handlers ---

// CommandBinding groups a command specification with its associated application command handlers.
// Using this, a feature can declare a command once and automatically bind the handler by name,
// eliminating duplication and the risk of mismatched identifiers.
type CommandBinding struct {
	Spec               CommandSpec
	AppCommandHandlers InteractionHandlersMap
}

// NewSlashCommandBinding creates a binding for a slash command with a single handler.
// The command name is used both in the spec and as the key in the handlers map, preventing drift.
func NewSlashCommandBinding(name, description string, scope CommandScope, handler InteractionHandler, dmPermission bool) CommandBinding {
	cmd := &discordgo.ApplicationCommand{
		Name:         name,
		Description:  description,
		DMPermission: func(b bool) *bool { return &b }(dmPermission),
	}
	return CommandBinding{
		Spec: CommandSpec{Command: cmd, Scope: scope},
		AppCommandHandlers: InteractionHandlersMap{
			name: handler,
		},
	}
}


