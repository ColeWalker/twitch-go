package main

//IncomingMessage from websocket client
type IncomingMessage struct {
	Type string `json:"type"`
	Data struct {
		Message string `json:"message"`
		Topic   string `json:"topic"`
	} `json:"data"`
	Nonce string `json:"nonce"`
	Error string `json:"error"`
}

//RedemptionWrapper is the outer container for a redemption
type RedemptionWrapper struct {
	Type string `json:"type"`
	Data struct {
		Redemption RedemptionData `json:"redemption"`
	} `json:"data"`
}

//RedemptionData stores information about points redemptions
type RedemptionData struct {
	User struct {
		DisplayName string `json:"display_name"`
	}
	Reward struct {
		Title  string `json:"title"`
		Prompt string `json:"prompt"`
		Cost   int    `json:"cost"`
	}
	Input string `json:"user_input"`
}

//OutgoingMessage is the message that is sent to the server
type OutgoingMessage struct {
	Type  string `json:"type,omitempty"`
	Nonce string `json:"nonce,omitempty"`
	Data  struct {
		Topics    []string `json:"topics,omitempty"`
		AuthToken string   `json:"auth_token,omitempty"`
	} `json:"data,omitempty"`
}
