package main

import (
	"context"

	"github.com/Sut103/discord-getting-messages-for-dynamodb/aws"
	"github.com/Sut103/discord-getting-messages-for-dynamodb/discord"
	"github.com/aws/aws-lambda-go/lambda"
)

func Setup() (discord.ChannelAPI, aws.DynamoDB, error) {
	ca := discord.NewChannelAPI()
	dd, err := aws.NewDynamoDB()
	if err != nil {
		return ca, dd, err
	}
	return ca, dd, nil
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	ca, dd, err := Setup()
	if err != nil {
		panic(err)
	}

	latest_id, err := dd.GetLatestId()
	if err != nil {
		panic(err)
	}

	var cms []discord.ChannelMessage
	if latest_id == "" {
		cms, err = ca.GetChannelMessageAll()
		if err != nil {
			panic(err)
		}
	} else {
		cms, err = ca.GetChannelMessagesNewer(latest_id)
		if err != nil {
			panic(err)
		}
	}

	err = dd.InsertImageMessages(cms)
	if err != nil {
		panic(err)
	}

	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
