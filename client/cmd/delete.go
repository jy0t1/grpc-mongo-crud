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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete document based on ID",
	Long: `gRPC based CRUD operation - deleting document based on ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("read called in cmd/delete.go")
		id, err := cmd.Flags().GetInt64("id")
		fmt.Println(id)
		if err != nil {
			return err
		}
		req := &bookpb.DeleteBookReq{
			Id: id,
		}
		fmt.Println("before readbook called in cmd/read.go ..............")
		res, err := client.DeleteBook(context.Background(), req)
		fmt.Println("After readbook called in cmd/read.go ..............")
		if err != nil {
			return err
		} 
		fmt.Println(res)
		return nil
	},
}

func init() {
	fmt.Println("inside init() in cmd/delete.go ...............")
	deleteCmd.Flags().Int64P("id", "i", 1, "The id of the book")
	deleteCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
