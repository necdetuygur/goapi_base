generator("Todo", ["Icerik", "Tamamlandi"]);

function generator(table, cols) {
  structRows = [];
  tableRows = [];
  insertCols = [];
  updateCols = [];
  questionMarks = [];
  insertMdls = [];
  scanItems = [];
  scanMdls = [];
  cols.forEach((item) => {
    structRows.push(`
    ${item} string \`json:"${item}"\``);

    tableRows.push(`
            ${item} TEXT`);

    insertCols.push(item);

    updateCols.push(`${item} = ?`);

    questionMarks.push("?");

    insertMdls.push(`mdl.${item}`);

    scanItems.push(`&item.${item}`);

    scanMdls.push(`&mdl.${item}`);
  });
  structRows = structRows.join("");
  tableRows = tableRows.join(",");
  insertCols = insertCols.join(", ");
  updateCols = updateCols.join(", ");
  questionMarks = questionMarks.join(", ");
  insertMdls = insertMdls.join(", ");
  scanItems = scanItems.join(", ");
  scanMdls = scanMdls.join(", ");
  /**/
  var model = `package model

import (
    "database/sql"
    "goapi_base/config"

    _ "github.com/mattn/go-sqlite3"
)

type ${table} struct {
    ${table}ID int    \`json:"${table}ID"\`${structRows}
}

func ${table}CreateTable() {
    db, _ := sql.Open("sqlite3", config.DB_NAME)
    defer db.Close()
    statement, _ := db.Prepare(\`
        CREATE TABLE IF NOT EXISTS ${table}
        (
            ${table}ID INTEGER PRIMARY KEY,${tableRows}
        )
    \`)
    statement.Exec()
    defer statement.Close()
}
`;
  var router = `package router

import (
	"goapi_base/service"

	"github.com/labstack/echo/v4"
)

func ${table}Router(e *echo.Group) {
	e.POST("/${table}", service.${table}Add)
	e.GET("/${table}", service.${table}List)
	e.GET("/${table}/:id", service.${table}Get)
	e.DELETE("/${table}/:id", service.${table}Delete)
	e.PUT("/${table}/:id", service.${table}Set)
}
`;
  var service = `package service

import (
	"database/sql"
	"goapi_base/config"
	"goapi_base/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ${table}Add(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	mdl := &model.${table}{}
	c.Bind(mdl)
	statement, _ := db.Prepare("INSERT INTO ${table} (${insertCols}) VALUES (${questionMarks})")
	statement.Exec(${insertMdls})
	defer statement.Close()
	return c.JSON(http.StatusCreated, mdl)
}

func ${table}List(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	rows, _ := db.Query("SELECT ${table}ID, ${insertCols} FROM ${table}")
	defer rows.Close()
	mdl := []model.${table}{}
	for rows.Next() {
		item := model.${table}{}
		rows.Scan(&item.${table}ID, ${scanItems})
		mdl = append(mdl, item)
	}
	return c.JSON(http.StatusOK, mdl)
}

func ${table}Get(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	mdl := model.${table}{}
	statement, _ := db.Prepare("SELECT ${table}ID, ${insertCols} FROM ${table} WHERE ${table}ID = ?")
	err = statement.QueryRow(id).Scan(&mdl.${table}ID, ${scanMdls})
	defer statement.Close()
	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, mdl)
}

func ${table}Delete(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	statement, _ := db.Prepare("DELETE FROM ${table} WHERE ${table}ID = ?")
	statement.Exec(id)
	defer statement.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func ${table}Set(c echo.Context) error {
	db, _ := sql.Open("sqlite3", config.DB_NAME)
	defer db.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}
	mdl := &model.${table}{}
	c.Bind(mdl)
	statement, _ := db.Prepare("UPDATE ${table} SET ${updateCols} WHERE ${table}ID = ?")
	statement.Exec(${insertMdls}, id)
	defer statement.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
`;

  const fs = require("fs");

  fs.mkdir("model", (err) => {
    if (err) {
      return console.error(err);
    }
    console.log("model klasör oluşturuldu.");
  });

  fs.mkdir("router", (err) => {
    if (err) {
      return console.error(err);
    }
    console.log("router klasör oluşturuldu.");
  });

  fs.mkdir("service", (err) => {
    if (err) {
      return console.error(err);
    }
    console.log("service klasör oluşturuldu.");
  });

  fs.writeFile(`./model/${table}.go`, model, function (e) {
    if (e) throw e;
    console.log(`${table} model oluşturuldu.`);
  });
  fs.writeFile(`./router/${table}.go`, router, function (e) {
    if (e) throw e;
    console.log(`${table} router oluşturuldu.`);
  });
  fs.writeFile(`./service/${table}.go`, service, function (e) {
    if (e) throw e;
    console.log(`${table} service oluşturuldu.`);
  });
}
