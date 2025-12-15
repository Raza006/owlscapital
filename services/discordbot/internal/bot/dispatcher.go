package bot

import (
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

// FeatureRegistry holds the aggregated, ready-to-use data from all active features.
type FeatureRegistry struct {
	Intents                    discordgo.Intent
	CommandSpecs               []CommandSpec
	ApplicationCommandHandlers InteractionHandlersMap
	ComponentHandlers          InteractionHandlersMap
	ModalSubmitHandlers        InteractionHandlersMap
	AutocompleteHandlers       InteractionHandlersMap
    EventHandlersByType        map[reflect.Type][]EventHandler
}

// LoadFeatures iterates through the active features, aggregates all their requirements, and returns a populated FeatureRegistry.
func LoadFeatures(activeFeatures []Feature) *FeatureRegistry {
	registry := &FeatureRegistry{
		ApplicationCommandHandlers: make(InteractionHandlersMap),
		ComponentHandlers:          make(InteractionHandlersMap),
		ModalSubmitHandlers:        make(InteractionHandlersMap),
		AutocompleteHandlers:       make(InteractionHandlersMap),
        EventHandlersByType:        make(map[reflect.Type][]EventHandler),
	}

    log.Println("Loading features...")
	for _, f := range activeFeatures {
		featureName := f.Name()
		log.Printf("- Loading feature: %s", featureName)

		registry.Intents |= f.Intents()
		registry.CommandSpecs = append(registry.CommandSpecs, f.CommandSpecs()...)

        // Register all handler types, ensuring no conflicts for handlers that require unique IDs.
        registerInteractionMap(registry.ApplicationCommandHandlers, f.ApplicationCommandHandlers(), "ApplicationCommand", featureName)

        // Enforce prefix convention and collision checks for components and modals
        enforcePrefixedIDsAndRegister(registry.ComponentHandlers, f.ComponentHandlers(), featureName, "Component")
        enforcePrefixedIDsAndRegister(registry.ModalSubmitHandlers, f.ModalSubmitHandlers(), featureName, "ModalSubmit")

        registerInteractionMap(registry.AutocompleteHandlers, f.AutocompleteHandlers(), "Autocomplete", featureName)

        // Register typed event handlers (gateway events)
        for t, handlers := range f.TypedEventHandlers() {
            registry.EventHandlersByType[t] = append(registry.EventHandlersByType[t], handlers...)
            log.Printf("  - Typed event handler registered for: %s", t.Name())
        }
	}
	log.Printf("Finished loading features.")
	return registry

}

// registerInteractionMap is a helper to populate a registry map and check for conflicts.
func registerInteractionMap(registryMap, featureMap InteractionHandlersMap, handlerType, featureName string) {
	for id, handler := range featureMap {
		if _, exists := registryMap[id]; exists {
			log.Fatalf("Error: %s handler with ID '%s' from feature '%s' is already registered.", handlerType, id, featureName)
		}
		registryMap[id] = handler
		log.Printf("  - %s handler registered for: %s", handlerType, id)
	}
}

// enforcePrefixedIDsAndRegister valida que todos os IDs com escopo global (componentes e modais)
// sigam a convenção "<FeatureName>:<id>" e que não haja colisões entre features.
// Em caso de violação, o processo é abortado com uma mensagem clara.
func enforcePrefixedIDsAndRegister(registryMap, featureMap InteractionHandlersMap, featureName, handlerType string) {
    prefix := featureName + ":"
    for id, handler := range featureMap {
        if !hasPrefix(id, prefix) {
            log.Fatalf("Error: %s handler ID '%s' from feature '%s' must start with prefix '%s'", handlerType, id, featureName, prefix)
        }
        if _, exists := registryMap[id]; exists {
            log.Fatalf("Error: %s handler with ID '%s' is already registered by another feature.", handlerType, id)
        }
        registryMap[id] = handler
        log.Printf("  - %s handler registered for: %s", handlerType, id)
    }
}

// hasPrefix é uma pequena ajuda para evitar importar strings e manter o estilo do arquivo.
func hasPrefix(s, prefix string) bool {
    if len(s) < len(prefix) {
        return false
    }
    return s[:len(prefix)] == prefix
}

// RegisterCommands uses the provided registry to register all application commands with Discord.
func RegisterCommands(s *discordgo.Session, guildID string, registry *FeatureRegistry) {
	if len(registry.CommandSpecs) == 0 {
		log.Println("No application commands to register.")
		return
	}

	var guildCommands []*discordgo.ApplicationCommand
	var globalCommands []*discordgo.ApplicationCommand
    for _, spec := range registry.CommandSpecs {
        if spec.Command == nil {
            continue
        }
        switch spec.Scope {
        case ScopeGuild:
            guildCommands = append(guildCommands, spec.Command)
        case ScopeGlobal:
            globalCommands = append(globalCommands, spec.Command)
        }
    }

	if len(guildCommands) > 0 {
		if guildID == "" {
			log.Printf("Warning: %d guild-scoped commands configured, but GUILD_ID is empty. Skipping guild registration.", len(guildCommands))
		} else {
			log.Printf("Registering %d guild commands...", len(guildCommands))
			if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, guildCommands); err != nil {
				log.Fatalf("Cannot register guild commands: %v", err)
			}
			log.Printf("Successfully registered %d guild commands.", len(guildCommands))
		}
	}

	if len(globalCommands) > 0 {
		log.Printf("Registering %d global commands...", len(globalCommands))
		if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", globalCommands); err != nil {
			log.Fatalf("Cannot register global commands: %v", err)
		}
		log.Printf("Successfully registered %d global commands.", len(globalCommands))
	}
}

