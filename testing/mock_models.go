package tests

import "mychat/db"

type MessageClientMock struct {
	Messages   map[string]db.Message
	CallLog    []string
	DeleteLogs []string
}

func NewMessageClientMock() *MessageClientMock {
	return &MessageClientMock{
		Messages: make(map[string]db.Message),
	}
}

func (m *MessageClientMock) Create(id string) db.Message {
	msg := db.Message{ID: len(m.Messages) + 1, Text: id}
	m.Messages[id] = msg
	m.CallLog = append(m.CallLog, "Create:"+id)
	return msg
}

func (m *MessageClientMock) Read(id string) (db.Message, bool) {
	m.CallLog = append(m.CallLog, "Read:"+id)
	if msg, ok := m.Messages[id]; ok {
		return msg, true
	}
	return db.Message{}, false
}

func (m *MessageClientMock) ReadAll() []db.Message {
	m.CallLog = append(m.CallLog, "ReadAll")
	all := []db.Message{}
	for _, msg := range m.Messages {
		all = append(all, msg)
	}
	return all
}

func (m *MessageClientMock) Delete(id string) (db.Message, bool) {
	m.DeleteLogs = append(m.DeleteLogs, id)
	if msg, ok := m.Messages[id]; ok {
		delete(m.Messages, id)
		m.CallLog = append(m.CallLog, "Delete:"+id)
		return msg, true
	}
	return db.Message{}, false
}

func (m *MessageClientMock) DeleteAll() int {
	count := len(m.Messages)
	m.Messages = make(map[string]db.Message)
	m.CallLog = append(m.CallLog, "DeleteAll")
	return count
}
