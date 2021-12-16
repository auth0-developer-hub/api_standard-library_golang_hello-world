set -xe

#make mod-down
make mod-tidy
make fmt
make swagger
make vet
make test
make hello-golang-api

./hello-golang-api

