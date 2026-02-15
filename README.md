# c6bank-transactions

Serviço web em Go que processa extratos de transações do C6 Bank e os converte para os formatos QIF/CSV para softwares de finanças pessoais.

## Funcionalidades

- **Múltiplos Formatos de Entrada**: Extratos PDF, arquivos CSV e capturas de tela de celular
- **Processamento Inteligente OCR**: Recorte inteligente para modelos de iPhone com OCR em português+inglês
- **Exportação QIF/CSV**: Geração de arquivos compatíveis com aplicativos de finanças pessoais
- **Interface Web**: Servidor HTTP simples para upload de arquivos
- **Suporte Docker**: Implantação em contêiner com Tesseract OCR

## Início Rápido

```sh
# Clonar e compilar
git clone <url-repositorio>
cd c6bank-transactions
go build -o bin/c6bank-transactions ./cmd/c6bank-transactions

# Executar o servidor
./bin/c6bank-transactions
# Servidor roda em http://localhost:4500

# Ou com Docker
docker build -t c6bank-transactions .
docker run -p 4500:4500 c6bank-transactions
```

## Uso

Faça upload dos arquivos através da interface web em `http://localhost:4500` ou use curl:

```sh
curl -X POST -F "file=@extrato.pdf" http://localhost:4500/upload
```

## Modelos de iPhone Suportados

| Modelo | Largura | Altura |
|--------|---------|--------|
| iPhone 16 Pro | 1206 | 2622 |
| iPhone 13 Pro Max | 1284 | 2778 |
| iPhone 13 | 1170 | 2532 |

## Desenvolvimento

```sh
# Rodar testes
go test -v ./...

# Rodar com cobertura
go test -v -race -coverprofile=coverage.txt ./...
go tool cover --func coverage.txt
```
