package core

import (
	"context"
	"github.com/shurcooL/githubv4"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var SampleLanguageSizes = []LanguageSize{
	{
		Name:       "TypeScript",
		Size:       545992,
		Percentage: 72.87,
	}, {
		Name:       "Vue",
		Size:       103000,
		Percentage: 13.75,
	}, {
		Name:       "HCL",
		Size:       100249,
		Percentage: 13.38,
	},
}

func TestOctclient_ConvertJson(t *testing.T) {
	type fields struct {
		user            string
		client          *githubv4.Client
		oldestUpdatedAt time.Time
		latestUpdatedAt time.Time
	}
	type args struct {
		r *Results
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "JSONで返ってくる",
			args: args{
				r: &Results{
					UpdatedRange: UpdatedRange{
						Oldest: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2020-08-02T16:43:48Z")
							return t
						}(),
						Latest: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2021-07-18T06:36:11Z")
							return t
						}(),
					},
					LanguageSizes: SampleLanguageSizes,
				},
			},
			want: `{
  "updated_range": {
    "oldest": "2020-08-02T16:43:48Z",
    "latest": "2021-07-18T06:36:11Z"
  },
  "language_sizes": [
    {
      "name": "TypeScript",
      "size": 545992,
      "percentage": 72.87
    },
    {
      "name": "Vue",
      "size": 103000,
      "percentage": 13.75
    },
    {
      "name": "HCL",
      "size": 100249,
      "percentage": 13.38
    }
  ]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &Octclient{
				user:            tt.fields.user,
				client:          tt.fields.client,
				oldestUpdatedAt: tt.fields.oldestUpdatedAt,
				latestUpdatedAt: tt.fields.latestUpdatedAt,
			}
			got, err := oc.ConvertJson(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOctclient_ConvertTableForMarkdown(t *testing.T) {
	type args struct {
		r *Results
		o *MarkdownOptions
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Markdownの形式が崩れない",
			args: args{
				r: &Results{
					UpdatedRange:  UpdatedRange{},
					LanguageSizes: SampleLanguageSizes,
				},
				o: &MarkdownOptions{
					IsEachExtension: false,
				},
			},
			want: `|language|percentage(%)|size(byte)|
|---|---|---|
|TypeScript|72.87|545992|
|Vue|13.75|103000|
|HCL|13.38|100249|
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &Octclient{}
			if got := oc.ConvertTableForMarkdown(tt.args.r, tt.args.o); got != tt.want {
				t.Errorf("ConvertTableForMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func mustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

func mustRead(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestOctclient_GetRepositoriesContributedTo(t *testing.T) {
	type fields struct {
		user            string
		client          *githubv4.Client
		oldestUpdatedAt time.Time
		latestUpdatedAt time.Time
	}
	type args struct {
		ctx          context.Context
		isSortBySize bool
		reverse      bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Results
		wantErr bool
	}{
		{
			name: "言語ごとの合計サイズを計算する",
			fields: fields{
				user: "hoge",
				client: func() *githubv4.Client {
					mux := http.NewServeMux()
					mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
						if got, want := req.Method, http.MethodPost; got != want {
							t.Errorf("got request method: %v, want: %v", got, want)
						}
						_ = mustRead(req.Body)
						//if got, want := body, `{"query":"{viewer{login,bio}}"}`+"\n"; got != want {
						//	t.Errorf("got body: %v, want %v", got, want)
						//}
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, `
{
  "data": {
    "user": {
      "repositoriesContributedTo": {
        "nodes": [
          {
            "updatedAt": {
              "time": "2020-08-03T17:43:50Z"
            },
            "languages": {
              "edges": [
                {
                  "size": 100000,
                  "node": {
                    "name": "TypeScript"
                  }
                }
              ]
            }
          },
          {
            "updatedAt": {
              "time": "2020-08-01T16:00:48Z"
            },
            "languages": {
              "edges": [
                {
                  "size": 445992,
                  "node": {
                    "name": "TypeScript"
                  }
                },
                {
                  "size": 103000,
                  "node": {
                    "name": "Vue"
                  }
                }
              ]
            }
          },
          {
            "updatedAt": {
              "time": "2020-07-29T10:10:38Z"
            },
            "languages": {
              "edges": [
                {
                  "size": 100249,
                  "node": {
                    "name": "HCL"
                  }
                }
              ]
            }
          }
        ]
      }
    }
  }
}
`)
					})
					return githubv4.NewClient(&http.Client{Transport: localRoundTripper{handler: mux}})
				}(),
				oldestUpdatedAt: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2020-08-02T16:43:48Z")
					return t
				}(),
				latestUpdatedAt: time.Time{},
			},
			args: args{
				ctx:          context.Background(),
				isSortBySize: true,
				reverse:      false,
			},
			want: &Results{
				UpdatedRange: UpdatedRange{
					Oldest: func() time.Time {
						t, _ := time.Parse(time.RFC3339, "2020-07-29T10:10:38Z")
						return t
					}(),
					Latest: func() time.Time {
						t, _ := time.Parse(time.RFC3339, "2020-08-03T17:43:50Z")
						return t
					}(),
				},
				LanguageSizes: SampleLanguageSizes,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Octclient{
				user:            tt.fields.user,
				client:          tt.fields.client,
				oldestUpdatedAt: tt.fields.oldestUpdatedAt,
				latestUpdatedAt: tt.fields.latestUpdatedAt,
			}
			got, err := o.GetRepositoriesContributedTo(tt.args.ctx, tt.args.isSortBySize, tt.args.reverse)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepositoriesContributedTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepositoriesContributedTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
