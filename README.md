# Weather API - CEP to Temperature Converter

Uma REST API em Go que recebe um CEP brasileiro, identifica a localização e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## Arquitetura

Este projeto segue os princípios da Clean Architecture com a seguinte estrutura:

```
.
├── domain/          # Entidades de domínio e regras de negócio
├── usecase/         # Casos de uso da aplicação
├── adapter/         # Adaptadores de infraestrutura
│   └── http/       # Adaptadores de entrada HTTP (handlers Gin)
└── cmd/            # Ponto de entrada da aplicação
```

## Funcionalidades

- Validação de CEP (8 dígitos)
- Busca de localização via ViaCEP API
- Recuperação de temperatura via WeatherAPI
- Conversão de temperatura (Celsius, Fahrenheit, Kelvin)
- REST API com framework Gin
- Design Clean Architecture
- 82% de cobertura de testes
- Suporte a Docker
- Documentação Swagger API

## Documentação da API

### Swagger UI

Ao executar a aplicação, a documentação Swagger fica disponível em:

- **Desenvolvimento local**: http://localhost:8080/swagger/index.html
- **Docker**: http://localhost:8080/swagger/index.html

O Swagger UI fornece uma interface interativa para explorar e testar os endpoints da API.

## Endpoints da API

### GET /weather/:zipcode

Retorna dados de temperatura para o CEP brasileiro fornecido.

**Resposta de Sucesso (200)**

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

**Respostas de Erro**

- `422 Unprocessable Entity` - Formato de CEP inválido

```json
{
  "message": "invalid zipcode"
}
```

- `404 Not Found` - CEP não encontrado

```json
{
  "message": "can not find zipcode"
}
```

## Requisitos

- Go 1.25+
- Docker & Docker Compose
- Chave WeatherAPI (tier gratuito disponível em [weatherapi.com](https://www.weatherapi.com/))

## Primeiros Passos

### 1. Obtenha uma Chave WeatherAPI

Cadastre-se em [weatherapi.com](https://www.weatherapi.com/) para obter uma chave de API gratuita.

### 2. Configure as Variáveis de Ambiente

Crie um arquivo `.env` baseado no `.env.example`:

```bash
cp .env.example .env
```

Edite o `.env` e adicione sua chave WeatherAPI:

```
WEATHER_API_KEY=your_actual_api_key_here
PORT=8080
```

### 3. Execute com Docker Compose

```bash
docker-compose up --build
```

API disponível em `http://localhost:8080`

Swagger disponível em `http://localhost:8080/swagger/index.html`

### 4. Teste a API

```bash
# CEP válido (Avenida Paulista, São Paulo)
curl http://localhost:8080/weather/01310100

# Formato inválido
curl http://localhost:8080/weather/123

# CEP não encontrado
curl http://localhost:8080/weather/99999999
```

## Desenvolvimento

### Executar Localmente

```bash
# Instalar dependências
go mod download

# Configurar variáveis de ambiente
export WEATHER_API_KEY=your_key_here

# Executar a aplicação
go run cmd/main.go
```

### Executar Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test ./... -coverprofile=coverage.txt

# Visualizar relatório de cobertura
go tool cover -html=coverage.txt
```

### Build

```bash
go build -o weather-api cmd/main.go
```

## Licença

Este projeto faz parte de um exercício de aprendizado para o programa Full Cycle Pós-Graduação GO Expert.
