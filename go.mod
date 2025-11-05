module star

go 1.25

require (
	github.com/a-h/templ v0.3.960
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

require github.com/google/go-cmp v0.7.0 // indirect

replace github.com/DraconDev/protos/auth-cerberus => ./github.com/DraconDev/protos/auth-cerberus

require (
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)
