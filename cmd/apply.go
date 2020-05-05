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
	"os"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply metadata (comments, description, labels) managed by tbls to the datasource",
	Long:  `apply metadata (comments, description, labels) managed by tbls to the datasource.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply(cmd, args)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runApply(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	from, to, driver, err := getSchemasAndDriver(ctx)
	if err != nil {
		return err
	}
	if err := driver.Apply(ctx, from, to); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
