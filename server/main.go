package main

import (
	"context"
	"fmt"
	"time"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"

	bookpb "github.com/jy0t1/grpc-mongo-crud/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookServiceServer struct {
}

// Select one book based on ID
func (s *BookServiceServer) ReadBook(ctx context.Context, req *bookpb.ReadBookReq) (*bookpb.ReadBookRes, error) {
	// convert string id (from proto) to mongoDB ObjectId
	// I have commented out this section as ID in Book collection (MongoDB) is numeric
	/*
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := bookdb.FindOne(ctx, bson.M{"_id": oid})
	*/
	fmt.Printf("In server/main.go begin readbook\n")
 
	fmt.Printf("Input parameter\n")
	fmt.Println(req)
	result := &bookpb.Book{}
	filter := req
	
	//var result primitive.D //  D is ordered and M is an unordered representation of a BSON document which is a Map
	fmt.Println("=====================in server/main.go=> ReadBook() ====================================")
	err1 := bookdb.FindOne(ctx, filter).Decode(&result)
	if err1 != nil {
		fmt.Println(err1)
	} 
	//fmt.Printf("In server/main.go readbook => Found a single document: %+v\n", result)
	
	// Cast to ReadBookRes type
	response := &bookpb.ReadBookRes{
		Book: &bookpb.Book{
			Id:   		 result.Id, //req.GetId(),
			Name: 		 result.Name,
			Author: 	 result.Author,
			AuthorEmail: result.AuthorEmail,
			Published: 	 result.Published,
			Pages: 		 result.Pages,
			Publisher: 	 result.Publisher,
			IsAvailable: result.IsAvailable,
			Category: 	 result.Category,
			BindType: 	 result.BindType,
			PhotoPath: 	 result.PhotoPath,
		},
	}
	fmt.Println(response)
	return response, nil
}


// Create a new book in library


func (s *BookServiceServer) CreateBook(ctx context.Context, req *bookpb.CreateBookReq) (*bookpb.CreateBookRes, error) {
	// Get the protobuf book type from the protobuf request type
	// Essentially doing req.Book to access the struct with a nil check
	fmt.Println("inside grpc-mongo-crud/server/main.go => createbook start")
	book := req.GetBook()
	fmt.Println(book)
	var _ID string
	fmt.Println("inside grpc-mongo-crud/server/main.go => createbook ............")
	result, err := bookdb.InsertOne(mongoCtx, book)
	//result, err := bookdb.InsertOne(context.TODO(), myBook)
	// check error
	if err != nil {
		// return internal gRPC error to be handled later
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	fmt.Println("inside grpc-mongo-crud/server/main.go => createbook after insert ............")
	oid := result.InsertedID.(primitive.ObjectID)
	_ID = oid.Hex()
	//fmt.Sprintf("Collection _ID : %s", _ID)
	fmt.Println(_ID)
	// return the book in a CreateBookRes type
	return &bookpb.CreateBookRes{Book: book}, nil
}

// Update a bnook in library

func (s *BookServiceServer) UpdateBook(ctx context.Context, req *bookpb.UpdateBookReq) (*bookpb.UpdateBookRes, error) {
	// Get the book data from the request
	book := req.GetBook()
    var err error
	// Convert the Id string to a MongoDB ObjectId
	
	/*
	oid, err := primitive.ObjectIDFromHex(book.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied book id to a MongoDB ObjectId: %v", err),
		)
	}
    */
	// document search filter
	filter := bson.M{"id": book.GetId()}
	// Convert the data to be updated into an unordered Bson document
	// insertone() always create elements in lowercase, so pass the element names in lower case here
	update := bson.M{
		"name": 		book.GetName(),
		"author": 	 	book.GetAuthor(),
		"authoremail": 	book.GetAuthorEmail(),
		"published": 	book.GetPublished(),
		"pages": 		book.GetPages(),
		"publisher": 	book.GetPublisher(),
		"isavailable":  book.GetIsAvailable(),
		"category": 	book.GetCategory(),
		"bindtype": 	book.GetBindType(),
		"photopath": 	book.GetPhotoPath(),
	}
	// Create an instance of an options and set the desired options
	// Upsert "true" means it will update if document found, otherwise insert new
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	//result := bookdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
    result := bookdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, &opt)
    
	// Decode result and write it to 'decoded'
	decoded := &bookpb.Book{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find book with supplied ID: %v", err),
		)
	}
	fmt.Println(decoded)
	return &bookpb.UpdateBookRes{
		Book: &bookpb.Book{
			Id:   		 decoded.Id,
			Name: 		 decoded.Name,
			Author: 	 decoded.Author,
			AuthorEmail: decoded.AuthorEmail,
			Published: 	 decoded.Published,
			Pages: 		 decoded.Pages,
			Publisher: 	 decoded.Publisher,
			IsAvailable: decoded.IsAvailable,
			Category: 	 decoded.Category,
			BindType: 	 decoded.BindType,
			PhotoPath: 	 decoded.PhotoPath,
		},
	}, nil
}

