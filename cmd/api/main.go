package main

import (
	"log"
	"net"
	"os"

	userpb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/Ostap00034/course-work-backend-user-service/db"
	"github.com/Ostap00034/course-work-backend-user-service/internal/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	dbString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		log.Fatal("not DB_CONN_STRING in .env file")
	}
	client := db.NewClient(dbString)
	defer client.Close()

	repo := user.NewRepo(client)
	svc := user.NewService(repo)
	srv := user.NewServer(svc)

	lis, _ := net.Listen("tcp", ":50052")
	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, srv)

	log.Println("UserService on :50052")
	log.Fatal(s.Serve(lis))
}
