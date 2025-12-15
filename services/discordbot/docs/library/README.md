# Biblioteca discordgo Starterkit

Este diretório concentra a documentação e o inventário da biblioteca de exemplos que demonstra **todas** as APIs do `discordgo`.

## Como Navegar
- Consulte `inventory.md` para ver tudo que precisa ser coberto.
- Use `index.md` para descobrir o subdiretório certo em `internal/library/`.
- Cada exemplo segue o template descrito em `example_template.md` e expõe uma função `Register...`.

## Como Utilizar com o Cursor
1. Abra este README e o índice correspondente.
2. Ao criar uma nova feature, peça para o Cursor ler o exemplo que cobre o mesmo tópico (ex.: `read_file internal/library/messaging/embed_basic.go`).
3. Siga o template e adapte o código para o seu caso.
4. Atualize `inventory.md` marcando o item com `✅` e adicione o link do novo arquivo.

## Estrutura Atual
- `inventory.md`: checklist mestre.
- `index.md`: mapa de categorias → subpastas.
- `example_template.md`: padrão para novos exemplos.
- `coverage.md`: plano de validação automática.
- embeds v2: consulte `internal/library/messaging/embed_v2.go` para o novo formato baseado em containers/sections exposto pelo discordgo (cobre Section=9, TextDisplay=10, Thumbnail=11, MediaGallery=12, FileComponent=13, Separator=14 e Container=17).

> Sempre que atualizar o `discordgo`, refaça o inventário e gere novos exemplos conforme necessário. A biblioteca deve ser a referência oficial para humanos e para a IA.
