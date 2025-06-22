package main

import (
	"database/sql"
	"fmt"
	"github.com/caricciy/ordersystem/internal/event"
	"github.com/caricciy/ordersystem/internal/infra/graph"
	"log"
	"net"
	"net/http"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caricciy/ordersystem/configs"
	"github.com/caricciy/ordersystem/internal/event/handler"
	"github.com/caricciy/ordersystem/internal/infra/grpc/pb"
	"github.com/caricciy/ordersystem/internal/infra/grpc/service"
	"github.com/caricciy/ordersystem/internal/infra/web/webserver"
	"github.com/caricciy/ordersystem/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(config.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	rabbitMQChannel := getRabbitMQChannel(*config)

	eventDispatcher := events.NewEventDispatcher()
	_ = eventDispatcher.Register(event.GetOrderCreatedEventName(), handler.NewOrderCreatedHandler(rabbitMQChannel))

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrdersUseCase(db)

	webServer := webserver.NewWebServer(config.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webServer.AddHandler("/order", webOrderHandler.Create)
	webServer.AddHandler("/orders", webOrderHandler.List)
	fmt.Println("Starting web server on port", config.WebServerPort)
	go webServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	srv := graphqlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	log.Fatal(http.ListenAndServe(":"+config.GraphQLServerPort, nil))
}

func getRabbitMQChannel(config configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQUser, config.RabbitMQPassword, config.RabbitMQHost, config.RabbitMQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
