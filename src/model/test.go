package model

type Test struct {
	Name    string  `yaml:"name,omitempty"`
	Trigger Trigger `yaml:"trigger,omitempty"`
	Specs   []*Spec `yaml:"specs,omitempty"`
}

type Trigger struct {
	Type        string      `yaml:"type"`
	HttpRequest HttpRequest `yaml:"httpRequest"`
}

type HttpRequest struct {
	Url     string        `yaml:"url"`
	Route   string        `yaml:"route"`
	Method  string        `yaml:"method"`
	Headers []HttpHeaders `yaml:"headers"`
	Body    string        `yaml:"body"`
}

type HttpHeaders struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Spec struct {
	Name          string    `yaml:"name"`
	Selectors     Selector  `yaml:"selectors"`
	Assertions    Assertion `yaml:"assertions"`
	MaxRetries    int       `yaml:"maxRetries,omitempty"`
	RetryInterval int       `yaml:"retryInterval,omitempty"`
}

type Selector struct {
	ServiceName string `yaml:"serviceName"` //Ex: goApp
	Name        string `yaml:"name"`        //Ex: /books/:id
	HttpMethod  string `yaml:"httpMethod"`  //Ex: GET
	HttpRoute   string `yaml:"httpRoute"`   //Ex: /books/:id
	HttpHost    string `yaml:"httpHost"`    //Ex: localhost:8090
	HttpCode    string `yaml:"httpCode"`    //Ex: 200
}

// select serviceName, name, httpMethod, httpCode, httpRoute, httpHost, responseStatusCode from signoz_traces.signoz_index_v2;
type Assertion struct {
	Body               string `yaml:"bodyContains"`
	ResponseStatusCode string `yaml:"responseStatusCode"`
}
