# Biblioteca discordgo Starterkit

Esta pasta hospeda exemplos e features de referência organizados por categoria:

- `session/`: ciclo de vida da `Session`, intents, shards, ratelimit, logging.
- `state/`: uso do cache (`State`), sincronização e coletores de dados.
- `interactions/`: slash commands, context menus, autocomplete, componentes, modals, localização.
- `messaging/`: envio/edição de mensagens, embeds, attachments, threads, reações.
- `events/`: handlers para todos os eventos do gateway.
- `voice/`: conexão de voz, streaming, recepção e monitoramento.
- `oauth/`: fluxos OAuth2, linked roles e autorização externa.
- `rest/`: operações administrativas da REST API (guilds, roles, stickers, automod, etc.).
- `webhooks/`: execução e gerenciamento de webhooks.
- `examples/`: adaptações dos exemplos oficiais do discordgo para o padrão do starterkit.

Cada subpasta deverá conter arquivos Go com exemplos autocontidos, além de documentação opcional (`README.md`) explicando como testar e quais intents/configurações são necessárias.

Consulte `docs/library/index.md` e `docs/library/inventory.md` para acompanhar o progresso e garantir cobertura completa.
