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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update MongoDB document based on ID",
	Long: `As part of CRUD operation using gPRC APIs this is update operation`,
	RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("update called")
			id, err := cmd.Flags().GetInt64("id")
			name, err := cmd.Flags().GetString("name")
			author, err := cmd.Flags().GetString("author")
			authorEmail, err := cmd.Flags().GetString("authorEmail")
			published, err := cmd.Flags().GetString("published")
			pages, err := cmd.Flags().GetInt64("pages")
			publisher, err := cmd.Flags().GetString("publisher")
			isAvailable, err := cmd.Flags().GetBool("isAvailable")
			category, err := cmd.Flags().GetString("category")
			bindType, err := cmd.Flags().GetString("bindType")
			photoPath, err := cmd.Flags().GetString("photoPath")
			fmt.Println(id)
			fmt.Println(pages)
			fmt.Println(isAvailable)
			fmt.Println(name)
			if err != nil {
				return err
			}
			book := &bookpb.Book{
				Id: id,
				Name:    name ,
				Author:  author,
				AuthorEmail:  authorEmail,
				Published:  published,
				Pages: pages,
				Publisher: publisher,
				IsAvailable: isAvailable,
				Category: category,
				BindType: bindType,
				PhotoPath: photoPath,
			}
			res, err := client.UpdateBook(
				context.TODO(),
				&bookpb.UpdateBookReq{
					Book: book,
				},
			)
			if err != nil {
				return err
			}
			fmt.Printf("Updated book with ID : %d [%T]\n", res.Book.Id, res.Book.Id)
			return nil
		},
}

func init() {
	fmt.Println("begining in init() in cmd/update.go ....................")
	updateCmd.Flags().Int64P("id", "i", 0, "Book ID")
	updateCmd.Flags().StringP("name", "n", "", "Name of the book")
	updateCmd.Flags().StringP("author", "a", "", "Author")
	updateCmd.Flags().StringP("authorEmail", "e", "", "Email-id of the author")
	updateCmd.Flags().StringP("published", "p", "", "Published On [mm/dd/yyyy]")
	updateCmd.Flags().Int64P("pages", "g", 0, "Pages")
	updateCmd.Flags().StringP("publisher", "s", "", "Publisher")
	updateCmd.Flags().BoolP("isAvailable", "v", false, "Available in library [true/false]")
	updateCmd.Flags().StringP("category", "c", "", "Category")
	updateCmd.Flags().StringP("bindType", "b", "", "Bind Type")
	updateCmd.Flags().StringP("photoPath", "o", "", "Image path")
	updateCmd.MarkFlagRequired("id")
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("author")
	updateCmd.MarkFlagRequired("pages")
	updateCmd.MarkFlagRequired("publisher")
	updateCmd.MarkFlagRequired("isAvailble")
	updateCmd.MarkFlagRequired("category")
	updateCmd.MarkFlagRequired("bindType")
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
