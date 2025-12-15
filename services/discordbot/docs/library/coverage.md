# Cobertura da Biblioteca discordgo

Este arquivo descreve como assegurar que cada item listado em `inventory.md` possua um exemplo correspondente.

## Processo Manual (provisório)
1. Mantenha `docs/library/inventory.md` como checklist mestre.
2. Ao criar um exemplo, adicione um link para o arquivo na linha correspondente e marque com `✅`.
3. Utilize `docs/library/index.md` para verificar se a categoria possui subdiretório e README de apoio.

## Automação Planejada
- Criar script Go (`cmd/librarycheck`) que:
  - Carrega `inventory.md` e procura por itens sem `✅`.
  - Verifica se existe arquivo correspondente em `internal/library/<categoria>`.
  - Opcional: falha com código ≠ 0 quando encontrar pendências (para uso em CI).
- Checklist mínimo para o script:
  - Ler inventário usando análise de markdown simples (linhas com `- [ ]` ou `- [✅]`).
  - Mapear diretórios de `internal/library/` para categorias (usar `docs/library/index.md`).
  - Gerar relatório com porcentagem de cobertura.

## Integração com CI (futuro)
- Adicionar job em GitHub Actions:
  1. `go run ./cmd/librarycheck`.
  2. Falhar se houver itens não marcados ou se diretórios/arquivos esperados não existirem.
  3. Publicar resumo nos logs.

> Enquanto o script não estiver implementado, use esta documentação para realizar revisão manual a cada mudança.
