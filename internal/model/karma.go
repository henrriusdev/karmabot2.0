package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Karma struct {
	User        int64
	FirstName   *string
	LastName    *string
	Count       int
	LastUpdated time.Time
}

type KarmaModel struct {
	DB *sql.DB
}

// CREATE & INSERT METHODS

func (m *KarmaModel) CreateTable(channel string) error {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(user_id BIGINT NOT NULL PRIMARY KEY, karma INT NOT NULL, first_name VARCHAR(100), last_name VARCHAR(100), last_updated DATETIME NOT NULL)", channel)
	if _, err := m.DB.Exec(query); err != nil {
		return err
	}

	infoLog.Println("Created table for: ", channel, "group")

	return nil
}

func (m KarmaModel) InsertUsers(userID int64, firstName, lastName, channel string) error {
	query := fmt.Sprintf("INSERT INTO `%s`(user_id, karma, first_name, last_name, last_updated) VALUES (?, ?, ?, ?, ?)", channel)
	if _, err := m.DB.Exec(query, userID, 0, firstName, lastName, time.Now()); err != nil {
		return err
	}
	return nil
}

// GET methods

func (m *KarmaModel) GetKarmas(channel string, top bool) ([]*Karma, error) {
	var query string
	if top {
		query = fmt.Sprintf("SELECT user_id, first_name, last_name, karma FROM `%s` ORDER BY karma DESC LIMIT 10", channel)
	} else {
		query = fmt.Sprintf("SELECT user_id, first_name, last_name, karma FROM `%s` ORDER BY karma ASC LIMIT 10", channel)
	}

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	karmas := []*Karma{}

	for rows.Next() {
		k := &Karma{}
		err := rows.Scan(&k.User, &k.FirstName, &k.LastName, &k.Count)
		if err != nil {
			return nil, err
		}

		karmas = append(karmas, k)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return karmas, nil
}

func (m *KarmaModel) GetActualKarma(userID int64, channel string) (int, error) {
	var user Karma
	query := fmt.Sprintf("SELECT karma FROM `%s` WHERE user_id = ?", channel)

	if err := m.DB.QueryRow(query, userID).Scan(&user.Count); err != nil {
		return 0, err
	}

	return user.Count, nil
}

func (m *KarmaModel) CanModify(userID int64, firstName, lastName, channel string) bool {
	var user Karma
	if channel == "" {
		return false
	}

	query := fmt.Sprintf("SELECT last_updated FROM `%s` WHERE user_id = ?", channel)

	if err := m.DB.QueryRow(query, userID).Scan(&user.LastUpdated); err != nil {
		if err == sql.ErrNoRows {
			err = m.InsertUsers(userID, firstName, lastName, channel)
			if err != nil {
				fmt.Println("105", err)
				return false
			}

			return true
		}
		fmt.Println("116", err)
		return false
	}
	timeElapsed := time.Since(user.LastUpdated)
	fmt.Println(timeElapsed)
	return timeElapsed > time.Second*60
}

// UPDATE METHODS

func (m *KarmaModel) AddKarma(karmaTransmitter, karmaReceiver int64, firstNameReceiver, lastNameReceiver, channel string) error {
	query := fmt.Sprintf("UPDATE `%s` SET karma = ?, first_name = ?, last_name = ? WHERE user_id = ?", channel)

	karma, err := m.GetActualKarma(karmaReceiver, channel)
	if err != nil {
		err := m.InsertUsers(karmaReceiver, firstNameReceiver, lastNameReceiver, channel)
		if err != nil {
			return err
		}
	}

	karma++
	_, err = m.DB.Exec(query, karma, firstNameReceiver, lastNameReceiver, karmaReceiver)
	if err != nil {
		return err
	}

	err = m.updateLastKarma(time.Now(), channel, karmaTransmitter)
	if err != nil {
		return err
	}

	return nil
}

func (m *KarmaModel) SubstractKarma(karmaTransmitter, karmaReceiver int64, firstNameReceiver, lastNameReceiver, channel string) error {
	query := fmt.Sprintf("UPDATE `%s` SET karma = ?, first_name = ?, last_name = ? WHERE user_id = ?", channel)

	karma, err := m.GetActualKarma(karmaReceiver, channel)
	if err != nil {
		err := m.InsertUsers(karmaReceiver, firstNameReceiver, lastNameReceiver, channel)
		if err != nil {
			return err
		}
	}

	karma--
	_, err = m.DB.Exec(query, karma, firstNameReceiver, lastNameReceiver, karmaReceiver)
	if err != nil {
		return err
	}

	err = m.updateLastKarma(time.Now(), channel, karmaTransmitter)
	if err != nil {
		return err
	}

	return nil
}

func (m *KarmaModel) updateLastKarma(date time.Time, channel string, username int64) error {
	query := fmt.Sprintf("UPDATE `%s` SET last_updated = ? WHERE user_id = ?", channel)
	_, err := m.DB.Exec(query, date, username)
	if err != nil {
		return err
	}

	return nil
}
