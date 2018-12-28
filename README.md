# avdb

# Ansible Variables Database

# Build binary
> go get -u github.com/ecray/avdb

> GOOS=linux go build -o release/avdb .
> docker build -t avdb-dev .
> docker-compose up

## docker-compose should output initial token

> curl -s -X POST -H "Auth-Token: $AVDB_TOKEN" http://127.0.0.1:3333/api/v1/hosts/tacotruck01 -d '{"comidas":"tacos"}' | jq '.'
> curl -s -X GET http://127.0.0.1:3333/api/v1/hosts | jq '.'

# References
```
http://www.golangprograms.com/advance-programs/golang-restful-api-using-grom-and-gorilla-mux.html
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand
```
