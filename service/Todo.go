package service

import (
	"database/sql"
	"goapi_base/config"
	"goapi_base/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func TodoAdd(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	mdl := &model.Todo{}
	c.Bind(mdl)
	statement, _ := db.Prepare("INSERT INTO Todo (Icerik, Tamamlandi) VALUES (?, ?)")
	statement.Exec(mdl.Icerik, mdl.Tamamlandi)
	defer statement.Close()
	return c.JSON(http.StatusCreated, mdl)
}

func TodoList(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	rows, _ := db.Query("SELECT TodoID, Icerik, Tamamlandi FROM Todo")
	defer rows.Close()
	mdl := []model.Todo{}
	for rows.Next() {
		item := model.Todo{}
		rows.Scan(&item.TodoID, &item.Icerik, &item.Tamamlandi)
		mdl = append(mdl, item)
	}
	return c.JSON(http.StatusOK, mdl)
}

func TodoGet(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	mdl := model.Todo{}
	statement, _ := db.Prepare("SELECT TodoID, Icerik, Tamamlandi FROM Todo WHERE TodoID = ?")
	err = statement.QueryRow(id).Scan(&mdl.TodoID, &mdl.Icerik, &mdl.Tamamlandi)
	defer statement.Close()
	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, mdl)
}

func TodoDelete(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	statement, _ := db.Prepare("DELETE FROM Todo WHERE TodoID = ?")
	statement.Exec(id)
	defer statement.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func TodoSet(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	mdl := &model.Todo{}
	c.Bind(mdl)
	statement, _ := db.Prepare("UPDATE Todo SET Icerik = ?, Tamamlandi = ? WHERE TodoID = ?")
	statement.Exec(mdl.Icerik, mdl.Tamamlandi, id)
	defer statement.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
