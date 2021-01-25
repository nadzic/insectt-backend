module insectt

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0 // indirect
	insectt.io/api v0.0.0-00010101000000-000000000000
)

replace insectt.io/api => ./api/

replace insectt.io/models => ./models/
