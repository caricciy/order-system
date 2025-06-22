# Order System
Repositório para implementação do módulo de clean architecture da Pós-graduação em GO.

## Desafio
Implementar os itens abaixo, utilizando os conceitos de Clean Architecture:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
- Criar as migrações necessárias para o banco de dados
- Criar script api.http com a request para criar e listar as orders.

## Rodando o projeto

### Pré-requisitos
- Ter o Go instalado na versão 1.20 ou superior.
- Ter o Docker instalado e rodando.
- Ter o Evans instalado (para testar o GRPC).
- Ter o Make instalado (para facilitar a execução dos comandos).

### Iniciando os containers
Antes de rodar o projeto, é necessário iniciar os containers do Docker. Para isso, execute o seguinte comando na raiz do projeto:
```bash
docker compose up -d
```

### Executando o projeto
Para executar o projeto, é necessário executar as migrações do banco de dados, criar a fila no RabbitMQ e iniciar o serviço. 
Para facilitar esse processo _(executar as migrações, criar a fila e iniciar o serviço)_, o projeto possui a task `start-server` no Makefile.

```bash
make start-server
```
>Antes de executar esse comando, certifique-se de que os containers do Docker estão rodando corretamente.

## Testando o projeto

### Pre-requisitos
- Ter o Evans instalado (para testar o GRPC).
- Ter alguma IDE com a extensão REST Client (para testar a API REST).
- Ter o Make instalado (para facilitar a execução dos comandos). 

### Testes unitários
Para rodar os testes unitários do projeto, execute a tarefa `test` do Makefile:
```bash
make test
```

### Rest API
Para testar a API REST, você pode utilizar o arquivo `api.http` localizado na pasta `api` na raíz do projeto. 
Basta abrir esse arquivo em alguma IDE que suporte REST Client (como o Visual Studio Code com a extensão REST Client) e executar as requisições contidas nele.


### GRPC
Para testar o serviço GRPC, você pode utilizar o Evans, que é uma ferramenta de linha de comando para interagir com serviços GRPC.
Para facilitar o processo utilize o seguinte comando:
```bash
make evans-repl
```

Uma vez dentro do REPL do Evans, você pode chamar os métodos disponíveis no serviço `OrderService`. Por exemplo:
```bash
pb.OrderService@127.0.0.1:50051> call CreateOrder
```
ou
```bash
pb.OrderService@127.0.0.1:50051> call ListOrders
```

### GraphQL
Para testar o serviço GraphQL, você pode utilizar o endereço `http://localhost:8080/graphql` e executar as queries disponíveis.
Abaixo estão alguns exemplos de queries e mutations que você pode executar:

#### Criando uma Order
```graphql
mutation createOrder {
  createOrder(input: { id: "f99a08cc", Price: 1.30, Tax: 0.10 }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

#### Listando Orders
```graphql
query listOrders {
  orders {
    id
    Tax
    Price
    FinalPrice
  }
}
```
## Makefile — Tasks Disponíveis

O projeto possui um conjunto de tasks automatizadas via **Makefile** para auxiliar no desenvolvimento, geração de código e gerenciamento da base de dados e mensageria.


### Variáveis configuráveis (defaults):

| Variável             | Default                                    | Descrição                                  |
|----------------------|--------------------------------------------|---------------------------------------------|
| `DB_PORT`            | `3307`                                     | Porta do banco MySQL                       |
| `DB_USER`            | `root`                                     | Usuário do banco                           |
| `DB_PASSWORD`        | `root`                                     | Senha do banco                             |
| `PROTO_DIR`          | `./internal/infra/grpc/protofiles`         | Diretório dos arquivos `.proto`             |
| `OUT_DIR`            | `./internal/infra/grpc/pb`                 | Diretório de saída dos arquivos gerados     |
| `PROTO_FILES`        | `$(wildcard $(PROTO_DIR)/*.proto)`         | Arquivos `.proto` a serem processados       |
| `RABBITMQ_CONTAINER` | `rabbitmq-3`                               | Nome do container Docker do RabbitMQ        |
| `RABBITMQ_USER`      | `guest`                                    | Usuário do RabbitMQ                        |
| `RABBITMQ_PASSWORD`  | `guest`                                    | Senha do RabbitMQ                          |

### Sobrescrevendo Variáveis nas Tasks do Makefile

As variáveis definidas no `Makefile` possuem valores padrão, mas podem ser sobrescritas diretamente na linha de comando durante a execução de qualquer task.

Para sobrescrever uma variável, utilize a seguinte sintaxe:
```bash
make <task> VARIAVEL=valor
```
Por exemplo, para sobrescrever a porta do banco de dados ao executar uma task, você pode fazer:
```bash
make migrate DB_PORT=3306
```

### Tasks disponíveis:

| Task                  | Descrição                                                                        |
|-----------------------|-----------------------------------------------------------------------------------|
| `make create-migration name=nome_migracao` | Cria uma nova migration SQL sequencial                     |
| `make migrate`        | Executa as migrations no banco de dados                                          |
| `make rollback`       | Desfaz a última migration                                                        |
| `make graphql-gen`    | Gera o código GraphQL via `gqlgen`                                               |
| `make grpc-gen`       | Gera o código Go a partir dos arquivos `.proto`                                  |
| `make grpc-clean`     | Remove os arquivos gRPC gerados                                                  |
| `make evans-repl`     | Abre um REPL interativo com o serviço gRPC via Evans                             |
| `make test`           | Executa todos os testes com verbose                                              |
| `make rabbitmq-queue` | Cria a fila `order.created` e faz o binding no RabbitMQ (via Docker)             |
| `make server-start`   | Executa o servidor com dependências (RabbitMQ e migrations aplicadas)            |

---

### Observações
> Em um projeto real o arquivo .env não deve ser versionado, mas para facilitar o processo de validação do projeto ele foi mantido.

> As `tags json` foram mantidas nos DTOs dos use cases exclusivamente para facilitar o processo de validação do projeto, 
> embora o ideal fosse evitá-las, já que atendem a uma necessidade pontual da camada REST, o que acaba por introduzir um
> acoplamento indesejado entre a lógica de negócio e a representação dos dados para um cliente específico.