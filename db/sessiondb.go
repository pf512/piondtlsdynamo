package sessiondb

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	"time"
)

var (
	svc *dynamodb.DynamoDB
	tableName = "dtls"
)
// Create struct to hold info about new item
type DtlsSession struct {
	ID   		string    `json:"id"`
	Secret  	string    `json:"secret"`
	Address  	string    `json:"address"`
	Expiration 	time.Time `json:"expire_at"`
}


func init() {
	InitDynamo()
}

func InitDynamo() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		println("dynamodb init failed")
		os.Exit(1)
	}

	// Create DynamoDB client
	svc = dynamodb.New(sess)
}

// 	sessiondb.StoreSession(d.ID, d.Secret, d.Addr, d.ExpireAt)
func StoreSession(id string, secret string, address string, expiration time.Time) {

	item := DtlsSession{
		ID: id,
		Secret: secret,
		Address: address,
		Expiration: expiration,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	// Create item in table Movies
	// tableName := "dtls"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added ID '" + item.ID + "' to table " + tableName)
}

func RetrieveSession(id string) (DtlsSession, error){

	//tableName := "dtls"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return DtlsSession{}, err
	}

	if result.Item == nil {
		msg := "Could not find '" + id + "'"
		//return nil, errors.New(msg)
		println(msg)
		return DtlsSession{}, errors.New(msg)
	}

	dtlsSession := DtlsSession{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &dtlsSession)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println("Found dtls session:")
	fmt.Println("ID:  ", dtlsSession.ID)
	fmt.Println("Secret: ", dtlsSession.Secret)
	fmt.Println("Address:  ", dtlsSession.Address)
	fmt.Println("Expiration:", dtlsSession.Expiration)

	return dtlsSession, nil
}


func ListAllTables() {


	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

func DeleteSessionID(id string) error {

	//tableName := "Movies"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}

	fmt.Println("Deleted '" + id + "' from table " + tableName)

	return nil

}
