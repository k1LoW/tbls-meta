/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/k1LoW/tbls-meta/drivers"
	"github.com/k1LoW/tbls-meta/drivers/bq"
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tbls-meta",
	Short: "tbls-meta is a CI-friendly tool for applying metadata managed by tbls to the datasource",
	Long:  `tbls-meta is a CI-friendly tool for applying metadata managed by tbls to the datasource.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getSchemasAndDriver(ctx context.Context) (*schema.Schema, *schema.Schema, drivers.Driver, error) {
	to, err := datasource.AnalyzeJSONString(os.Getenv("TBLS_SCHEMA"))
	if err != nil {
		return nil, nil, nil, err
	}
	dsn := os.Getenv("TBLS_DSN")
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	var driver drivers.Driver

	switch u.Scheme {
	case "bq", "bigquery":
		client, _, datasetID, err := datasource.NewBigqueryClient(ctx, dsn)
		if err != nil {
			return nil, nil, nil, err
		}
		driver = bq.New(client, datasetID)
	default:
		return nil, nil, nil, fmt.Errorf("unsupported driver '%s'", u.Scheme)
	}
	from, err := datasource.Analyze(config.DSN{
		URL: dsn,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return from, to, driver, nil
}

func init() {}