// dispatchEvent is a helper to safely dispatch events to their handlers in separate goroutines.
func dispatchEvent(s *discordgo.Session, handlers []EventHandler, eventData interface{}, eventName string) {
	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic recovered in event handler for %s: %v", eventName, r)
				}
			}()
			h(s, eventData)
		}(handler)
	}
}

// CreateDispatcher creates a single, comprehensive event handler that routes all incoming events
// and interactions from Discord to the appropriate, registered feature handlers.
// This version is DRY, using the type switch only to determine the event name.
func CreateDispatcher(registry *FeatureRegistry) func(s *discordgo.Session, v interface{}) {
	return func(s *discordgo.Session, v interface{}) {
		// First, handle interactions, which have a specific structure.
		if i, ok := v.(*discordgo.InteractionCreate); ok {
			var handler InteractionHandler
			var id string

			// This switch is now cleaner, directly assigning the handler and relying
			// on the fact that a failed map lookup returns a nil handler.
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				id = i.ApplicationCommandData().Name
				handler = registry.ApplicationCommandHandlers[id]
			case discordgo.InteractionMessageComponent:
				id = i.MessageComponentData().CustomID
				handler = registry.ComponentHandlers[id]
			case discordgo.InteractionModalSubmit:
				id = i.ModalSubmitData().CustomID
				handler = registry.ModalSubmitHandlers[id]
			case discordgo.InteractionApplicationCommandAutocomplete:
				id = i.ApplicationCommandData().Name
				handler = registry.AutocompleteHandlers[id]
			case discordgo.InteractionPing:
				// The library usually handles this automatically.
				// Responding here is only necessary if you want to override the default behavior.
				// For now, we'll just log that we received it.
				log.Println("Received a Ping interaction.")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponsePong,
				})
				return
			}

			// The check is now simpler, as a non-existent handler will be nil.
			if handler != nil {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Panic recovered in interaction handler for ID '%s': %v", id, r)
						}
					}()
					handler(s, i)
				}()
			}
			// Interactions are self-contained, so we can return early.
			return
		}

        // For all other gateway events, determine the concrete type via reflection and dispatch.
        t := reflect.TypeOf(v)
        if t == nil {
            return
        }
        if t.Kind() == reflect.Ptr {
            t = t.Elem()
        }
        if handlers, ok := registry.EventHandlersByType[t]; ok && len(handlers) > 0 {
            dispatchEvent(s, handlers, v, t.Name())
        }
	}
}


