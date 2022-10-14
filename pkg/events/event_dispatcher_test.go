package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetName() string {
	return e.Name
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface) {
}

type EventDispatecherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher EventDispatcherInterface
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Register() {
	ed := NewEventDispatcher()
	eh := &TestEventHandler{}
	eh2 := &TestEventHandler{}
	ed.Register("test", eh)
	suite.Equal(1, len(ed.handlers["test"]))

	ed.Register("test", eh2)
	suite.Equal(2, len(ed.handlers["test"]))

	ed.Register("test2", &TestEventHandler{})
	suite.Equal(1, len(ed.handlers["test2"]))

	assert.Equal(suite.T(), eh, ed.handlers["test"][0])
	assert.Equal(suite.T(), eh2, ed.handlers["test"][1])
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Clear() {
	ed := NewEventDispatcher()
	ed.Register("test", &TestEventHandler{})
	ed.Register("test", &TestEventHandler{})
	ed.Register("test2", &TestEventHandler{})
	ed.Clear()
	suite.Equal(0, len(ed.handlers["test"]))
	suite.Equal(0, len(ed.handlers["test2"]))
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Remove() {
	ed := NewEventDispatcher()
	eh := &TestEventHandler{ID: 1}
	eh2 := &TestEventHandler{ID: 2}
	ed.Register("test", eh)
	ed.Register("test", eh2)
	suite.Equal(2, len(ed.handlers["test"]))
	ed.Remove("test", eh)
	suite.Equal(1, len(ed.handlers["test"]))
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Has() {
	ed := NewEventDispatcher()
	eh := &TestEventHandler{ID: 1}
	eh2 := &TestEventHandler{ID: 2}
	ed.Register("test", eh)
	ed.Register("test", eh2)
	suite.Equal(2, len(ed.handlers["test"]))
	assert.True(suite.T(), ed.Has("test", eh))
	assert.True(suite.T(), ed.Has("test", eh2))
	assert.False(suite.T(), ed.Has("test", &TestEventHandler{ID: 3}))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Dispatch() {
	ed := NewEventDispatcher()
	event := &TestEvent{Name: "test", Payload: "test"}
	eh := &MockHandler{}
	eh.On("Handle", event)
	ed.Register("test", eh)
	ed.Dispatch(event)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatecherTestSuite))
}
