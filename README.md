# JWT Authentication with Golang, Gin, and GORM.

## Quick Start

To start, clone this project:

```bash
git clone git@github.com:yaolongt/jwt-auth-go.git
```

## Development

To deploy locally for development, you will need the following dependencies:

- Docker / Orbstack

Fill up the `.env` file with environment variables and generate public/private keys:

```bash
openssl genrsa -out private.pem 4096
openssl rsa -in private.pem -pubout -out public.pem
```

To run:

```bash
docker compose up -d
```
