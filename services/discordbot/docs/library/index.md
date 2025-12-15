# Índice da Biblioteca discordgo Starterkit

Use este índice para navegar pelos exemplos que serão adicionados em `internal/library`. Cada categoria aponta para um subdiretório e descreve quais itens do inventário devem ser cobertos.

| Categoria | Subdiretório Planejado | Cobertura Necessária (ver `inventory.md`) |
|-----------|------------------------|-------------------------------------------|
| Sessão & Gateway | `internal/library/session/` | Construção de `Session`, intents, shards, reconexão, ratelimit, logging |
| Estado & Cache | `internal/library/state/` | Uso de `State`, caching, sincronização de guilds e membros |
| Comandos & Interações | `internal/library/interactions/` | Slash commands, context menus, autocomplete, componentes, modals, localização |
| Mensagens & Embeds | `internal/library/messaging/` | Envio/edição, embeds, attachments, threads, reações |
| Eventos Gateway | `internal/library/events/` | Todos os eventos listados em `events.go`/`eventhandlers.go` (guild, canal, mensagem, interação, voz, auto-moderation, subscription etc.) |
| Voz & Áudio | `internal/library/voice/` | Conexão de voz, streaming, voice state/stats |
| OAuth2 & Linked Roles | `internal/library/oauth/` | Fluxos OAuth2, linked roles, autorização externa |
| REST Admin | `internal/library/rest/` | Guilds, channels, members, roles, emojis, stickers, scheduled events, webhooks, auto-moderation |
| Webhooks & Integrações | `internal/library/webhooks/` | Execução de webhooks, interaction responses, integração externa |
| Exemplos Oficiais | `internal/library/examples/` | Portar/adaptar exemplos de `libs/discordgo/examples` com padrão do starterkit |
| Documentação Auxiliar | `docs/library/` | Guia rápido, processos de testes, anotações adicionais |

### Convenção de Arquivos
Cada exemplo deverá seguir o padrão:
- `internal/library/<categoria>/<nome_exemplo>.go`
- Comentários iniciais descrevendo propósito, intents necessárias, passos de teste.
- Função `Register()` ou feature pronta para ser importada em `features_enabled_library.go`.

### Progresso Visual
Considere atualizar o inventário com `✅` quando um exemplo for criado. O índice serve apenas como mapa de diretórios e deve apontar para qualquer material complementar (README, diagramas, etc.).
