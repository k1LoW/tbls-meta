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
	"os"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "show the difference between metadata of the datasource and metadata managed by tbls",
	Long:  `show the difference between metadata of the datasource and metadata managed by tbls.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runPlan(cmd, args)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runPlan(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	from, to, driver, err := getSchemasAndDriver(ctx)
	if err != nil {
		return err
	}
	diffs, err := driver.Plan(ctx, from, to)
	if err != nil {
		return err
	}

	diff := "tbls meta plan:\n"
	for _, df := range diffs {
		d := difflib.UnifiedDiff{
			A:        difflib.SplitLines(df.FromContent),
			B:        difflib.SplitLines(df.ToContent),
			FromFile: df.From,
			ToFile:   df.To,
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(d)
		if text != "" {
			diff += text
		}
	}
	fmt.Println(diff)
	return nil
}

func init() {
	rootCmd.AddCommand(planCmd)
}
