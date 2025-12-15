# Template de Exemplo (`internal/library`)

Cada arquivo de exemplo deve seguir o padrão abaixo:

```go
// Pacote: internal/library/<categoria>
// Nome do arquivo: <slug>_example.go

// 1) Comentário de cabeçalho com:
//    - Resumo em linguagem simples (o que o exemplo demonstra).
//    - Intents necessárias.
//    - Passos para testar no Discord.
//    - Referências ao inventário (`docs/library/inventory.md`).
//
// 2) Estrutura principal:
//    type <NomeFeature> struct { bot.BaseFeature }
//
// 3) Função construtora:
//    func New<NomeFeature>() bot.Feature
//
// 4) Métodos obrigatórios implementados (CommandSpecs, Handlers etc.).
//
// 5) Função auxiliar para registrar manualmente:
//    func Register<NomeFeature>() { bot.RegisterFeature(New<NomeFeature>()) }
//
// 6) Comentários inline explicando pontos importantes (ex.: por que usar certo campo).
//
// 7) Se o exemplo precisar de strings constantes (IDs, textos), definir variáveis no topo.
//
// 8) Ao final, incluir comentário apontando para próximos passos ou variações.
```

Observações:
- **Não** usar `init()` para registrar automaticamente; a ideia é que o usuário importe e chame explicitamente o registrador do exemplo.
- Utilize nomes descritivos e prefixo `Library` para evitar colisões (`LibraryEmbedBasicFeature`).
- Inclua `lint:disable` somente se necessário e sempre com justificativa.
- Pré-requisite que qualquer helper adicional seja referenciado (quando existir) para manter consistência.
