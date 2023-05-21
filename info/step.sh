# install gin
$ go get -u github.com/gin-gonic/gin

# install gorm
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# buat file migrasi
migrate create -ext sql -dir db/migrations create_"namatable"

# migrasi data
migrate -path db/migrations -database "mysql://root:@tcp(127.0.0.1:3306)/go_foundation?charset=utf8mb4&parseTime=True&loc=Local" up

