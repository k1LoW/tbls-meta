package bq

import (
	"context"
	"fmt"
	"sort"
	"strings"

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

func (t *Table) FindColumnDescByName(name string) (string, error) {
	for _, i := range t.Columns {
		if i.Key.(string) == name {
			return i.Value.(string), nil
		}
	}
	return "", fmt.Errorf("not found column '%s'", name)
}

type Schema struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"description,omitempty"`
	Labels []string `yaml:"labels,omitempty"`
	Tables []*Table `yaml:"tables,omitempty"`
}

func (s Schema) FindTableByName(name string) (*Table, error) {
	for _, t := range s.Tables {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("not found table '%s'", name)
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
	f := makeSchema(from)
	t := makeSchema(to)

	ds := b.client.Dataset(b.datasetID)
	md, err := ds.Metadata(ctx)
	if err != nil {
		return err
	}
	mdu := bigquery.DatasetMetadataToUpdate{}

	// dataset.Description
	if f.Desc != t.Desc {
		fmt.Printf("Update: dataset.Description '%s' -> '%s'\n", f.Desc, t.Desc)
		mdu.Description = t.Desc
	}

	// dataset.Labels
	create, remove := diffLabels(f.Labels, t.Labels)
	if len(remove) > 0 {
		for _, l := range remove {
			fmt.Printf("Remove: dataset.Labels '%s' -> ''\n", l)
			splited := strings.Split(l, ":")
			mdu.DeleteLabel(splited[0])
		}
	}
	if len(create) > 0 {
		for _, l := range create {
			fmt.Printf("Create: dataset.Labels '' -> '%s'\n", l)
			splited := strings.Split(l, ":")
			mdu.SetLabel(splited[0], splited[1])
		}
	}

	if _, err := ds.Update(ctx, mdu, md.ETag); err != nil {
		return err
	}

	// tables
	for _, tt := range t.Tables {
		ft, err := f.FindTableByName(tt.Name)
		if err != nil {
			return err
		}

		bt := ds.Table(tt.Name)
		md, err := bt.Metadata(ctx)
		if err != nil {
			return err
		}
		mdu := bigquery.TableMetadataToUpdate{}

		if tt.Desc != ft.Desc {
			fmt.Printf("Update: table.%s.Description '%s' -> '%s'\n", tt.Name, ft.Desc, tt.Desc)
			mdu.Description = tt.Desc
		}
		// table.Labels
		create, remove := diffLabels(ft.Labels, tt.Labels)
		if len(remove) > 0 {
			for _, l := range remove {
				fmt.Printf("Remove: table.%s.Labels '%s' -> ''\n", tt.Name, l)
				splited := strings.Split(l, ":")
				mdu.DeleteLabel(splited[0])
			}
		}
		if len(create) > 0 {
			for _, l := range create {
				fmt.Printf("Create: table.%s.Labels '' -> '%s'\n", tt.Name, l)
				splited := strings.Split(l, ":")
				mdu.SetLabel(splited[0], splited[1])
			}
		}

		ts := md.Schema

		// columns
		for _, i := range tt.Columns {
			name := i.Key.(string)
			td := i.Value.(string)
			fd, err := ft.FindColumnDescByName(name)
			if err != nil {
				return err
			}
			if td != fd {
				fmt.Printf("Update: table.%s.%s.Description '%s' -> '%s'\n", tt.Name, name, fd, td)
				if err := setColumnDescription(ts, name, fd, td); err != nil {
					return err
				}
			}
		}

		mdu.Schema = ts

		if _, err := bt.Update(ctx, mdu, md.ETag); err != nil {
			return err
		}
	}

	return nil
}

func setColumnDescription(ts bigquery.Schema, name, fd, td string) error {
	for _, fs := range ts {
		if fs.Name == name {
			fs.Description = td
			return nil
		}
		if strings.Contains(name, ".") && len(fs.Schema) > 0 {
			splited := strings.SplitN(name, ".", 2)
			if err := setColumnDescription(fs.Schema, splited[1], fd, td); err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("not found column '%s'", name)
}

func diffLabels(from, to []string) ([]string, []string) {
	create := []string{}
	remove := []string{}

	for _, fl := range from {
		exist := false
		for _, tl := range to {
			if fl == tl {
				exist = true
				break
			}
		}
		if !exist {
			remove = append(remove, fl)
		}
	}

	for _, tl := range to {
		exist := false
		for _, fl := range from {
			if fl == tl {
				exist = true
				break
			}
		}
		if !exist {
			create = append(create, tl)
		}
	}

	return create, remove
}

func makeSchema(in *schema.Schema) *Schema {
	s := Schema{
		Name:   in.Name,
		Desc:   in.Desc,
		Labels: []string{},
		Tables: []*Table{},
	}

	for _, l := range in.Labels {
		s.Labels = append(s.Labels, l.Name)
	}
	sort.Strings(s.Labels)

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
		sort.Strings(table.Labels)

		for _, c := range t.Columns {
			item := yaml.MapItem{Key: c.Name, Value: c.Comment}
			table.Columns = append(table.Columns, item)
		}

		s.Tables = append(s.Tables, &table)
	}

	return &s
}
