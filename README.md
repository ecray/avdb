# avdb

# Ansible Variables Database

## Build binary
```
$ make build
$ ls -l release/avdb
```
## Build container and run in docker compose
```
$ make dev-env
```

## docker-compose should output initial token
```
$ curl -s -X POST -H "Auth-Token: $AVDB_TOKEN" http://127.0.0.1:3333/api/v1/hosts/tacotruck01 -d '{"comidas":"tacos"}' | jq '.'
$ curl -s -X GET http://127.0.0.1:3333/api/v1/hosts | jq '.'
```

# References
```
http://www.golangprograms.com/advance-programs/golang-restful-api-using-grom-and-gorilla-mux.html
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand
```
