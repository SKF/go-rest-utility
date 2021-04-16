package problems

import (
	"context"
	"crypto/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/sanity-io/litter"
	"go.opencensus.io/trace"
)

func TestBasicProblem_DecorateWithRequest(t *testing.T) {
	rand.Reader = strings.NewReader("abcdefghijklmnopqrstuvwxyz")

	type fields struct {
		Type          string
		Title         string
		Status        int
		Detail        string
		Instance      string
		CorrelationID string
	}

	type args struct {
		ctx context.Context
		r   *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Problem
	}{
		{
			name: "context without span return copy of problem with instance set from request",
			fields: fields{
				Type:   "internal",
				Title:  "internal error",
				Status: 500,
				Detail: "an internal error occurred",
			},
			args: args{
				ctx: context.Background(),
				r: &http.Request{URL: &url.URL{
					Scheme: "https",
					Host:   "example.org",
					Path:   "/some/path",
				}},
			},
			want: BasicProblem{
				Type:     "internal",
				Title:    "internal error",
				Status:   500,
				Detail:   "an internal error occurred",
				Instance: "https://example.org/some/path",
			},
		},
		{
			name: "correlation id from context is used if present",
			fields: fields{
				Type:   "internal",
				Title:  "internal error",
				Status: 500,
				Detail: "an internal error occurred",
			},
			args: args{
				ctx: spanContext(),
				r: &http.Request{URL: &url.URL{
					Scheme: "https",
					Host:   "example.org",
					Path:   "/some/path",
				}},
			},
			want: BasicProblem{
				Type:          "internal",
				Title:         "internal error",
				Status:        500,
				Detail:        "an internal error occurred",
				Instance:      "https://example.org/some/path",
				CorrelationID: "15612916543113841209",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problem := BasicProblem{
				Type:          tt.fields.Type,
				Title:         tt.fields.Title,
				Status:        tt.fields.Status,
				Detail:        tt.fields.Detail,
				Instance:      tt.fields.Instance,
				CorrelationID: tt.fields.CorrelationID,
			}
			if got := problem.DecorateWithRequest(tt.args.ctx, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecorateWithRequest() = %v, want %v", litter.Sdump(got), litter.Sdump(tt.want))
			}
		})
	}
}

func spanContext() context.Context {
	ctx, _ := trace.StartSpan(context.Background(), "test")
	return ctx
}
