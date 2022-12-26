package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type GetTracesDBResponse struct {
	TraceID            string `ch:"traceID"`
	SpanID             string `ch:"spanID"`
	ServiceName        string `ch:"serviceName"`
	Name               string `ch:"name"`
	HttpMethod         string `ch:"httpMethod"`
	HttpRoute          string `ch:"httpRoute"`
	HttpHost           string `ch:"httpHost"`
	HttpCode           string `ch:"httpCode"`
	ResponseStatusCode string `ch:"responseStatusCode"`
}

type SearchSpanDBReponseItem struct {
	Timestamp time.Time `ch:"timestamp"`
	TraceID   string    `ch:"traceID"`
	Model     string    `ch:"model"`
}

type SearchSpansResult struct {
	Columns []string        `json:"columns"`
	Events  [][]interface{} `json:"events"`
}

type SearchSpanReponseItem struct {
	TimeUnixNano uint64            `json:"timestamp"`
	SpanID       string            `json:"spanID"`
	TraceID      string            `json:"traceID"`
	ServiceName  string            `json:"serviceName"`
	Name         string            `json:"name"`
	Kind         int32             `json:"kind"`
	References   []OtelSpanRef     `json:"references,omitempty"`
	DurationNano int64             `json:"durationNano"`
	TagMap       map[string]string `json:"tagMap"`
	Events       []string          `json:"event"`
	HasError     bool              `json:"hasError"`
}

type OtelSpanRef struct {
	TraceId string `json:"traceId,omitempty"`
	SpanId  string `json:"spanId,omitempty"`
	RefType string `json:"refType,omitempty"`
}

func (ref *OtelSpanRef) toString() string {

	retString := fmt.Sprintf(`{TraceId=%s, SpanId=%s, RefType=%s}`, ref.TraceId, ref.SpanId, ref.RefType)

	return retString
}

func (item *SearchSpanReponseItem) GetValues() []interface{} {

	references := []OtelSpanRef{}
	jsonbody, _ := json.Marshal(item.References)
	json.Unmarshal(jsonbody, &references)

	referencesStringArray := []string{}
	for _, item := range references {
		referencesStringArray = append(referencesStringArray, item.toString())
	}

	if item.Events == nil {
		item.Events = []string{}
	}
	keys := make([]string, 0, len(item.TagMap))
	values := make([]string, 0, len(item.TagMap))

	for k, v := range item.TagMap {
		keys = append(keys, k)
		values = append(values, v)
	}
	returnArray := []interface{}{item.TimeUnixNano, item.SpanID, item.TraceID, item.ServiceName, item.Name, strconv.Itoa(int(item.Kind)), strconv.FormatInt(item.DurationNano, 10), keys, values, referencesStringArray, item.Events, item.HasError}

	return returnArray
}
