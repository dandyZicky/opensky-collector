package producer

type DefaultMessage struct {
	MessageKey   string
	MessageValue string
}

func (m DefaultMessage) Key() []byte {
	return []byte(m.MessageKey)
}

func (m DefaultMessage) Value() []byte {
	return []byte(m.MessageValue)
}
