package mgorm

import (
	"github.com/globalsign/mgo"
)

type findOpts struct {
	Skip  int
	Limit int
	Sort  []string
}

func applyFindOpts(query *mgo.Query, opts ...FindOpt) *mgo.Query {
	fo := &findOpts{}
	for _, opt := range opts {
		opt(fo)
	}
	if len(fo.Sort) != 0 {
		query = query.Sort(fo.Sort...)
	}
	if fo.Skip != 0 {
		query = query.Skip(fo.Skip)
	}
	if fo.Limit != 0 {
		query = query.Limit(fo.Limit)
	}
	return query
}

// FindOpt is option for find
type FindOpt func(*findOpts)

// Skip on cursor
func Skip(skip int) FindOpt {
	return func(opts *findOpts) {
		opts.Skip = skip
	}
}

// Limit for result
func Limit(limit int) FindOpt {
	return func(opts *findOpts) {
		opts.Limit = limit
	}
}

// Sort result
func Sort(fields ...string) FindOpt {
	return func(opts *findOpts) {
		opts.Sort = append(opts.Sort, fields...)
	}
}

type updateOpts struct {
	Upsert bool
	Many   bool
}

// UpdateOpt is option for update
type UpdateOpt func(*updateOpts)

// Upsert option
func Upsert() UpdateOpt {
	return func(opts *updateOpts) {
		opts.Upsert = true
	}
}
