package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName string = "lists"
)

// DaoBuilds serializes to/from AWS DynamoDB
type DaoLists struct {
	svc *dynamodb.DynamoDB
}

type ListKey struct {
	CustomerID string
	ListName   string
}

type List struct {
	ListKey
	LastUpdated time.Time
	Items       []string
}

type ListRecordKey struct {
	UserID_ListName string
}
type ListRecord struct {
	ListRecordKey
	LastUpdated string
	ListItems   string
}

var awsRegion = flag.String("region", "us-east-1", "AWS region")
var svc = dynamodb.New(session.New(&aws.Config{Region: awsRegion}))

// NewDaoBuilds are used to perform CRUD operations on Builds
func NewDaoLists() (*DaoLists, error) {
	if svc != nil {
		return &DaoLists{svc: svc}, nil
	}
	return nil, fmt.Errorf("unable to construct DaoLists object with a <nil> DynamoDB attribute")
}

func formatListRecordKey(key ListKey) string {
	return key.CustomerID + ":" + key.ListName
}

func parseListRecordKey(key string) (customerID string, listName string, err error) {
	i := strings.Index(key, ":")
	if i == -1 {
		return "", "",
			fmt.Errorf("parseListRecordKey was unable to find the separator character ':' in the key %s!  We should all panic!", key)
	}
	if i == 0 || i >= len(key)-1 {
		return "", "", fmt.Errorf("parseListRecordKey found the separator character ':' at position %d in key %s", i, key)
	}

	return key[:i], key[i+1:], nil
}

func (dao *DaoLists) Persist(obj *List) error {

	o := &ListRecord{
		ListRecordKey: ListRecordKey{UserID_ListName: formatListRecordKey(obj.ListKey)},
		LastUpdated:   time.Now().Format(time.RFC3339Nano),
	}

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("item: %+v\n", item)

	params := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = dao.svc.PutItem(params)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (dao *DaoLists) Fetch(customerID string, listName string) (*List, error) {
	key := ListKey{CustomerID: customerID, ListName: listName}
	o := &ListRecordKey{
		UserID_ListName: formatListRecordKey(key),
	}

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       item,
	}

	log.Printf("params: %+v\n", params)

	resp, err := dao.svc.GetItem(params)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(resp.Item) == 0 {
		//		return nil, NewErrorCode(404, fmt.Errorf("object with key %s not found", o.ListRecordKey))
		return nil, fmt.Errorf("object with key %s not found", o.UserID_ListName)
	}

	out := &ListRecord{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, out)
	if err != nil {
		log.Fatal(err)
	}
	updated, err := time.Parse(time.RFC3339Nano, out.LastUpdated)
	if err != nil {
		updated = time.Time{}
	}
	log.Printf("ListRecord: %+v\n", out)

	customerID, listName, err = parseListRecordKey(out.UserID_ListName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("customerID: %s; listName: %s\n", customerID, listName)
	result := &List{ListKey: ListKey{CustomerID: customerID, ListName: listName}, LastUpdated: updated}
	return result, nil
}
