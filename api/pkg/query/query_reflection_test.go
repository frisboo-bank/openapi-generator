package query_test

import (
	"frisboo-bank/openapi-generator-service/pkg/query"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAllowedQueryFieldsFor(t *testing.T) {
	type Model struct {
		ID        uuid.UUID
		Title     string     `query:"title:search,filter,order"`
		Name      string     `query:"name:search"`
		HiddenAt  *time.Time `query:"hidden_at:filter"`
		CreatedAt *time.Time `query:"created_at:order"`
	}

	examples := []struct {
		name       string
		fn         func() query.AllowedQueryFields
		wantSearch []string
		wantFilter []string
		wantOrder  []string
		wantPanic  string
	}{
		{
			name:       "happy path",
			fn:         func() query.AllowedQueryFields { return query.AllowedQueryFieldsFor[Model]() },
			wantSearch: []string{"title", "name"},
			wantFilter: []string{"title", "hidden_at"},
			wantOrder:  []string{"title", "created_at"},
		},
		{
			name:       "pointer to struct",
			fn:         func() query.AllowedQueryFields { return query.AllowedQueryFieldsFor[*Model]() },
			wantSearch: []string{"title", "name"},
			wantFilter: []string{"title", "hidden_at"},
			wantOrder:  []string{"title", "created_at"},
		},
		{
			name: "no query tags",
			fn: func() query.AllowedQueryFields {
				type NoTagModel struct {
					ID    uuid.UUID
					Title string
				}

				return query.AllowedQueryFieldsFor[NoTagModel]()
			},
			wantSearch: nil,
			wantFilter: nil,
			wantOrder:  nil,
		},
		{
			name: "whitespace trimming",
			fn: func() query.AllowedQueryFields {
				type WhitespaceModel struct {
					Name string `query:"  name  :  filter  ,  order  "`
				}
				return query.AllowedQueryFieldsFor[WhitespaceModel]()
			},
			wantSearch: nil,
			wantFilter: []string{"name"},
			wantOrder:  []string{"name"},
		},
		{
			name: "unexported fields are still processed",
			fn: func() query.AllowedQueryFields {
				type UnexportedModel struct {
					Public  string `query:"public:search"`
					private string `query:"private:filter"`
				}

				return query.AllowedQueryFieldsFor[UnexportedModel]()
			},
			wantSearch: []string{"public"},
			wantFilter: []string{"private"},
			wantOrder:  nil,
		},

		// validation
		{
			name:      "non-struct model",
			fn:        func() query.AllowedQueryFields { return query.AllowedQueryFieldsFor[string]() },
			wantPanic: "model must be a struct, got string",
		},
		{
			name: "duplicate capability in same tag panics",
			fn: func() query.AllowedQueryFields {
				type DuplicateCapsModel struct {
					Name string `query:"name:filter,filter,order,order"`
				}
				return query.AllowedQueryFieldsFor[DuplicateCapsModel]()
			},
			wantPanic: "field \"name\" already present in filter capabilities for DuplicateCapsModel.Name",
		},
		{
			name: "duplicate field name from different struct fields panics",
			fn: func() query.AllowedQueryFields {
				type DuplicateFieldModel struct {
					A string `query:"x:filter"`
					B string `query:"x:filter"`
				}
				return query.AllowedQueryFieldsFor[DuplicateFieldModel]()
			},
			wantPanic: "field \"x\" already present in filter capabilities for DuplicateFieldModel.B",
		},
		{
			name: "invalid tag missing colon",
			fn: func() query.AllowedQueryFields {
				type BadTagModel struct {
					A string `query:"filter"`
				}
				return query.AllowedQueryFieldsFor[BadTagModel]()
			},
			wantPanic: "invalid query tag \"filter\" on field BadTagModel.A: expected format \"name:cap1,cap2\"",
		},
		{
			name: "invalid capability",
			fn: func() query.AllowedQueryFields {
				type BadCapModel struct {
					A string `query:"a:unknown"`
				}
				return query.AllowedQueryFieldsFor[BadCapModel]()
			},
			wantPanic: "invalid capability \"unknown\" in query tag on field BadCapModel.A (allowed: filter, order, search)",
		},
		{
			name: "invalid tag empty field name",
			fn: func() query.AllowedQueryFields {
				type BadTagModel struct {
					A string `query:":filter"`
				}
				return query.AllowedQueryFieldsFor[BadTagModel]()
			},
			wantPanic: "invalid query tag \":filter\" on field BadTagModel.A: expected format \"name:cap1,cap2\"",
		},
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			if e.wantPanic != "" {
				assert.PanicsWithError(t, e.wantPanic, func() { e.fn() })
				return
			}

			got := e.fn()

			assert.Equal(t, e.wantSearch, got.Search)
			assert.Equal(t, e.wantFilter, got.Filter)
			assert.Equal(t, e.wantOrder, got.Order)
		})
	}
}
