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

//CommerceMsg stores information about responses from points redemptions
type CommerceMsg struct {
	UserName        string `json:"user_name"`
	DisplayName     string `json:"display_name"`
	ChannelName     string `json:"channel_name"`
	UserID          string `json:"user_id"`
	ChannelID       string `json:"channel_id"`
	Time            string `json:"time"`
	ItemImageURL    string `json:"item_image_url"`
	ItemDescription string `json:"item_description"`
	SupportsChannel bool   `json:"supports_channel"`
	PurchaseMessage struct {
		Message string `json:"message"`
		Emotes  []struct {
			Start int `json:"start"`
			End   int `json:"end"`
			ID    int `json:"id"`
		} `json:"emotes"`
	} `json:"purchase_message"`
}
