package clickhousereader

import (
	"context"
	"encoding/json"

	"fmt"
	"os"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kube-tarian/quality-trace/server/model"
)

const (
	primaryNamespace      = "clickhouse"
	archiveNamespace      = "clickhouse-archive"
	signozTraceDBName     = "signoz_traces"
	signozSpansTable      = "signoz_spans"
	signozErrorIndexTable = "signoz_error_index_v2"
	signozTraceTableName  = "signoz_index_v2"
)

// SpanWriter for reading spans from ClickHouse
type ClickHouseReader struct {
	db clickhouse.Conn

	traceDB                 string
	operationsTable         string
	durationTable           string
	indexTable              string
	errorTable              string
	usageExplorerTable      string
	spansTable              string
	dependencyGraphTable    string
	topLevelOperationsTable string
	logsDB                  string
	logsTable               string
	logsAttributeKeys       string
	logsResourceKeys        string

	promConfigFile string

	liveTailRefreshSeconds int
}

// NewTraceReader returns a TraceReader for the database
func NewReader(clickHouseUrl string) *ClickHouseReader {
	datasource := clickHouseUrl
	options := NewOptions(datasource, primaryNamespace, archiveNamespace)
	db, err := initialize(options)

	if err != nil {
		fmt.Printf("failed to initialize ClickHouse: %v", err)
		os.Exit(1)
	}

	return &ClickHouseReader{
		db:                      db,
		traceDB:                 options.primary.TraceDB,
		operationsTable:         options.primary.OperationsTable,
		indexTable:              options.primary.IndexTable,
		errorTable:              options.primary.ErrorTable,
		usageExplorerTable:      options.primary.UsageExplorerTable,
		durationTable:           options.primary.DurationTable,
		spansTable:              options.primary.SpansTable,
		dependencyGraphTable:    options.primary.DependencyGraphTable,
		topLevelOperationsTable: options.primary.TopLevelOperationsTable,
		logsDB:                  options.primary.LogsDB,
		logsTable:               options.primary.LogsTable,
		logsAttributeKeys:       options.primary.LogsAttributeKeysTable,
		logsResourceKeys:        options.primary.LogsResourceKeysTable,
		liveTailRefreshSeconds:  options.primary.LiveTailRefreshSeconds,
	}
}

func initialize(options *Options) (clickhouse.Conn, error) {
	db, err := connect(options.getPrimary())
	if err != nil {
		return nil, fmt.Errorf("error connecting to primary db: %v", err)
	}

	return db, nil
}

func connect(cfg *namespaceConfig) (clickhouse.Conn, error) {
	if cfg.Encoding != EncodingJSON && cfg.Encoding != EncodingProto {
		return nil, fmt.Errorf("unknown encoding %q, supported: %q, %q", cfg.Encoding, EncodingJSON, EncodingProto)
	}

	return cfg.Connector(cfg)
}

func (r *ClickHouseReader) GetConn() clickhouse.Conn {
	return r.db
}

func (r *ClickHouseReader) ClearTracesTable(ctx context.Context) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s.%s", r.traceDB, signozTraceTableName)
	err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error in processing truncate sql query: %v", err)
	}

	return nil
}

func (r *ClickHouseReader) ListTraces(ctx context.Context) (*[]model.SearchSpansResult, error) {
	var searchScanReponses []model.SearchSpanDBReponseItem
	query := fmt.Sprintf("SELECT * FROM %s.%s", r.traceDB, signozSpansTable)
	fmt.Println("query: ", query)
	err := r.db.Select(ctx, &searchScanReponses, query)
	if err != nil {
		return nil, fmt.Errorf("error in processing sql query: %v", err)
	}

	searchSpansResult := []model.SearchSpansResult{{
		Columns: []string{"__time", "SpanId", "TraceId", "ServiceName", "Name", "Kind", "DurationNano", "TagsKeys", "TagsValues", "References", "Events", "HasError"},
		Events:  make([][]interface{}, len(searchScanReponses)),
	},
	}

	for i, item := range searchScanReponses {
		var jsonItem model.SearchSpanReponseItem
		json.Unmarshal([]byte(item.Model), &jsonItem)
		jsonItem.TimeUnixNano = uint64(item.Timestamp.UnixNano() / 1000000)
		spanEvents := jsonItem.GetValues()
		searchSpansResult[0].Events[i] = spanEvents
	}

	return &searchSpansResult, nil

}

// addCondition ...
func addCondition(varName, varValue string, includeAND bool) string {
	// includeAND should only be used if there are multiple selectors in the WHERE clause.
	if varValue != "" {
		andPrefix := ""
		if includeAND {
			andPrefix += " AND "
		}
		return fmt.Sprintf(`%s(%s='%s')`, andPrefix, varName, varValue)
	}
	return ""
}

func (r *ClickHouseReader) GetTraces(ctx context.Context, selectors model.Selector) (*[]model.GetTracesDBResponse, error) {
	// Sample query
	// select serviceName, name, httpMethod, httpCode, httpRoute, httpHost, responseStatusCode
	// from signoz_traces.signoz_index_v2
	// where serviceName=goApp

	whereClause := "where " +
		addCondition("serviceName", selectors.ServiceName, false) +
		addCondition("name", selectors.Name, true) +
		addCondition("httpMethod", selectors.HttpMethod, true) +
		addCondition("httpCode", selectors.HttpCode, true) +
		addCondition("httpRoute", selectors.HttpRoute, true) +
		addCondition("httpHost", selectors.HttpHost, true)

	var getResponses []model.GetTracesDBResponse
	query := fmt.Sprintf("select traceID, spanID, serviceName, name, httpMethod, httpCode, httpRoute, httpHost, responseStatusCode from signoz_traces.signoz_index_v2 %s ", whereClause)
	// fmt.Println("query: ", query)
	err := r.db.Select(ctx, &getResponses, query)
	if err != nil {
		return nil, fmt.Errorf("error in processing sql query: %v", err)
	}

	return &getResponses, nil
}
