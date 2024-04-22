package service

import (
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/orgs"
	"github.com/les-cours/learning-service/api/users"
	"github.com/les-cours/learning-service/env"
	"github.com/les-cours/learning-service/resolvers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"runtime"
)

var (
	registry       = prometheus.NewRegistry()
	requestCounter = prometheus.NewGauge(prometheus.GaugeOpts{Name: "request_counter", Help: "request counter"})
	memoryUsage    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "memory_usage", Help: "memory usage"})
	goRoutineNum   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_routines_num", Help: "the number of go routine "})
)

func monitoring_middleware(originalHandler http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memoryUsage.Set(float64(m.Alloc))
		goRoutineNum.Set(float64(runtime.NumGoroutine()))
		requestCounter.Inc()
		originalHandler.ServeHTTP(w, r)
	})
}

func Start() {
	registry.MustRegister(requestCounter, memoryUsage, goRoutineNum)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.HandleFunc("/metrics", monitoring_middleware(promHandler))
	log.Printf("Starting http server on port " + env.Settings.HttpPort)
	go func() {
		err := http.ListenAndServe(":"+env.Settings.HttpPort, nil)
		if err != nil {
			log.Fatalf("Error http server on port %v: %v", env.Settings.HttpPort, err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+env.Settings.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", env.Settings.GrpcPort, err)
	}

	db, err := StartDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()
	usersConnectionService, err := grpc.Dial(env.Settings.UserService.Host+":"+env.Settings.UserService.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to users service %v", err)
	}

	defer usersConnectionService.Close()

	defer db.Close()
	orgsConnectionService, err := grpc.Dial(env.Settings.OrgService.Host+":"+env.Settings.OrgService.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to orgs service %v", err)
	}

	defer orgsConnectionService.Close()

	server := resolvers.Server{
		DB:    db,
		Users: users.NewUserServiceClient(usersConnectionService),
		Orgs:  orgs.NewOrgServiceClient(orgsConnectionService),
	}
	grpcServer := grpc.NewServer()
	learning.RegisterLearningServiceServer(grpcServer, &server)
	log.Printf("Starting grpc server on port " + env.Settings.GrpcPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start gRPC server on port %v: %v", env.Settings.GrpcPort, err)
	}
}
