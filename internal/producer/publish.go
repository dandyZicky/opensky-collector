package producer

type DefaultMessage struct {
	MessageKey   string
	MessageValue []byte
}

func (m DefaultMessage) Key() []byte {
	return []byte(m.MessageKey)
}

func (m DefaultMessage) Value() []byte {
	return m.MessageValue
}
