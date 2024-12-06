package kimsufi

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	kimsufiauthentication "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/authentication"
	kimsufiavailability "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/availability"
	kimsuficatalog "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/catalog"
)

// MultiService is a map of OVH endpoints to Services.
type MultiService map[string]Service

// Service is a wrapper around ovh.Client
// with optional caching and logging.
type Service struct {
	cache  *cache.Cache
	client *ovh.Client
	logger *Logger
}

// NewMultiService creates a new MultiService
// with a Service for each OVH endpoint.
func NewMultiService(l *logrus.Logger, c *cache.Cache) (MultiService, error) {
	m := make(MultiService, 0)

	for _, endpoint := range GetOVHEndpoints() {
		s, err := NewService(endpoint, l, c)
		if err != nil {
			log.Errorf("failed to create OVH client for %s: %v", endpoint, err)
			return nil, err
		}

		m[endpoint] = *s
	}

	return m, nil
}

// Endpoint returns the Service for the given endpoint.
func (m MultiService) Endpoint(endpoint string) *Service {
	s, found := m[endpoint]
	if !found {
		return nil
	}

	return &s
}

// NewService creates a new Service for the given endpoint.
// logger is optional, if nil a no-op logger will be used.
// c is optional, if nil no caching will be used.
func NewService(endpoint string, logger *logrus.Logger, c *cache.Cache) (*Service, error) {
	e, found := ovh.Endpoints[endpoint]
	if !found {
		return nil, fmt.Errorf("invalid endpoint %s", endpoint)
	}

	client, err := ovh.NewClient(e, "none", "none", "none")
	if err != nil {
		return nil, err
	}

	s := &Service{
		cache:  c,
		client: client,
		logger: NewRequestLogger(logger),
	}

	client.Logger = s.logger

	return s, nil
}

// GetOVHEndpoints returns a list of OVH endpoints.
// It keeps only the ones starting with "ovh-".
func GetOVHEndpoints() []string {
	endpoints := slices.Sorted(maps.Keys(ovh.Endpoints))
	endpoints = slices.DeleteFunc(endpoints, func(s string) bool {
		return !strings.HasPrefix(s, "ovh-")
	})

	return endpoints
}

// GetAvailabilities returns servers availabilities.
// datacenters is a list of datacenters to filter on.
// planCode is the plan code to filter on.
// options is a map of additional query parameters.
// see https://eu.api.ovh.com/console/?section=%2Fdedicated%2Fserver&branch=v1#get-/dedicated/server/datacenter/availabilities
func (s *Service) GetAvailabilities(datacenters []string, planCode string, options map[string]string) (*kimsufiavailability.Availabilities, error) {
	path := "/dedicated/server/datacenter/availabilities"

	queryArgs := make(map[string]string)
	if len(datacenters) > 0 {
		queryArgs["datacenters"] = strings.Join(datacenters, ",")
	}
	if planCode != "" {
		queryArgs["planCode"] = planCode
	}
	for key, value := range options {
		queryArgs[key] = value
	}

	var availabilities *kimsufiavailability.Availabilities
	err := s.request(http.MethodGet, path, queryArgs, nil, &availabilities, false)
	if err != nil {
		return nil, err
	}

	return availabilities, nil
}

// ListServers returns the catalog for the given OVH subsidiary.
// ovhSubsidiary is the country code to filter on, given in a two-letter format.
// see https://eu.api.ovh.com/console/?section=%2Forder&branch=v1#get-/order/catalog/public/eco
func (s *Service) ListServers(ovhSubsidiary string) (*kimsuficatalog.Catalog, error) {
	path := "/order/catalog/public/eco"

	queryArgs := map[string]string{
		"ovhSubsidiary": ovhSubsidiary,
	}

	var catalog *kimsuficatalog.Catalog
	err := s.request(http.MethodGet, path, queryArgs, nil, &catalog, false)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

// GetAuthDetails performs a test API request to check if the client is authenticated.
func (s *Service) GetAuthDetails() error {
	path := "/auth/details"

	err := s.client.Get(path, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCurrentCredential() (*kimsufiauthentication.CurrentCredentialResponse, error) {
	path := "/auth/currentCredential"

	var resp kimsufiauthentication.CurrentCredentialResponse
	err := s.client.Get(path, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// WithAuth returns a new authenticated Service with the given credentials.
func (s *Service) WithAuth(appKey, appSecret, consumerKey string) (*Service, error) {
	authClient, err := ovh.NewClient(s.client.Endpoint(), appKey, appSecret, consumerKey)
	if err != nil {
		return nil, err
	}

	newService := &Service{
		cache:  s.cache,
		logger: s.logger,
		client: authClient,
	}

	return newService, nil
}

// request performs an API request.
// this is a wrapper around ovh.Client.CallAPI, it allows for caching when set on the Service.
// path and queryArgs are combined to form the request URL.
// method, body, response, and needAuth are passed as is.
// response must be a pointer.
func (s *Service) request(method, path string, queryArgs map[string]string, body any, response any, needAuth bool) error {
	// Ensure response is a pointer
	rv := reflect.ValueOf(response)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("response must be a pointer")
	}

	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	q := u.Query()
	for key, value := range queryArgs {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	cacheKey := fmt.Sprintf("%s%s", s.client.Endpoint(), u.String())

	var found bool = false
	var cacheEntry interface{}
	if s.cache != nil {
		cacheEntry, found = s.cache.Get(cacheKey)
	}
	if found {
		s.logger.Tracef("cache hit: %s", cacheKey)
		ce := reflect.ValueOf(cacheEntry)
		rv.Elem().Set(ce.Elem())
	} else {
		s.logger.Tracef("cache miss: %s", cacheKey)
		err = s.client.CallAPI(method, u.String(), body, response, needAuth)
		if err != nil {
			return err
		}
		if s.cache != nil {
			s.cache.Set(cacheKey, response, cache.DefaultExpiration)
		}
	}

	return nil
}
