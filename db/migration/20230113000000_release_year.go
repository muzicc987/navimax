package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddRelRecYear, downAddRelRecYear)
}

func upAddRelRecYear(tx *sql.Tx) error {
	_, err := tx.Exec(`
alter table media_file
    add date date null default null;
alter table media_file
    add release_date date null default null;
alter table media_file
    add release_year integer;
	
alter table album
    add date date null default null;
alter table album
    add release_date date null default null;
alter table album
    add min_release_year integer;
alter table album
    add max_release_year integer;
create index if not exists media_file_track_number
	on media_file (release_year, disc_number, track_number);
`)
	if err != nil {
		return err
	}

	notice(tx, "A full rescan needs to be performed to import more tags")
	return forceFullRescan(tx)
}

func downAddRelRecYear(tx *sql.Tx) error {
	return nil
}
