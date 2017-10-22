package data

import "fmt"

type table struct {
	name    string
	size    int
	players []Player
	admin   int64
	mc      chan string
}

type Player struct {
	id    int64
	name  string
	money int64
}

var tm map[string]table

func Init() {
	tm = make(map[string]table)

	// load from disk

}

func JoinTable(tableName string, player Player) error {
	// check to make sure table exists
	if t, ok := tm[tableName]; !ok {
		// does not exist
		return fmt.Errorf("table named %s does not exist", tableName)
	} else if t.size == len(t.players) {
		return fmt.Errorf("table named %s is full", tableName)
	} else {
		t.players = append(t.players, player)
	}

	return nil
}

func NewTable(tableName string, size int, admin int64) error {
	// ensure that table does not exist
	if _, ok := tm[tableName]; ok {
		return fmt.Errorf("table named %s already exists", tableName)
	}

	t := table{tableName, size, make([]Player, 9, 9), admin, make(chan string)}
	tm[tableName] = t
	// start message receiving service
	go ReceiveMessages(t)
	return nil
}

func GetMessageChan(tableName string, player int64) (chan string, error) {
	// check to make sure table exists
	if t, ok := tm[tableName]; !ok {
		return nil, fmt.Errorf("table named %s does not exist", tableName)
	} else {
		// make sure player is a member of group
		isMember := func() bool {
			for _, p := range t.players {
				if player == p.id {
					return true
				}
			}
			if player == t.admin {
				return true
			}
			return false
		}()

		if !isMember {
			return nil, fmt.Errorf("player %v is not a member of table %d", player, tableName)
		}

		return t.mc, nil

	}
}

func MessageTable(tableName string, player Player, message string) error {
	// check to make sure table exists
	if t, ok := tm[tableName]; !ok {
		return fmt.Errorf("table named %s does not exist", tableName)
	} else {
		// make sure player is a member of group
		isMember := func() bool {
			for _, p := range t.players {
				if player.id == p.id {
					return true
				}
			}
			if player.id == t.admin {
				return true
			}
			return false
		}()

		if !isMember {
			return fmt.Errorf("player %v is not a member of table %d", player, tableName)
		}

		// send message through channel
		t.mc <- player.name + ": " + message
	}
	return nil
}

// TODO may have to use struct pointer
func ReceiveMessages(t table) error {
	return nil
}
