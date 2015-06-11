package etcd

import "skipper/settings"

type EndpointType string

const (
	Ephttp  EndpointType = "http"
	Ephttps EndpointType = "https"
)

type Server struct {
	Url string
}

type Backend struct {
	Typ     EndpointType
	Name    string
	Servers []Server
}

type Frontend struct {
	Typ        EndpointType
	Name       string
	BackendId  string
	Route      string
	Middleware []string
}

type Middleware struct {
	Id       string
	Priority int
	Typ      string // ???
	Args     interface{}
}

type Settings interface {
	GetBackends() map[string]Backend
	GetFrontends() map[string]Frontend
	GetMiddleware() map[string]Middleware
}

type settingsStruct struct {
	backends   map[string]Backend
	frontends  map[string]Frontend
	middleware map[string]Middleware
}

type Client interface {
	GetSettings() <-chan Settings
	Start()
}

type etcdClient struct {
	settings chan Settings
}

func (s *settingsStruct) GetBackends() map[string]Backend {
	return s.backends
}

func (s *settingsStruct) GetFrontends() map[string]Frontend {
	return s.frontends
}

func (s *settingsStruct) GetMiddleware() map[string]Middleware {
	return s.middleware
}

func MakeEtcdClient() Client {
	return &etcdClient{make(chan Settings)}
}

func (ec *etcdClient) feedSettings() {
	testSettings := &settingsStruct{
		backends: map[string]Backend{
			"test": Backend{
				Typ: Ephttp,
				Servers: []Server{
					Server{Url: "http://localhost:9999/slow"}}}},
		frontends:  map[string]Frontend{},
		middleware: map[string]Middleware{}}

	for {
		ec.settings <- testSettings
	}
}

func (ec *etcdClient) Start() {
	go ec.feedSettings()
}

func (ec *etcdClient) GetSettings() <-chan Settings {
	return ec.settings
}

type RawData struct {
	mapping map[string]string
}

type Mock struct {
	RawData *RawData
	get     chan settings.RawData
}

func TempMock() *Mock {
	m := &Mock{
		&RawData{
			map[string]string{
				"Path(\"/hello<v>\")": "http://localhost:9999/slow"}},
		make(chan settings.RawData)}
	go m.feed()
	return m
}

func (m *Mock) feed() {
	for {
		m.get <- m.RawData
	}
}

func (m *Mock) Get() <-chan settings.RawData {
	return m.get
}

func (rd *RawData) GetTestMapping() map[string]string {
	return rd.mapping
}