package assertions

import (
	"github.com/ShashankGirish/clickhouse-integration/src/model"
)

func RunAssertions(traces *[]model.GetTracesDBResponse, asserts model.Assertion) (bool, *model.GetTracesDBResponse, error) {
	for _, trace := range *traces {
		if asserts.ResponseStatusCode == trace.ResponseStatusCode {
			// Found a matching trace, return true for success and the matching trace.
			return true, &trace, nil
		}
	}

	// Couldn't find any successful assertions.
	return false, nil, nil
}
