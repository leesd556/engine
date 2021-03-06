/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package pbft

import (
	"errors"
	"sync"
)

var ErrInvalidSave = errors.New("Invalid Save Error")
var ErrEmptyRepo = errors.New("Repository has empty state")

type StateRepository struct {
	state State
	sync.RWMutex
}

func NewStateRepository() StateRepository {
	return StateRepository{
		state:   State{},
		RWMutex: sync.RWMutex{},
	}
}
func (repo *StateRepository) Save(state State) error {

	repo.Lock()
	defer repo.Unlock()
	id := repo.state.StateID.ID
	if id == state.StateID.ID || id == "" {
		repo.state = state
		return nil
	}
	return ErrInvalidSave
}
func (repo *StateRepository) Load() (State, error) {

	if repo.state.StateID.ID == "" {
		return repo.state, ErrEmptyRepo
	}

	return repo.state, nil
}

func (repo *StateRepository) Remove() {
	repo.state = State{}
}
