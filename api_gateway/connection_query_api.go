package api_gateway

import (
	"errors"
	"sync"

	"github.com/it-chain/engine/common/event"
)

var ErrConnectionExists = errors.New("Connection already exists")

type ConnectionQueryApi struct {
	mux                  *sync.Mutex
	connectionRepository *ConnectionRepository
}

func NewConnectionQueryApi(connRepo *ConnectionRepository) *ConnectionQueryApi {
	return &ConnectionQueryApi{
		mux:                  &sync.Mutex{},
		connectionRepository: connRepo,
	}
}

func (q ConnectionQueryApi) GetAllConnectionList() ([]Connection, error) {
	return q.connectionRepository.FindAll()
}

func (q ConnectionQueryApi) GetConnectionByID(connID string) (Connection, error) {
	return q.connectionRepository.FindByID(connID)
}

type ConnectionRepository struct {
	mux             *sync.RWMutex
	ConnectionTable map[string]Connection
}

func NewConnectionRepository() *ConnectionRepository {
	return &ConnectionRepository{
		mux:             &sync.RWMutex{},
		ConnectionTable: make(map[string]Connection),
	}
}

func (cr *ConnectionRepository) Save(conn Connection) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	_, exist := cr.ConnectionTable[conn.ConnectionID]
	if exist {
		return ErrConnectionExists
	}

	cr.ConnectionTable[conn.ConnectionID] = conn

	return nil
}

func (cr *ConnectionRepository) Remove(connID string) error {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	delete(cr.ConnectionTable, connID)

	return nil
}

func (cr *ConnectionRepository) FindAll() ([]Connection, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	connectionList := []Connection{}

	for _, conn := range cr.ConnectionTable {
		connectionList = append(connectionList, conn)
	}

	return connectionList, nil
}

func (cr *ConnectionRepository) FindByID(connID string) (Connection, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	for _, conn := range cr.ConnectionTable {
		if connID == conn.ConnectionID {
			return conn, nil
		}
	}

	return Connection{}, nil
}

type ConnectionEventListener struct {
	connectionRepository *ConnectionRepository
}

func NewConnectionEventListener(connRepo *ConnectionRepository) *ConnectionEventListener {
	return &ConnectionEventListener{
		connectionRepository: connRepo,
	}
}

func (cel *ConnectionEventListener) HandleConnectionCreatedEvent(event event.ConnectionCreated) error {
	connection := Connection{
		ConnectionID: event.ConnectionID,
		Address:      event.Address,
	}

	err := cel.connectionRepository.Save(connection)
	if err != nil {
		return err
	}
	return nil
}

func (cel *ConnectionEventListener) HandleConnectionClosedEvent(event event.ConnectionClosed) error {
	cel.connectionRepository.Remove(event.ConnectionID)
	return nil
}

type Connection struct {
	ConnectionID string
	Address      string
}
