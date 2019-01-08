export PATH=$PATH:/home/$(USER)/go/bin
export GOROOT=/home/$(USER)/go
export GOPATH=$(GOROOT)/src

go get github.com/lib/pq
go get github.com/dgrijalva/jwt-go
go get github.com/SermoDigital/jose
go get github.com/Azure/go-ntlmssp
go get github.com/julienschmidt/httprouter
go get github.com/metakeule/fmtdate
go get github.com/gorilla/sessions
go get github.com/buger/jsonparser
go get github.com/rs/cors