package main

import (
	"fmt"
	"strings"
)

//Topic for PubSubClient
type Topic string

//Type of Topic
type Type string

const (
	// Bits event prefix
	Bits = Type("channel-bits-events-v2")
	// ChannelPoints event prefix
	ChannelPoints = Type("channel-points-channel-v1")
	// BitsBadgeNotification event prefix
	BitsBadgeNotification = Type("channel-bits-badge-unlocks")
	// Subscriptions event prefix
	Subscriptions = Type("channel-subscribe-events-v1")
	//Whispers event prefix
	Whispers = Type("whispers")
	//ModerationAction event prefix
	ModerationAction = Type("chat_moderator_actions")
	//Invalid event prefix
	Invalid = Type("invalid")
	//Pong event prefix
	Pong = Type("PONG")
	//Reconnect event prefix
	Reconnect = Type("RECONNECT")
)

//GetType of message
func GetType(topic string) Type {
	pieces := strings.Split(topic, ".")

	if len(pieces) < 2 {
		return Invalid
	}

	switch Type(pieces[0]) {
	case Bits:
		return Bits
	case BitsBadgeNotification:
		return BitsBadgeNotification
	case Subscriptions:
		return Subscriptions
	case ChannelPoints:
		return ChannelPoints
	case Whispers:
		return Whispers
	case ModerationAction:
		return ModerationAction
	}

	return Invalid
}

// GetBitsTopic -> Topic
func GetBitsTopic(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-bits-events-v2.%d", channelID))
}

// GetBitsBadgeNotificationTopic -> Topic
func GetBitsBadgeNotificationTopic(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-bits-badge-unlocks.%d", channelID))
}

// GetSubscriptionsTopic -> Topic
func GetSubscriptionsTopic(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-subscribe-events-v1.%d", channelID))
}

// GetCommerceTopic -> Topic
func GetCommerceTopic(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-commerce-events-v1.%d", channelID))
}

// GetWhispersTopic -> Topic
func GetWhispersTopic(channelID int) Topic {
	return Topic(fmt.Sprintf("whispers.%d", channelID))
}

//GetModerationActionTopic -> Topic
func GetModerationActionTopic(userID int, channelID int) Topic {
	return Topic(fmt.Sprintf("chat_moderator_actions.%d.%d", userID, channelID))
}
