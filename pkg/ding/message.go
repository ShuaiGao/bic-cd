package ding

func NewMsg() *Message {
	return &Message{}
}

func NewMsgWithWebhook(webhook, secret string) *Message {
	return &Message{
		webhook: webhook,
		secret:  secret,
	}
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type FeedCard struct {
	Links struct {
		PicURL     string `json:"picURL"`
		MessageURL string `json:"messageURL"`
		Title      string `json:"title"`
	} `json:"links"`
}

type ActionCard struct {
	HideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	Btns           []struct {
		ActionURL string `json:"actionURL"`
		Title     string `json:"title"`
	} `json:"btns"`
	Text      string `json:"text"`
	SingleURL string `json:"singleURL"`
	Title     string `json:"title"`
}

type Link struct {
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl"`
	Text       string `json:"text"`
	Title      string `json:"title"`
}

type Message struct {
	MsgType    string                 `json:"msgtype"`
	Text       map[string]interface{} `json:"text,omitempty"`
	Markdown   Markdown               `json:"markdown,omitempty"`
	ActionCard ActionCard             `json:"actionCard"`
	Link       Link                   `json:"link"`
	FeedCard   FeedCard               `json:"feedCard"`
	At         At                     `json:"at"`
	webhook    string
	secret     string
}

func (m *Message) SetWebhook(webhook string) *Message {
	m.webhook = webhook
	return m
}

func (m *Message) SetSecret(secret string) *Message {
	m.secret = secret
	return m
}

func (m *Message) SetText(text string) *Message {
	m.MsgType = "text"
	m.Text = map[string]interface{}{
		"content": text,
	}
	return m
}

func (m *Message) SetMarkdown(markdown Markdown) *Message {
	m.MsgType = "markdown"
	m.Markdown = markdown
	return m
}

func (m *Message) SetAtAll() *Message {
	m.At.IsAtAll = true
	return m
}

// SetAt 不支持纯@消息，@信息需附带到其他消息发送
func (m *Message) SetAt(at At) *Message {
	m.At = at
	return m
}

func (m *Message) SetLink(link Link) *Message {
	m.MsgType = "link"
	m.Link = link
	return m
}

func (m *Message) SetActionCard(actionCard ActionCard) *Message {
	m.MsgType = "actionCard"
	m.ActionCard = actionCard
	return m
}

func (m *Message) SetFeedCard(feedCard FeedCard) *Message {
	m.MsgType = "feedCard"
	m.FeedCard = feedCard
	return m
}
