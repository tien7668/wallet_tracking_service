## Global Environments

| Environment variable | Note                                                                              |
| -------------------- | --------------------------------------------------------------------------------- |
| LOG_LEVEL            | Log level (`trace`,`debug`,`info`,`warn`,`error`,`fatal`,`panic`) default: `info` |
| DATABASE_HOST        | Database host                                                                     |
| DATABASE_PORT        | Database host                                                                     |
| DATABASE_DBNAME      | Database name                                                                     |
| DATABASE_USER        | Database user                                                                     |
| DATABASE_PASSWORD    | Database password                                                                 |

---

## Docker

docker-compose --env-file .env up

## Admin dashboard

localhost:8080 | server = mysql_service_name = mysql | other credentials in .env file

go build -o app ./cmd/app/main.go

DATABASE_HOST="localhost" DATABASE_PORT="3306" DATABASE_DBNAME="kyberswap_user_monitor" DATABASE_USER="admin" DATABASE_PASSWORD="123456" ./app -c internal/pkg/config/config.yaml api

localhost:8081

go build -o app ./cmd/app/main.go

DATABASE_HOST="localhost" DATABASE_PORT="3306" DATABASE_DBNAME="kyberswap_user_monitor" DATABASE_USER="admin" DATABASE_PASSWORD="123456" ./app -c internal/pkg/config/config.yaml fetcher
