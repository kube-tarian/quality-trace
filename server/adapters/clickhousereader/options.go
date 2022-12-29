package clickhousereader

import (
	"context"
	"fmt"
	"net/url"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

type Encoding string

const (
	// EncodingJSON is used for spans encoded as JSON.
	EncodingJSON Encoding = "json"
	// EncodingProto is used for spans encoded as Protobuf.
	EncodingProto Encoding = "protobuf"
)

const (
	defaultDatasource             string        = "tcp://localhost:9000"
	defaultTraceDB                string        = "signoz_traces"
	defaultSpansTable             string        = "signoz_spans"
	defaultLiveTailRefreshSeconds int           = 10
	defaultWriteBatchDelay        time.Duration = 5 * time.Second
	defaultWriteBatchSize         int           = 10000
	defaultEncoding               Encoding      = EncodingJSON
)

// NamespaceConfig is Clickhouse's internal configuration data
type namespaceConfig struct {
	namespace               string
	Enabled                 bool
	Datasource              string
	TraceDB                 string
	OperationsTable         string
	IndexTable              string
	DurationTable           string
	UsageExplorerTable      string
	SpansTable              string
	ErrorTable              string
	DependencyGraphTable    string
	TopLevelOperationsTable string
	LogsDB                  string
	LogsTable               string
	LogsAttributeKeysTable  string
	LogsResourceKeysTable   string
	LiveTailRefreshSeconds  int
	WriteBatchDelay         time.Duration
	WriteBatchSize          int
	Encoding                Encoding
	Connector               Connector
}

// Connecto defines how to connect to the database
type Connector func(cfg *namespaceConfig) (clickhouse.Conn, error)

func defaultConnector(cfg *namespaceConfig) (clickhouse.Conn, error) {
	ctx := context.Background()
	dsnURL, err := url.Parse(cfg.Datasource)
	fmt.Println("URL:", dsnURL)
	options := &clickhouse.Options{
		Addr: []string{dsnURL.Host},
	}
	if dsnURL.Query().Get("username") != "" {
		auth := clickhouse.Auth{
			Username: "admin",
			Password: "admin",
		}
		options.Auth = auth
	}

	fmt.Printf("Options:%v\n", options)

	db, err := clickhouse.Open(options)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// Options store storage plugin related configs
type Options struct {
	primary *namespaceConfig

	others map[string]*namespaceConfig
}

// NewOptions creates a new Options struct.
func NewOptions(datasource string, primaryNamespace string, otherNamespaces ...string) *Options {

	if datasource == "" {
		datasource = defaultDatasource
	}

	options := &Options{
		primary: &namespaceConfig{
			namespace:              primaryNamespace,
			Enabled:                true,
			Datasource:             datasource,
			TraceDB:                defaultTraceDB,
			SpansTable:             defaultSpansTable,
			LiveTailRefreshSeconds: defaultLiveTailRefreshSeconds,
			WriteBatchDelay:        defaultWriteBatchDelay,
			WriteBatchSize:         defaultWriteBatchSize,
			Encoding:               defaultEncoding,
			Connector:              defaultConnector,
		},
		others: make(map[string]*namespaceConfig, len(otherNamespaces)),
	}

	for _, namespace := range otherNamespaces {
		if namespace == archiveNamespace {
			options.others[namespace] = &namespaceConfig{
				namespace:              namespace,
				Datasource:             datasource,
				TraceDB:                "",
				OperationsTable:        "",
				IndexTable:             "",
				ErrorTable:             "",
				LogsDB:                 "",
				LogsTable:              "",
				LogsAttributeKeysTable: "",
				LogsResourceKeysTable:  "",
				LiveTailRefreshSeconds: defaultLiveTailRefreshSeconds,
				WriteBatchDelay:        defaultWriteBatchDelay,
				WriteBatchSize:         defaultWriteBatchSize,
				Encoding:               defaultEncoding,
				Connector:              defaultConnector,
			}
		} else {
			options.others[namespace] = &namespaceConfig{namespace: namespace}
		}
	}

	return options
}

// GetPrimary returns the primary namespace configuration
func (opt *Options) getPrimary() *namespaceConfig {
	return opt.primary
}
