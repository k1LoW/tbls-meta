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
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/k1LoW/tbls-meta/drivers"
	"github.com/k1LoW/tbls-meta/drivers/bq"
	"github.com/k1LoW/tbls-meta/version"
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tbls-meta",
	Short: "tbls-meta is an external subcommand of tbls for applying metadata managed by tbls to the datasource",
	Long:  `tbls-meta is an external subcommand of tbls for applying metadata managed by tbls to the datasource.`,
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	log.SetOutput(ioutil.Discard)
	if env := os.Getenv("DEBUG"); env != "" {
		debug, err := os.Create(fmt.Sprintf("%s.debug", version.Name))
		if err != nil {
			printFatalln(rootCmd, err)
		}
		log.SetOutput(debug)
	}

	if err := rootCmd.Execute(); err != nil {
		printFatalln(rootCmd, err)
	}
}

// https://github.com/spf13/cobra/pull/894
func printErrln(c *cobra.Command, i ...interface{}) {
	c.PrintErr(fmt.Sprintln(i...))
}

func printErrf(c *cobra.Command, format string, i ...interface{}) {
	c.PrintErr(fmt.Sprintf(format, i...))
}

func printFatalln(c *cobra.Command, i ...interface{}) {
	printErrln(c, i...)
	os.Exit(1)
}

func printFatalf(c *cobra.Command, format string, i ...interface{}) {
	printErrf(c, format, i...)
	os.Exit(1)
}

func getSchemasAndDriver(ctx context.Context) (*schema.Schema, *schema.Schema, drivers.Driver, error) {
	to, err := datasource.AnalyzeJSONStringOrFile(os.Getenv("TBLS_SCHEMA"))
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
