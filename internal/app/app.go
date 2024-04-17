package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"

	"github.com/mistandok/chat-server/internal/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegerConfig "github.com/uber/jaeger-client-go/config"

	"github.com/mistandok/chat-server/internal/interceptor"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mistandok/platform_common/pkg/closer"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"

	"github.com/mistandok/chat-server/internal/config"
	desc "github.com/mistandok/chat-server/pkg/chat_v1"
	_ "github.com/mistandok/chat-server/statik" //nolint:golint,unused
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	serviceName = "chat"
)

// App ..
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	prometheusServer *http.Server
	swaggerServer    *http.Server
	configPath       string
}

// NewApp ..
func NewApp(ctx context.Context, configPath string) (*App, error) {
	a := &App{configPath: configPath}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run ..
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	runActions := []struct {
		action func() error
		errMsg string
	}{
		{action: a.runGRPCServer, errMsg: "ошибка при запуске GRPC сервера"},
		{action: a.runHTTPServer, errMsg: "ошибка при запуске HTTP сервера"},
		{action: a.runSwaggerServer, errMsg: "ошибка при запуске Swagger сервера"},
		{action: a.runPrometheusServer, errMsg: "ошибка при запуске Prometheus сервера"},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(runActions))

	for _, runAction := range runActions {
		currentRunAction := runAction
		go func() {
			defer wg.Done()

			err := currentRunAction.action()
			if err != nil {
				log.Fatalf(currentRunAction.errMsg)
			}
		}()
	}

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	initDepFunctions := []func(context.Context) error{
		a.initConfig,
		a.initGlobalTracer,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initMetrics,
		a.initPrometheusServer,
		a.initSwaggerServer,
	}

	for _, f := range initDepFunctions {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(a.configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
			interceptor.ServerTracingInterceptor,
			interceptor.MetricsInterceptor,
			interceptor.ValidateInterceptor,
			a.serviceProvider.AccessCheckInterceptor(ctx).Get,
		),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterChatV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 2 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/chat_api.swagger.json", serveSwaggerFile("/chat_api.swagger.json"))
	mux.HandleFunc("/auth_api.swagger.json", serveSwaggerFile("/auth_api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              a.serviceProvider.PrometheusConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return nil
}

func (a *App) initGlobalTracer(_ context.Context) error {
	cfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	return err
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC запущен на %s", a.serviceProvider.GRPCConfig().Address())
	listener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP запущен на %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger сервер запущен на %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

func (a *App) runPrometheusServer() error {
	log.Printf("Prometheus запущен на %s", a.serviceProvider.PrometheusConfig().Address())
	log.Printf("Prometheus доступен по адресу %s", a.serviceProvider.PrometheusConfig().PublicAddress())

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
