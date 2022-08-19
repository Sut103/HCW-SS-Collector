package main

import (
	"log"

	"github.com/Sut103/discord-getting-messages-for-dynamodb/aws"
	"github.com/Sut103/discord-getting-messages-for-dynamodb/discord"
)

func Setup() (discord.ChannelAPI, aws.DynamoDB, error) {
	ca := discord.NewChannelAPI()
	dd, err := aws.NewDynamoDB()
	if err != nil {
		return ca, dd, err
	}
	return ca, dd, nil
}

func main() {
	ca, dd, err := Setup()
	if err != nil {
		log.Fatalln(err)
	}

	latest_id, err := dd.GetLatestId()
	if err != nil {
		panic("")
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
}
