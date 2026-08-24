package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Knoppiix/InstagramEmbedResolver/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
var appConfig Config

var routes = []string{
	"/p/{id}",
	"/p/{id}/{$}",
	"/reels/{id}",
	"/reels/{id}/{$}",
	"/reel/{id}",
	"/reel/{id}/{$}",
	"/{username}/p/{id}",
	"/{username}/p/{id}/{$}",
	"/{username}/reel/{id}",
	"/{username}/reel/{id}/{$}",
	"/share/{id}",
	"/share/{id}/{$}",
}
var defaultRes Resolver

func startServer(resolvers []Resolver, cfg Config) {
	template := template.Must(template.ParseFiles("templates/index.html"))
	// home page handler
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, cfg.templateData())
	})

	for _, route := range routes {
		mux.HandleFunc("GET "+route, reqHandler(resolvers))
	}
	// expose the prometheus metrics route
	mux.Handle("GET /metrics", promhttp.Handler())

	// stats endpoints backing the dashboard charts on "/", proxying pre-set PromQL queries to Prometheus
	mux.HandleFunc("GET /api/stats/success", statsSuccessHandler)
	mux.HandleFunc("GET /api/stats/latency", statsLatencyHandler)
	mux.HandleFunc("GET /api/stats/requests-timeseries", statsRequestsTimeseriesHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux))
}

func reqHandler(resolvers []Resolver) http.HandlerFunc {
	// handler function for proxifying the requests to the resolvers
	return func(w http.ResponseWriter, r *http.Request) {

		// If the request is from discord OR telegram, we proxify the request through the resolver
		ua := r.Header.Get("User-Agent")
		if strings.Contains(ua, "Discordbot") || strings.Contains(ua, "Telegram") {
			proxyRequest(w, r, resolvers)
			return
		}
		// Else, we simply redirect the user to the instagram post

		// Extract the img_index query parameter if it exists
		imgIndex := r.URL.Query().Get("img_index")
		var query string
		if imgIndex != "" {
			query = fmt.Sprintf("?img_index=%s", imgIndex)
		}
		// Construct the redirect URL
		redirectURL := "https://instagram.com" + r.URL.Path + query

		// Send HTTP 302 redirect
		http.Redirect(w, r, redirectURL, http.StatusFound)
		log.Printf("User redirection toward %s", redirectURL)
	}
}

func main() {
	port := flag.Int("p", envInt(8080, "PROXY_PORT", "PORT"), "port to run the server on")
	proxyDomain := flag.String("proxy-domain", envHost("PROXY_BASE_DOMAIN", defaultProxyBaseDomain), "base domain used by this proxy")
	normalSubdomain := flag.String("normal-subdomain", normalizeSubdomain(envString("PROXY_NORMAL_SUBDOMAIN", "n")), "subdomain used for normal embeds")
	gallerySubdomain := flag.String("gallery-subdomain", normalizeSubdomain(envString("PROXY_GALLERY_SUBDOMAIN", "g")), "subdomain used for gallery embeds")
	directSubdomain := flag.String("direct-subdomain", normalizeSubdomain(envString("PROXY_DIRECT_SUBDOMAIN", "d")), "subdomain used for direct embeds")
	resolversFile := flag.String("resolvers-file", envString("RESOLVERS_FILE", "resolvers.json"), "path to the resolvers configuration file")
	flag.Parse()
	log.SetOutput(os.Stdout)
	appConfig = Config{
		Port:             *port,
		ProxyBaseDomain:  normalizeHost(*proxyDomain),
		NormalSubdomain:  normalizeSubdomain(*normalSubdomain),
		GallerySubdomain: normalizeSubdomain(*gallerySubdomain),
		DirectSubdomain:  normalizeSubdomain(*directSubdomain),
		ResolversFile:    *resolversFile,
	}
	resolvers, err := loadResolvers(appConfig.ResolversFile)
	if err != nil {
		log.Fatalf("Error reading the file: %v", err)
	}
	applyResolverConfig(resolvers, appConfig)
	metrics.Init()
	go monitorResolvers(resolvers)
	startServer(resolvers, appConfig)
}
