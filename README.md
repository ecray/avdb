# avdb

# Ansible Variables Database

## Build binary
```
$ make build
$ ls -l release/avdb
```
## Build binary, container and run in docker-compose for dev environment
```
$ make dev-env
...
db_1    | LOG:  database system is ready to accept connections
db_1    | LOG:  autovacuum launcher started
avdb_1  | 2018/12/28 19:56:09 Auto-migrating schema...
avdb_1  | 2018/12/28 19:56:09 Checking token credential..
avdb_1  | 2018/12/28 19:56:09 Initial token: NXSABjq0pDVENA9wlVfY5GF3BLJ1-3lwUG95QaPquk4=
avdb_1  | 2018/12/28 19:56:09 Running on 0.0.0.0:3333
```

## Use token reported from docker-compose logs 
```
$ export AVDB_TOKEN=NXSABjq0pDVENA9wlVfY5GF3BLJ1-3lwUG95QaPquk4=
$ curl -s -X POST -H "Auth-Token: $AVDB_TOKEN" http://127.0.0.1:3333/api/v1/hosts/tacotruck01 -d '{"comidas":"tacos"}' | jq '.'
$ curl -s -X GET http://127.0.0.1:3333/api/v1/hosts | jq '.'
```

# References
```
http://www.golangprograms.com/advance-programs/golang-restful-api-using-grom-and-gorilla-mux.html
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand
```
