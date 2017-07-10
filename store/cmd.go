package store

import (
	"database/sql"

	"github.com/arahmanhamdy/elvish/store/storedefs"
)

func init() {
	initDB["initialize command history table"] = func(db *sql.DB) error {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS cmd (content text)`)
		return err
	}
}

// NextCmdSeq returns the next sequence number of the command history.
func (s *Store) NextCmdSeq() (int, error) {
	row := s.db.QueryRow(`SELECT ifnull(max(rowid), 0) + 1 FROM cmd`)
	var seq int
	err := row.Scan(&seq)
	return seq, err
}

// AddCmd adds a new command to the command history.
func (s *Store) AddCmd(cmd string) (int, error) {
	r, err := s.db.Exec(`INSERT INTO cmd (content) VALUES(?)`, cmd)
	if err != nil {
		return -1, err
	}
	i, err := r.LastInsertId()
	return int(i), err
}

// Cmd queries the command history item with the specified sequence number.
func (s *Store) Cmd(seq int) (string, error) {
	row := s.db.QueryRow(`SELECT content FROM cmd WHERE rowid = ?`, seq)
	var cmd string
	err := row.Scan(&cmd)
	return cmd, err
}

// IterateCmds iterates all the commands in the specified range, and calls the
// callback with the content of each command sequentially.
func (s *Store) IterateCmds(from, upto int, f func(string) bool) error {
	rows, err := s.db.Query(`SELECT content FROM cmd WHERE rowid >= ? AND rowid < ?`, from, upto)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cmd string
		err = rows.Scan(&cmd)
		if err != nil {
			break
		}
		if !f(cmd) {
			break
		}
	}
	return err
}

// Cmds returns the contents of all commands within the specified range.
func (s *Store) Cmds(from, upto int) ([]string, error) {
	var cmds []string
	err := s.IterateCmds(from, upto, func(cmd string) bool {
		cmds = append(cmds, cmd)
		return true
	})
	return cmds, err
}

// NextCmd finds the first command after the given sequence number (inclusive)
// with the given prefix.
func (s *Store) NextCmd(from int, prefix string) (int, string, error) {
	row := s.db.QueryRow(`SELECT rowid, content FROM cmd WHERE rowid >= ? AND substr(content, 1, ?) = ? ORDER BY rowid asc LIMIT 1`, from, len(prefix), prefix)
	return convertCmd(row)
}

// PrevCmd finds the last command before the given sequence number (exclusive)
// with the given prefix.
func (s *Store) PrevCmd(upto int, prefix string) (int, string, error) {
	var upto64 = int64(upto)
	if upto < 0 {
		upto64 = 0x7FFFFFFFFFFFFFFF
	}
	row := s.db.QueryRow(`SELECT rowid, content FROM cmd WHERE rowid < ? AND substr(content, 1, ?) = ? ORDER BY rowid DESC LIMIT 1`, upto64, len(prefix), prefix)
	return convertCmd(row)
}

func convertCmd(row *sql.Row) (seq int, cmd string, err error) {
	err = row.Scan(&seq, &cmd)
	if err == sql.ErrNoRows {
		err = storedefs.ErrNoMatchingCmd
	}
	return
}
