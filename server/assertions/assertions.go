package assertions

import (
	"github.com/kube-tarian/quality-trace/server/model"
)

func RunAssertions(traces *[]model.GetTracesDBResponse, asserts model.Assertion) *[]model.Result {
	var res []model.Result
	for _, trace := range *traces {
		if asserts.HttpMethod != "" {
			if asserts.HttpMethod == trace.HttpMethod {
				res = append(res, model.Result{
					Name: "Http Method",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "Http Method",
					Pass: "Fail",
				})
			}
		}
		if asserts.ResponseStatusCode != "" {
			if asserts.ResponseStatusCode == trace.ResponseStatusCode {
				res = append(res, model.Result{
					Name: "Response Status Code",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "Response Status Code",
					Pass: "Fail",
				})
			}
		}
		if asserts.ServiceName != "" {
			if asserts.ServiceName == trace.ServiceName {
				res = append(res, model.Result{
					Name: "Service Name",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "Service Name",
					Pass: "Fail",
				})
			}
		}
		if asserts.HttpHost != "" {
			if asserts.HttpHost == trace.HttpHost {
				res = append(res, model.Result{
					Name: "HTTP Host",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "HTTP Host",
					Pass: "Fail",
				})
			}
		}
		if asserts.HttpRoute != "" {
			if asserts.HttpRoute == trace.HttpRoute {
				res = append(res, model.Result{
					Name: "HTTP Route",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "HTTP Route",
					Pass: "Fail",
				})
			}
		}
		if asserts.Name != "" {
			if asserts.Name == trace.Name {
				res = append(res, model.Result{
					Name: "Name",
					Pass: "Pass",
				})
			} else {
				res = append(res, model.Result{
					Name: "Name",
					Pass: "Fail",
				})
			}
		}
	}
	return &res
}
