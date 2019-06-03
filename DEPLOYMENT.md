
# Single Node

## MongoDB

```
docker run -d \
    --restart always \
    --name mongodb-wt \
    -v $HOME/Programs/mongodb-wt/data:/data/db \
    -p 27017:27017 \
    mongo:4.0
```

Index
```
{tenant_id:1, id:1}
{tenant_id:1, created_at:-1}
```

### Prep

1. Create a database named `engine`

## ArangoDB

### RocksDB

```
docker run -d \
    --restart always \
    --name arangodb-rock \
    -e ARANGO_ROOT_PASSWORD=Test1244 \
    -p 8529:8529 \
    -v $HOME/Programs/arangodb-rock:/var/lib/arangodb3 \
    arangodb:3.4.5
```

### MMFiles

```
docker run -d \
    --restart always \
    --name arangodb-mm \
    -e ARANGO_STORAGE_ENGINE=mmfiles \
    -e ARANGO_ROOT_PASSWORD=Test1244 \
    -p 8539:8529 \
    -v $HOME/Programs/arangodb-mm:/var/lib/arangodb3 \
    arangodb:3.4.5
```

Index
```
{tenant_id, id}
{tenant_id, created_at}
```

### Prep

1. Create a user with name `engine` and password `engine`
1. Create a database named `engine`. set user to `engine`

## Couchbase

```
docker run -d \
    --restart always \
    --name couchbase \
    -p 8091-8094:8091-8094 \
    -p 11210:11210 \
    -v $HOME/Programs/couchbase:/opt/couchbase/var \
    couchbase:community-6.0.0

```

Index
```
CREATE PRIMARY INDEX `default_engine_index` ON `engine` USING GSI;
CREATE INDEX entity_tenant_id ON engine(entity, tenant_id, id);
CREATE INDEX entity_tenant_created_at ON engine(entity, tenant_id, -created_at);
```

### Prep

1. Create bucket named `engine`. Allocate 512GB memory 
1. Create user named `engine` and password `engine`. Set to bucket `engine` 

## Postgres

```
docker run -d \
    --restart always \
    --name postgres-bench \
    -p 5436:5432 \
    -e "POSTGRES_USER=engine" \
    -e "POSTGRES_PASSWORD=engine" \
    -e "POSTGRES_DB=engine" \
    -e "POSTGRES_PORT=5432" \
    -v $HOME/Programs/postgres/data:/var/lib/postgresql/data \
    postgres:9.4
```

```
docker run -d \
    --restart always \
    --name postgres-bench-11 \
    -p 5438:5432 \
    -e "POSTGRES_USER=engine" \
    -e "POSTGRES_PASSWORD=engine" \
    -e "POSTGRES_DB=engine" \
    -e "POSTGRES_PORT=5432" \
    -v $HOME/Programs/postgres11/data:/var/lib/postgresql/data \
    postgres:11.1
```

Migration
```
migrate -database postgres://engine:engine@localhost:5436/engine?sslmode=disable -path migration/postgres up

migrate -database postgres://engine:engine@localhost:5436/engine?sslmode=disable -path migration/postgres down
```

psql
```
psql -h localhost -p 5436 -U engine -W engine engine
```

### Prep

Nothing. initialize will do it

## MySQL

```
docker run -d \
    --restart always \
    --name mysql \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=engine \
    -e MYSQL_USER=engine \
    -e MYSQL_PASSWORD=engine \
    -e MYSQL_DATABASE=engine \
    -v $HOME/Programs/mysql:/var/lib/mysql \
    mysql:8.0
```

Migration
```
migrate -database "mysql://engine:engine@tcp(localhost:3306)/engine" -path migration/mysql up

migrate -database "mysql://engine:engine@tcp(localhost:3306)/engine" -path migration/mysql down
```

### Prep

Nothing. initialize will do it

## cAdvisor

```
docker run -d \
  --restart always \
  --name=cadvisor \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --publish=8080:8080 \
  google/cadvisor:v0.33.0
```

## Prometheus

```
docker run -d  \
    --restart always \
    --publish 9090:9090 \
    --name prometheus \
    --volume $HOME/Programs/prometheus:/prometheus \
    --volume $(pwd)/prometheus:/etc/prometheus \
    prom/prometheus:v2.10.0 --web.enable-lifecycle --config.file=/etc/prometheus/prometheus.yml
```

## Grafana

```
docker run -d \
    --restart always \
    --publish 3000:3000 \
    --name grafana \
    -e "GF_SECURITY_ADMIN_USER=admin" \
    -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
    grafana/grafana:6.2.1
```
