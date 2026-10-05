# Load Tester

CLI em Go, estruturado com [Cobra](https://github.com/spf13/cobra), para executar testes de carga HTTP com concorrência configurável e apresentar um resumo dos resultados.

## Requisitos

- Go 1.22 ou superior para executar localmente.
- Docker para construir e executar a imagem.

## Executar localmente

```powershell
go run ./cmd --url=http://localhost:8080 --requests=1000 --concurrency=10
```

Os parâmetros `--url`, `--requests` e `--concurrency` são obrigatórios. A URL deve usar `http` ou `https`, e requests e concorrência devem ser maiores que zero. O método usado é `GET`; a concorrência efetiva nunca ultrapassa o número total de requests. O comando raiz Cobra também disponibiliza `--help`.

## Construir a imagem Docker

Na raiz do projeto:

```powershell
docker build -t stress .
```

## Executar via Docker

```powershell
docker run --rm stress --url=http://google.com --requests=1000 --concurrency=10
```

Para consultar as opções:

```powershell
docker run --rm stress --help
```

Se o serviço estiver em execução na máquina host, em Docker Desktop para Windows use `host.docker.internal` em vez de `localhost`, por exemplo `--url=http://host.docker.internal:8080`.

## Relatório

Ao concluir, o programa exibe o tempo total, a quantidade exata de requests tentados, a contagem de respostas HTTP 200, a distribuição dos demais status HTTP e as falhas de conexão ou timeout. Respostas HTTP, inclusive erros como 404 e 500, contam como requests realizados. Cada chamada tem timeout de 30 segundos.

 
## Autora
Diane Elo Rieck

## Licença
Projeto para fins educacionais e estudo.