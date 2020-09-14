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


/*library_db => book_records
int64 id = 1;
    string name = 2;
    string author = 3;
    string authorEmail = 4;
    //google.protobuf.Timestamp published = 5;
    string published = 5;
    int64 pages = 6;
    string publisher = 7;
    bool isAvailable = 8;
    string category = 9;
    string bindType = 10;
    string photoPath = 11;
*/

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new book in library",
	Long: `Using gRPC this will create new book in library`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create called")
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
		fmt.Println(book)
		res, err := client.CreateBook(
			context.TODO(),
			&bookpb.CreateBookReq{
				Book: book,
			},
		)
		if err != nil {
			return err
		}
		fmt.Printf("New book created: %d [%T] %s\n", res.Book.Id, res.Book.Id, res.Book.Name)
		return nil
	},
}

func init() {
	fmt.Println("begining in init() in cmd/create.go ....................")
	createCmd.Flags().Int64P("id", "i", 0, "Book ID")
	createCmd.Flags().StringP("name", "n", "", "Name of the book")
	createCmd.Flags().StringP("author", "a", "", "Author")
	createCmd.Flags().StringP("authorEmail", "e", "", "Email-id of the author")
	createCmd.Flags().StringP("published", "p", "", "Published On [mm/dd/yyyy]")
	createCmd.Flags().Int64P("pages", "g", 0, "Pages")
	createCmd.Flags().StringP("publisher", "s", "", "Publisher")
	createCmd.Flags().BoolP("isAvailable", "v", false, "Available in library [true/false]")
	createCmd.Flags().StringP("category", "c", "", "Category")
	createCmd.Flags().StringP("bindType", "b", "", "Bind Type")
	createCmd.Flags().StringP("photoPath", "o", "", "Image path")
	createCmd.MarkFlagRequired("id")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("author")
	createCmd.MarkFlagRequired("pages")
	createCmd.MarkFlagRequired("publisher")
	createCmd.MarkFlagRequired("isAvailble")
	createCmd.MarkFlagRequired("category")
	createCmd.MarkFlagRequired("bindType")
	rootCmd.AddCommand(createCmd)
	fmt.Println("Ending in init() in cmd/create.go ....................")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