// delete book from library

func (s *BookServiceServer) DeleteBook(ctx context.Context, req *bookpb.DeleteBookReq) (*bookpb.DeleteBookRes, error) {
	// DeleteOne returns DeleteResult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	var err error
	_, err = bookdb.DeleteOne(ctx, bson.M{"id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete book with id %d: %v", req.GetId(), err))
	}
	return &bookpb.DeleteBookRes{
		Success: true,
	}, nil
}

// list of books from library

func (s *BookServiceServer) ListBooks(req *bookpb.ListBooksReq, stream bookpb.BookService_ListBooksServer) error {
	// Initiate a BookItem type to write decoded data to
	var err error
	data := &BookItem{}
	// collection.Find returns a cursor for our (empty) query
	cursor, err := bookdb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// An expression with defer will be called at the end of the function
	defer cursor.Close(context.Background())
	// cursor.Next() returns a boolean, if false there are no more items and loop will break
	for cursor.Next(context.Background()) {
		// Decode the data at the current pointer and write it to data
		err := cursor.Decode(data)
		// check error
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		// If no error is found send book over stream
		stream.Send(&bookpb.ListBooksRes{
			Book: &bookpb.Book{
				Id:   		 data.id,
				Name: 		 data.name,
				Author: 	 data.author,
				AuthorEmail: data.authorEmail,
				Published: 	 data.published,
				Pages: 		 data.pages,
				Publisher: 	 data.publisher,
				IsAvailable: data.isAvailable,
				Category: 	 data.category,
				BindType: 	 data.bindType,
				PhotoPath: 	 data.photoPath,
			},
		})
	}
	// Check if the cursor has any errors
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

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

//insertone() will always create elements with lowercase based on the following structure
type BookItem struct {
	id          int64 	`bson:"id,omitempty"`
	name 		string  `bson:"name"`
	author  	string  `bson:"author"`
	authorEmail string  `bson:"authorEmail"`
	published 	string 	`bson:"published"`
	pages 		int64  	`bson:"pages"`
	publisher 	string 	`bson:"publisher"`
    isAvailable bool 	`bson:"isAvailable"`
    category 	string 	`bson:"category"`
    bindType 	string 	`bson:"bindType"`
    photoPath 	string 	`bson:"photoPath"`
}

var db *mongo.Client
var bookdb *mongo.Collection
var mongoCtx context.Context

func main() {

	// Configure 'log' package to give file name and line number on eg. log.Fatal
	// just the filename & line number:
	// log.SetFlags(log.Lshortfile)
	// Or add timestamps and pipe file name and line number to it:
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting server on port :50051...")

	// 50051 is the default port for gRPC
	// Ideally we'd use 0.0.0.0 instead of localhost as well
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}
	// var s *grpc.Server 
	s := grpc.NewServer(opts...)
	// var srv *BookServiceServer
	srv := &BookServiceServer{}

	bookpb.RegisterBookServiceServer(s, srv)

	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://test_admin:test123@clusterjs-u8ha7.mongodb.net/library_db?retryWrites=true&w=majority"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = db.Connect(ctx)
	if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")


	bookdb = db.Database("library_db").Collection("book_records")
	fmt.Println("Collection type:", reflect.TypeOf(bookdb), "\n")

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server successfully started on port :50051")

	// Bad way to stop the server
	// if err := s.Serve(listener); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }
	// Right way to stop the server using a SHUTDOWN HOOK

	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed
	// And our main routine keeps running
	// If the main routine were to shutdown so would the child routine that is Serving the server
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}