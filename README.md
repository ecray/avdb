# avdb-go

# Ansible Variables Database - Go Edition

# Build binary
> go get -u github.com/ecray/avdb-go

# Setup environment variables for backend connection
> echo "export DB_HOST=127.0.0.1 \\
        export DB_PORT=5432 \\
        export DB_NAME=avdb \\
        export DB_USER=avdb \\
        export DB_PASS=avdb" >> env.sh

> source env.sh

> ./avdb

# References
http://www.golangprograms.com/advance-programs/golang-restful-api-using-grom-and-gorilla-mux.html
https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand
