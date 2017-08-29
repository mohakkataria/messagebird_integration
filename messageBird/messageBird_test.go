package messageBird

import (
    "testing"
    "github.com/mohakkataria/messagebird_integration/models"
    "github.com/spf13/viper"
    "fmt"
)

func TestQueueMessage(t *testing.T) {
    msg := &models.Message{
        Recipients:[]string{"123, 123"},
        Originator:"MessageBird",
        MessageBody:"test",
    }

    QueueMessage(msg)

    msg = &models.Message{
        Recipients:[]string{"123, 123"},
        Originator:"MessageBird",
        MessageBody:"test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test test ",
    }

    QueueMessage(msg)

    msg = &models.Message{
        Recipients:[]string{"123, 123"},
        Originator:"MessageBird",
        MessageBody:"日本語 ",
    }

    QueueMessage(msg)

    msg = &models.Message{
        Recipients:[]string{"123, 123"},
        Originator:"MessageBird",
        MessageBody:"тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест тестестестест ",
    }

    QueueMessage(msg)

    msg = &models.Message{
        Recipients:[]string{"123, 123"},
        Originator:"MessageBird",
        MessageBody:"日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 日本語 ",
    }

    QueueMessage(msg)
}

func init() {
    InitializeAPIHits()
    viper.SetConfigFile("./../config.json")
    err := viper.ReadInConfig()
    if err != nil {
        fmt.Println("No configuration file loaded")
    }
}