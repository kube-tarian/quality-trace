package assertions

import (
	"fmt"

	"github.com/kube-tarian/quality-trace/server/model"
)

func RunAssertions(traces *[]model.GetTracesDBResponse, asserts model.Assertion) (bool, *model.GetTracesDBResponse, error) {
	for _, trace := range *traces {
		fmt.Println(asserts.ResponseStatusCode, trace.ResponseStatusCode, trace.HttpCode)
		if asserts.ResponseStatusCode == trace.ResponseStatusCode || asserts.ResponseStatusCode == trace.HttpCode {
			// Found a trace that matches the assertion.
			fmt.Println("Trace Status is :", trace.ResponseStatusCode)
			return true, &trace, nil
		}
	}

	// No traces found which match the assertion.
	return false, nil, nil
}
