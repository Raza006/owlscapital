# Inventário discordgo Starterkit

Este checklist lista todas as superfícies públicas do `discordgo` presentes em `libs/discordgo`. Use para garantir que a biblioteca de exemplos cubra **cada área** antes de avançar para as demais fases.

## Núcleo da Sessão e Conexão
- [✅](../internal/library/session/session_config.go) Construção e gerenciamento de `Session` (`discord.go`)
- [✅](../internal/library/session/session_config.go) Login, tokens e reconexão (`discord.go`, `wsapi.go`)
- [✅](../internal/library/session/reconnect_strategy.go) Gateway Websocket: intents, shards, heartbeats (`wsapi.go`)
- [✅](../internal/library/state/guild_cache_lookup.go) Estado interno e cache (`state.go`)
- [✅](../internal/library/session/session_config.go) Logging configurável (`logging.go`)
- [✅](../internal/library/session/session_config.go) Ratelimit global/local (`ratelimit.go`)
- [✅](../internal/library/session/utilities.go) Utilidades gerais (`util.go`)

## Interações e Comandos
- [✅](../internal/library/interactions/button_basic.go) Slash commands completos (`interactions.go`, `components.go`)
- [✅](../internal/library/interactions/user_context_info.go) Context menu (User) (`interactions.go`)
- [✅](../internal/library/interactions/message_context_quote.go) Context menu (Message) (`interactions.go`)
- [✅](../internal/library/interactions/autocomplete_role.go) Autocomplete (`interactions.go`)
- [✅](../internal/library/interactions/select_menu_multi.go) Componentes (botões, selects, menus) (`components.go`)
- [✅](../internal/library/interactions/modal_feedback.go) Modals (`interactions.go`)
- [✅](../internal/library/interactions/localization_multilang.go) Localização / idiomas (`locales.go`)

## Mensagens e Conteúdo
- [✅](../internal/library/messaging/message_edit_delete.go) Envio/edição/apagamento de mensagens (`message.go`, `restapi.go`)
- [✅](../internal/library/messaging/embed_basic.go) Embeds + rich content (`message.go`)
- [✅](../internal/library/messaging/embed_v2.go) Embeds v2 (containers/sections) (`components.go`)
- [✅](../internal/library/messaging/thread_create.go) Threads (criação, gerenciamento) (`restapi.go`, `events.go`)
- [✅](../internal/library/messaging/reaction_manage.go) Reações (add/remove, remove-all) (`restapi.go`)
- [✅](../internal/library/messaging/attachment_upload.go) Stickers, attachments, componentes especiais (`message.go`)

## Eventos Gateway (Eventos de Discord)
- [ ] Enumerar todos os structs em `events.go`
- [ ] Handlers específicos em `eventhandlers.go`
- [✅](../internal/library/events/guild_events.go) Eventos de guild (create/update/delete, member add/remove, role updates)
- [✅](../internal/library/events/channel_events.go) Eventos de canal/threads/stage (`events.go`)
- [✅](../internal/library/events/thread_events.go) Eventos específicos de thread (`events.go`)
- [✅](../internal/library/events/message_events.go) Eventos de mensagem (create/update/delete, reaction add/remove, bulk delete)
- [✅](../internal/library/events/reaction_events.go) Eventos de reação (`events.go`)
- [✅](../internal/library/events/interaction_events.go) Eventos de interação (`events.go`)
- [✅](../internal/library/events/voice_events.go) Eventos de voz (`voice.go`, `events.go`)
- [✅](../internal/library/events/automod_events.go) Auto moderation, webhooks (`events.go`)
- [✅](../internal/library/events/subscription_events.go) Scheduled events / subscription update (`events.go`)
- [✅](../internal/library/events/stage_events.go) Stage instance (`events.go`)
- [✅](../internal/library/events/gateway_lifecycle.go) Gateway reconnect/ready/resumed heartbeat (`events.go`)

## Voice e Áudio
- [✅](../internal/library/voice/voice_connect.go) Conexão de voz (`voice.go`)
- [✅](../internal/library/voice/voice_receive.go) Streaming de áudio (send/receive) (`voice.go`, `examples/voice_receive`)
- [✅](../internal/library/events/voice_events.go) Voice states e updates (`events.go`)

## OAuth2 e Autenticação
- [✅](../internal/library/oauth/oauth_authorize_link.go) OAuth2 flows (`oauth2.go`)
- [✅](../internal/library/oauth/linked_roles_metadata.go) Linked roles (`examples/linked_roles`)

## REST / Admin
- [✅](../internal/library/rest/application_command_sync.go) Aplicação (commands, guild commands) (`restapi.go`)
- [✅](../internal/library/rest/guild_overview.go) Guilds (create, modify, widgets, onboarding) (`restapi.go`)
- [✅](../internal/library/rest/channel_list.go) Channels (text, voice, threads, stage) (`restapi.go`)
- [✅](../internal/library/rest/member_lookup.go) Members e bans (`restapi.go`)
- [✅](../internal/library/rest/role_list.go) Roles e permissions (`restapi.go`)
- [✅](../internal/library/rest/emoji_list.go) Emojis, stickers (`restapi.go`)
- [✅](../internal/library/rest/scheduled_event_list.go) Scheduled events (`restapi.go`)
- [✅](../internal/library/rest/automod_rules.go) Auto moderation rules (`restapi.go`)
- [✅](../internal/library/webhooks/webhook_execute.go) Webhooks (`webhook.go`, `restapi.go`)
- [✅](../internal/library/events/subscription_events.go) Application/Subscribers (`restapi.go`, `SubscriptionUpdate`)

## Webhook e Integrações Externas
- [✅](../internal/library/webhooks/webhook_execute.go) Webhook execution e gerenciamento (`webhook.go`)
- [✅](../internal/library/webhooks/interaction_followup.go) Interaction responses via webhook (`restapi.go`)

## Exemplos Oficiais (libs/discordgo/examples)
- [ ] Ping/Pong básico
- [ ] Slash commands
- [ ] Autocomplete
- [ ] Components
- [ ] Modals
- [ ] Threads
- [ ] Voice (send/receive)
- [ ] Stage instance / scheduled events
- [ ] Auto moderation
- [ ] Context menus
- [ ] Linked roles
- [ ] Echo / DM pingpong

## Documentação Auxiliar
- [ ] Referências em `libs/discordgo/docs`
- [ ] README do projeto oficial (`libs/discordgo/README.md`)

> Use este inventário como "fonte da verdade". Cada item deve receber um exemplo/documento correspondente na futura biblioteca (`internal/library`).
