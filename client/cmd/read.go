/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	
	bookpb "github.com/jy0t1/grpc-mongo-crud/proto"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Find a book based on ID",
	Long: `Find a book based on ID from MOngoDB collection`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("read called in cmd/read.go")
		id, err := cmd.Flags().GetInt64("id")
		fmt.Println(id)
		if err != nil {
			return err
		}
		req := &bookpb.ReadBookReq{
			Id: id,
		}
		fmt.Println("before readbook called in cmd/read.go ..............")
		res, err := client.ReadBook(context.Background(), req)
		fmt.Println("After readbook called in cmd/read.go ..............")
		if err != nil {
			return err
		} 
		fmt.Println(res)
		return nil
		},
}

func init() {
	fmt.Println("inside init() in cmd/read.go ...............")
	readCmd.Flags().Int64P("id", "i", 1, "The id of the book")
	readCmd.MarkFlagRequired("id")
	fmt.Println("inside init() in cmd/read.go after id capture ............... ")
	rootCmd.AddCommand(readCmd)
	fmt.Println("inside init() in cmd/read.go after rootCmd.AddCommand(readCmd) ............... ")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
