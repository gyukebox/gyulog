go get github.com/russross/blackfriday
go get github.com/go-sql-driver/mysql
go get github.com/gyukebox/gyulog/post
go get github.com/gyukebox/gyulog/admin
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bin/application .
