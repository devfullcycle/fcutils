package events

import (
	"errors"
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

func (h *TestEventHandler) Handle(event EventInterface,  errs chan error) {
	
}

type EventDispatecherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher EventDispatcher
}

func (suite *EventDispatecherTestSuite) SetupTest() {
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.eventDispatcher = *NewEventDispatcher()
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Register() {
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)	
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))


	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])
	assert.Equal(suite.T(), &suite.handler3, suite.eventDispatcher.handlers[suite.event2.GetName()][0])
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Clear() {
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	suite.eventDispatcher.Clear()
	
	assert.Equal(suite.T(), 0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Remove() {
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	assert.Equal(suite.T(), 1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	assert.Equal(suite.T(), 0, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	assert.Equal(suite.T(), 0, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Has() {
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	assert.Equal(suite.T(), true, suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.Equal(suite.T(), true, suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	assert.Equal(suite.T(), true, suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler3))
	assert.Equal(suite.T(), false, suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, errs chan error) {
	m.Called(event, errs)
}

func (suite *EventDispatecherTestSuite) TestEventDispatcher_Dispatch() {
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	errsChan := make(chan error)
	
	suite.eventDispatcher.Dispatch(&suite.event, errsChan)

	assert.Equal(suite.T(), 0, len(errsChan))
}

type HandlerWithError struct {
	ID int
}

func (h *HandlerWithError) Handle(event EventInterface, errs chan error) {
	errs <- errors.New("error")
}


func (suite *EventDispatecherTestSuite) TestEventDispatcher_DispatchWithError() {
	handler1 := HandlerWithError{ID: 1}
	handler2 := HandlerWithError{ID: 2}
	handler3 := HandlerWithError{ID: 3}
	

	suite.eventDispatcher.Register(suite.event.GetName(), &handler1)
	suite.eventDispatcher.Register(suite.event.GetName(), &handler2)
	suite.eventDispatcher.Register(suite.event.GetName(), &handler3)

	errsChan := make(chan error)
	
	suite.eventDispatcher.Dispatch(&suite.event, errsChan)

	time.Sleep(3 * time.Second)
	assert.Equal(suite.T(), 3, len(errsChan))

	// for err := range errsChan {
	// 	assert.Equal(suite.T(), "error", err.Error())
	// }
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatecherTestSuite))
}
