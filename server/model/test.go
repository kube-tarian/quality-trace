package model

type Test struct {
	Id      int     `yaml:"index,omitempty" json:"index,omitempty"`
	Status  string  `yaml:"status,omitempty" json:"status,omitempty"`
	Name    string  `yaml:"name,omitempty" json:"name,omitempty"`
	Trigger Trigger `yaml:"trigger,omitempty" json:"trigger,omitempty"`
	Specs   []*Spec `yaml:"specs,omitempty" json:"specs,omitempty"`
}

type Trigger struct {
	Type        string      `yaml:"type" json:"type"`
	HttpRequest HttpRequest `yaml:"httpRequest" json:"httpRequest"`
}

type HttpRequest struct {
	Url     string        `yaml:"url" json:"url"`
	Route   string        `yaml:"route" json:"route"`
	Method  string        `yaml:"method" json:"method"`
	Headers []HttpHeaders `yaml:"headers" json:"headers"`
	Body    string        `yaml:"body" json:"body"`
}

type HttpHeaders struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

type Spec struct {
	Name          string    `yaml:"name" json:"name"`
	Selectors     Selector  `yaml:"selectors" json:"selectors"`
	Assertions    Assertion `yaml:"assertions" json:"assertions"`
	MaxRetries    int       `yaml:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	RetryInterval int       `yaml:"retryInterval,omitempty" json:"retryInterval,omitempty"`
}

type Selector struct {
	//Ex: goApp
	ServiceName string `yaml:"serviceName" json:"serviceName"`
	//Ex: /books/:id
	Name string `yaml:"name" json:"name"`
	//Ex: GET
	HttpMethod string `yaml:"httpMethod" json:"httpMethod"`
	//Ex: /books/:id
	HttpRoute string `yaml:"httpRoute" json:"httpRoute"`
	//Ex: localhost:8090
	HttpHost string `yaml:"httpHost" json:"httpHost"`
	//Ex: 200
	HttpCode string `yaml:"httpCode" json:"httpCode"`
}

// select serviceName, name, httpMethod, httpCode, httpRoute, httpHost, responseStatusCode from signoz_traces.signoz_index_v2;
type Assertion struct {
	Body               string `yaml:"bodyContains" json:"bodyContains"`
	ResponseStatusCode string `yaml:"responseStatusCode" json:"responseStatusCode"`
}
