package main

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

const defaultProxyBaseDomain = "zzinstagram.com"

type Config struct {
	Port             int
	ProxyBaseDomain  string
	NormalSubdomain  string
	GallerySubdomain string
	DirectSubdomain  string
	PrometheusURL    string
	ResolversFile    string
}

type homePageData struct {
	BaseDomain       string
	ExampleURL       string
	NormalHost       string
	GalleryHost      string
	DirectHost       string
	StatsEnabled     bool
	NormalSubdomain  string
	GallerySubdomain string
	DirectSubdomain  string
}

func envHost(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return normalizeHost(value)
}

func envString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envOptionalString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return strings.TrimSpace(value)
}

func envInt(fallback int, keys ...string) int {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func normalizeHost(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return value
	}
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	if parsed, err := url.Parse(value); err == nil && parsed.Host != "" {
		value = parsed.Hostname()
	}
	return strings.TrimSuffix(value, ".")
}

func normalizeSubdomain(value string) string {
	return strings.Trim(strings.ToLower(strings.TrimSpace(value)), ".")
}

func (c Config) hostFor(subdomain string) string {
	subdomain = normalizeSubdomain(subdomain)
	if subdomain == "" {
		return c.ProxyBaseDomain
	}
	return subdomain + "." + c.ProxyBaseDomain
}

func (c Config) templateData() homePageData {
	return homePageData{
		BaseDomain:       c.ProxyBaseDomain,
		ExampleURL:       "https://" + c.hostFor("") + "/p/{id}",
		NormalHost:       c.hostFor(c.NormalSubdomain),
		GalleryHost:      c.hostFor(c.GallerySubdomain),
		DirectHost:       c.hostFor(c.DirectSubdomain),
		StatsEnabled:     c.statsEnabled(),
		NormalSubdomain:  c.NormalSubdomain,
		GallerySubdomain: c.GallerySubdomain,
		DirectSubdomain:  c.DirectSubdomain,
	}
}

func (c Config) statsEnabled() bool {
	return strings.TrimSpace(c.PrometheusURL) != ""
}

func applyResolverConfig(resolvers []Resolver, cfg Config) {
	for i := range resolvers {
		resolvers[i].Url = applyResolverTokens(resolvers[i].Url, cfg)
		applyRenderModeConfig(resolvers[i].Normal, cfg)
		applyRenderModeConfig(resolvers[i].Gallery, cfg)
		applyRenderModeConfig(resolvers[i].Direct, cfg)
	}
}

func applyRenderModeConfig(mode *RenderMode, cfg Config) {
	if mode == nil {
		return
	}
	mode.Subdomain = applyResolverTokens(mode.Subdomain, cfg)
}

func applyResolverTokens(value string, cfg Config) string {
	replacements := map[string]string{
		"{BASE_DOMAIN}":       cfg.ProxyBaseDomain,
		"{NORMAL_SUBDOMAIN}":  cfg.NormalSubdomain,
		"{GALLERY_SUBDOMAIN}": cfg.GallerySubdomain,
		"{DIRECT_SUBDOMAIN}":  cfg.DirectSubdomain,
	}
	for token, replacement := range replacements {
		value = strings.ReplaceAll(value, token, replacement)
	}
	return value
}
