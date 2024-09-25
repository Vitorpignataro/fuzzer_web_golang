# Fuzzer Web Golang

Fuzzer Web Golang é uma ferramenta de linha de comando para enumerar diretórios web usando técnicas de fuzzing. Ela permite fazer requisições HTTP simultâneas para um host, substituindo a palavra-chave `FUZZER` pelo nome de diretórios listados em um arquivo de texto, e filtrar os resultados por códigos de status HTTP e tamanhos de conteúdo específicos.

## Funcionalidades

- Suporte a múltiplas threads para execução concorrente de requisições.
- Filtragem de resultados com base em códigos de status HTTP e tamanho do conteúdo.
- Exibição colorida e organizada dos resultados no terminal.

## Como Buildar

Antes de tudo, certifique-se de ter o Go instalado em sua máquina. Você pode verificar se o Go está instalado usando o comando:

```bash
go version
```

## Para buildar a aplicação, siga os passos abaixo:

1) Clone o repositório:
```bash
git clone git@github.com:Vitorpignataro/fuzzer_web_golang.git
cd fuzzer_web_golang
```

2) Compile o projeto:
```bash
go build -o fuzzer_web
```
Isso gerará um binário chamado `fuzzer_web` na pasta raiz do projeto.


# Como Usar

Você pode usar o binário gerado ou executar diretamente o código-fonte com go run.

## Executando a Aplicação
Usando o binário compilado:
```bash
./fuzzer_web --threads 40 --host https://google.com.br/FUZZER --file arquivo.txt --hdc 404,301
```

Usando `go run`:
```bash
go run main.go --threads 40 --host https://google.com.br/FUZZER --file arquivo.txt --hdc 404,301
```

## Parâmetros

- `--threads` -> Define o número de threads a serem usadas para as requisições simultâneas (padrão: 1).
- `--host` -> URL do site alvo, com a palavra-chave FUZZER sendo substituída pelos caminhos no arquivo de entrada.
- `--file` -> Caminho para o arquivo de texto contendo os nomes dos diretórios a serem testados.
- `--hdc` -> Códigos de status HTTP a serem ocultados (opcional).
- `--hcl` -> Content lenght a serem ocultados (opcional).

## Exemplo de Uso
```bash
go run main.go --threads 40 --host https://google.com.br/FUZZER --file arquivo.txt --hdc 404,301 --hcl 5123,2211
```

Nesse exemplo, a aplicação fará requisições simultâneas (usando 40 threads) para o host `https://google.com.br/`, substituindo a palavra `FUZZER` pelos diretórios listados no arquivo `arquivo.txt`. Os resultados com códigos de status 404 e 301 e tamanho de resposta 5123 e 2211 serão ocultados.

