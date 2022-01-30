package model

import (
    "database/sql"
    "goapi_base/config"

    _ "github.com/mattn/go-sqlite3"
)

type Todo struct {
    TodoID int    `json:"TodoID"`
    Icerik string `json:"Icerik"`
    Tamamlandi string `json:"Tamamlandi"`
}

func TodoCreateTable() {
    db, _ := sql.Open("sqlite3", config.DB_NAME)
    defer db.Close()
    statement, _ := db.Prepare(`
        CREATE TABLE IF NOT EXISTS Todo
        (
            TodoID INTEGER PRIMARY KEY,
            Icerik TEXT,
            Tamamlandi TEXT
        )
    `)
    statement.Exec()
    defer statement.Close()
}
