package aws

import (
	"context"
	"log"
	"os"

	"github.com/Sut103/discord-getting-messages-for-dynamodb/discord"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func getEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	endpoint := aws.Endpoint{}

	if url, exists := os.LookupEnv("DYNAMO_ENDPOINT"); exists {
		endpoint.URL = url
	}
	return endpoint, nil
}

type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDB() (DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(getEndpoint)))
	if err != nil {
		return DynamoDB{}, err
	}

	return DynamoDB{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: os.Getenv("DYNAMO_TABLE_NAME")}, nil
}

type image struct {
	Id   string                 `dynamodbav:"id"`
	Url  string                 `dynamodbav:"url"`
	Info map[string]interface{} `dynamodbav:"info"`
}

func (dd *DynamoDB) GetLatestId() (string, error) {
	response, err := dd.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(dd.tableName),
		Limit:     aws.Int32(1),
	})
	if err != nil {
		return "", err
	}

	var images []image
	err = attributevalue.UnmarshalListOfMaps(response.Items, &images)
	if err != nil {
		return "", err
	}

	if len(images) <= 0 {
		return "", nil
	}
	return images[0].Id, nil
}

type WChannelMessage struct {
	URL            string                 `json:"url"`
	ChannelMessage discord.ChannelMessage `json:"channel_message"`
}

func (dd *DynamoDB) InsertImageMessages(channelMessages []discord.ChannelMessage) error {
	var err error
	var item map[string]types.AttributeValue
	written := 0
	batchSize := 25
	start := 0
	end := start + batchSize
	for start < len(channelMessages) {
		var writeReqs []types.WriteRequest

		if end > len(channelMessages) {
			end = len(channelMessages)
		}

		for _, channelMessage := range channelMessages[start:end] {
			for _, attachment := range channelMessage.Attachments {
				var wChannelMessage WChannelMessage
				wChannelMessage.ChannelMessage = channelMessage
				wChannelMessage.URL = attachment.ProxyURL

				item, err = attributevalue.MarshalMap(wChannelMessage)
				if err != nil {
					log.Println(wChannelMessage.ChannelMessage.ID, err)

				} else {
					writeReqs = append(
						writeReqs,
						types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
					)
				}
			}
		}

		_, err = dd.client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{dd.tableName: writeReqs}})
		if err != nil {
			log.Println(dd.tableName, err)
		} else {

			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	return err
}
