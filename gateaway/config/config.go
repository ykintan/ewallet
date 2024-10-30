// config/config.go
package config

const (
	UserAddress        = "localhost:50052"
	TransactionAddress = "localhost:50051"
	HTTPPort           = ":8080"
	BasicAuthUsername  = "admin"
	BasicAuthPassword  = "password"
)

func GetUserAddress() string {
	return UserAddress
}

func GetTransactionAddress() string {
	return TransactionAddress
}

func GetHTTPPort() string {
	return HTTPPort
}

func GetBasicAuthUsername() string {
	return BasicAuthUsername
}

func GetBasicAuthPassword() string {
	return BasicAuthPassword
}
