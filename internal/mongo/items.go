package mongo

import (
	"fmt"
)

type Item struct {
}

func ItemsForPlayerUuid(uuid string) ([]*Item, error) {
	return nil, fmt.Errorf("not implemented")
}

func ItemsForProfileUuid(uuid string) ([]*Item, error) {
	return nil, fmt.Errorf("not implemented")
}
