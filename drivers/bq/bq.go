package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/k1LoW/tbls-meta/drivers"
	"github.com/k1LoW/tbls/schema"
	"gopkg.in/yaml.v2"
)

type Table struct {
	Name    string        `yaml:"name"`
	Desc    string        `yaml:"description,omitempty"`
	Labels  []string      `yaml:"labels,omitempty"`
	Columns yaml.MapSlice `yaml:"columns,omitempty"`
}

type Schema struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"description,omitempty"`
	Labels []string `yaml:"labels,omitempty"`
	Tables []Table  `yaml:"tables,omitempty"`
}

type Bigquery struct {
	client    *bigquery.Client
	datasetID string
}

// New returns new Bigquery
func New(client *bigquery.Client, datasetID string) *Bigquery {
	return &Bigquery{
		client:    client,
		datasetID: datasetID,
	}
}

func (b *Bigquery) Plan(ctx context.Context, from, to *schema.Schema) (drivers.Diffs, error) {
	diffs := drivers.Diffs{}

	f, err := yaml.Marshal(makeSchema(from))
	if err != nil {
		return diffs, err
	}

	t, err := yaml.Marshal(makeSchema(to))
	if err != nil {
		return diffs, err
	}

	// Description
	diffs = append(diffs, drivers.Diff{
		From:        "remove or replace",
		To:          "create",
		FromContent: string(f),
		ToContent:   string(t),
	})

	return diffs, nil
}

func (b *Bigquery) Apply(ctx context.Context, from, to *schema.Schema) error {
	return nil
}

func makeSchema(in *schema.Schema) *Schema {
	s := Schema{
		Name:   in.Name,
		Desc:   in.Desc,
		Labels: []string{},
		Tables: []Table{},
	}

	for _, l := range in.Labels {
		s.Labels = append(s.Labels, l.Name)
	}

	for _, t := range in.Tables {
		table := Table{
			Name:    t.Name,
			Desc:    t.Comment,
			Labels:  []string{},
			Columns: yaml.MapSlice{},
		}

		for _, l := range t.Labels {
			table.Labels = append(table.Labels, l.Name)
		}

		for _, c := range t.Columns {
			item := yaml.MapItem{Key: c.Name, Value: c.Comment}
			table.Columns = append(table.Columns, item)
		}

		s.Tables = append(s.Tables, table)
	}

	return &s
}
