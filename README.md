## Generation
```
alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
swagger generate model -f ne-v1-catalog-ne_v1_PATCHED.yml -t internal -m api
```
