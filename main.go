package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cristalhq/acmd"
	_ "github.com/mattn/go-sqlite3"
)

func dbCheckOrCreate() (error error) {
	_, err := os.Stat(os.Getenv("HOME") + "/.dirnote/dirnotes.sqlite")
	if os.IsNotExist(err) {
		fmt.Print("db file does not exist, creating...")

		_, err := os.Stat(("HOME") + "/.dirnote/")
		if os.IsNotExist(err) {
			err = os.Mkdir(os.Getenv("HOME")+"/.dirnote/", 0700)
			if err != nil {
				fmt.Print(err)
			}
		}

		_, err = os.Create(os.Getenv("HOME") + "/.dirnote/dirnotes.sqlite")
		if err != nil {
			if os.IsNotExist(err) {
				log.Print(err)
			}

			if err == os.ErrExist {
				log.Print(err)
			} else {
				log.Fatal(err)
			}
		}
		err = dbPrep()
		if err != nil {
			log.Fatal(err)
		}
	}
	//err = nil
	return err
}
func dbPrep() (err error) {
	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.dirnote/dirnotes.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	//dbPrepTables(db)

	sqlStmt := `
				CREATE TABLE IF NOT EXISTS directories 
				(
					id integer not null
						primary key,
					dirname text);
				`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	sqlStmt = `
				CREATE TABLE IF NOT EXISTS notes 
				(
					id integer not null 
						primary key,
					directory int
						constraint notes_directories_id_fk
							references directories,
					note_data TEXT
					);
				`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	//db.Close()
	return err
}

func checkOrAddDir(db *sql.DB, directory string) (directoryResult int, directoryName string) {
	selectstmt, err := db.Prepare("select id, dirname from directories where dirname = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer selectstmt.Close()

	var dir int
	var dirname string

	err = selectstmt.QueryRow(directory).Scan(&dir, &dirname)
	if errors.Is(err, sql.ErrNoRows) {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		insertstmt, err := tx.Prepare("INSERT INTO directories(dirname) VALUES (?)")
		if err != nil {
			log.Fatal(err)
		}

		defer insertstmt.Close()

		_, err = insertstmt.Exec(directory)
		if err != nil {
			log.Fatal(err)
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

//		err = selectstmt.QueryRow(directory).Scan(&dir, &dirname)
//		if err != nil {
//			log.Fatal(err)
//		}

		//err = nil
	} else {
		if err != nil {
			log.Fatal(err)
		}
	}

//	err = dbPrep()
//	if err != nil {
//		log.Println(err)
//	}
	fmt.Println(directory + " added")
	return dir, dirname
}

func addNote(db *sql.DB, dir int, noteData string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	insertstmt, err := tx.Prepare("INSERT INTO notes(directory, note_data) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	defer insertstmt.Close()

	_, err = insertstmt.Exec(dir, noteData)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	//err = selectstmt.QueryRow(noteData).Scan(&dir)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func findDirID(db *sql.DB, path string) (id int, err error) {
	selectstmt, err := db.Prepare("select id from directories where dirname = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer selectstmt.Close()

	var dir int
	err = selectstmt.QueryRow(path).Scan(&dir)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Directory not in dirnote")
		}else{
			log.Fatal(err)
		}
	}
	//log.Print(dir)
	return dir, err
}
func getNotesForDirectory(db *sql.DB, noteDirectory int) (err error) {
	rows, err := db.Query("select id, note_data from notes where directory = ?", noteDirectory)
	if err != nil {
			log.Fatal(err)
	}
	defer rows.Close()

	var id int
	var note_data string
	for rows.Next() {
		err := rows.Scan(&id, &note_data)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
			fmt.Print("No notes for current directory")
			} else {
				log.Fatal(err)
			}
		}
		fmt.Printf("Note: %d -- %s \n", id, note_data)
		//log.Println(id, note_data)
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {} else {
				log.Fatal(err)
			}
	}
	
	return err
}

func deleteNoteByID(db *sql.DB, ID int) (err error) {
	res, err := db.Exec("DELETE FROM notes WHERE id=$1", ID)

	if err == nil {

		count, err := res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		println(count)
	}

	return err

}

func main() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	err = dbCheckOrCreate()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.dirnote/dirnotes.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()



	cmds := []acmd.Command{
		{
			Name:        "add",
			Description: "adds a new note in the current dir",
			ExecFunc: func(ctx context.Context, args []string) error {
				fmt.Println(path)
				checkOrAddDir(db, path)
				dir, err := findDirID(db, path)
				//fmt.Println(dir)
				if err != nil {
					log.Fatal(err)
				}

				err = addNote(db, dir, args[0])
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:        "del",
			Description: "deletes a note by global ID",
			ExecFunc: func(ctx context.Context, args []string) error {

				i, err := strconv.Atoi(args[0])
				if err != nil {
					log.Fatal(err)
				}

				err = deleteNoteByID(db, i)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:        "get",
			Description: "gets notes for the dir",
			ExecFunc: func(ctx context.Context, args []string) error {
				if len(args) > 0 {
					path = args[0]
				}
				fmt.Println(path)

				id, err := findDirID(db, path)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
					} else {
						//log.Fatal(err)
					}
				}
				// fmt.Print(id)
				err = getNotesForDirectory(db, id)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
					} else {
						//log.Fatal(err)
					}
				}
				return nil
			},
		},
	}
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "dirnote",
		AppDescription: "An app to keep notes for your dirs",
		Version:        "the best v0.0.2",
		// Context - if nil `signal.Notify` will be used
		// Args - if nil `os.Args[1:]` will be used
		// Usage - if nil default print will be used
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}
