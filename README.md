# CLI de Teste de Carga (Desafio Técnico - Pós Golang Full Cycle)

Aplicação CLI em Go para realizar testes de carga em serviços HTTP.

Requisitos
- Go 1.20+ instalado
- Docker (obrigatório, para rodar o container)

Estrutura principal
- `src/cmd/main.go` - entrypoint do executável (empacota o pacote `loadtest`)
- `src/main.go` - pacote `loadtest` com a lógica e função testável `Run(args, out)`
- `test/` - testes BDD (Godog features + steps)
- `go.mod` - módulo `pos-goexpert-desafio5`

Executando via Docker

Use exatamente o formato abaixo para executar o teste:

```bash
docker build -t pos-goexpert-teste-de-carga .
docker run pos-goexpert-teste-de-carga --url=http://google.com --requests=1000 --concurrency=10
```

Testes

- Testes unitários + BDD (Godog) executados com `go test`:

```bash
# roda todos os testes (unidade + BDD)
go test ./... -v

# roda apenas os testes BDD
go test ./test -v
```

Observações sobre BDD
- As features estão em `test/features` e os steps em `test/steps`.
- As features usam `github.com/cucumber/godog`; `go test ./...` já baixa a dependência (veja `go.mod`).

Relatório gerado
- Ao final da execução a CLI imprime no console:
	- Tempo total gasto na execução
	- Quantidade total de requests realizados
	- Quantidade de requests com status HTTP 200
	- Distribuição de outros códigos de status HTTP (ex: 404, 500)