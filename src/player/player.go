package player

import "net"
import "time"
import "encoding/binary"

type UserData struct {
	mq chan string;
	id int;
	name string;
    owner_id int;
	is_capital int;
	achievement int;
	x int;
	y int;
	gold int;
	wood int;
	food int;
	iron int;
	stone int;
	scout int;
	swordsman int;
	crossbow_archer int;
	squire int;
	templar int;
	paladin int;
	archer_cavalry int;
	royal_knight int;
	action_events_count int;
	recruit_events_count int;
	deals_count int;
	lock_version int;
	last_move_time int;
	durability int;
	arcane_mage int;
	battle_mage  int;
	holy_mage int;
	is_auto_fix int;
	revive_time int;
	item_warehouse_lv int;
	item_transport_lv int;
	skeleton int;
	ghost_rider int;
	ram int;
	zeppelin int;
	steel_golem int;
	cruiser int
}


func send(conn net.Conn, p string) error {
	header := make([]byte,2)
	binary.BigEndian.PutUint16(header, uint16(len(p)));
	_, err := conn.Write(header)
	if err != nil {
		println("Error send reply header:", err.Error())
		return err
	}

	_, err = conn.Write([]byte(p))
	if err != nil {
		println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}

func (user *UserData) flush_timer() {
	for {
		time.Sleep(10*time.Second)
		if user.id != 0 {
			DB.Flush(user)
		}
		time.Sleep(4*time.Second)
	}
}

func NewPlayer(in chan string, conn net.Conn) {
	var user UserData
	user.mq = make(chan string, 100)

	if send(conn, "Welcome") != nil {
		return
	}

	go user.flush_timer()
L:
	for {
		select {
		case msg := <-in:
			if msg == "" {
				break L
			}

			result := user.exec_cli(msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg := <-user.mq:
			if msg == "" {
				break L
			}

			result := user.exec_srv(msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}
		}
	}
}
