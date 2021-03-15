package main

//go:generate npx prisma generate

import (
	"context"
	"log"
	"mime"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mattn/go-prisma-example/prisma/db"
)

func main() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	mime.AddExtensionType(".js", "application/javascript")

	e := echo.New()

	e.POST("/tasks", func(c echo.Context) error {
		var task db.TaskModel
		if err := c.Bind(&task); err != nil {
			c.Logger().Error("Bind: ", err)
			return c.String(http.StatusBadRequest, "Bind: "+err.Error())
		}
		text := &task.Text
		completed := &task.Completed
		newTask, err := client.Task.CreateOne(
			db.Task.Text.SetIfPresent(text),
			db.Task.Completed.SetIfPresent(completed),
		).Exec(context.Background())
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, newTask)
	})

	e.GET("/tasks", func(c echo.Context) error {
		tasks, err := client.Task.FindMany().OrderBy(
			db.Task.ID.Order(db.ASC),
		).Exec(context.Background())
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, tasks)
	})

	e.POST("/tasks/:id", func(c echo.Context) error {
		var task db.TaskModel
		if err := c.Bind(&task); err != nil {
			c.Logger().Error("Bind: ", err)
			return c.String(http.StatusBadRequest, "Bind: "+err.Error())
		}
		text := &task.Text
		completed := &task.Completed
		newTask, err := client.Task.FindUnique(
			db.Task.ID.Equals(task.ID),
		).Update(
			db.Task.Text.SetIfPresent(text),
			db.Task.Completed.SetIfPresent(completed),
		).Exec(context.Background())
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, newTask)
	})

	e.DELETE("/tasks/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		task, err := client.Task.FindUnique(
			db.Task.ID.Equals(id),
		).Delete().Exec(context.Background())
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, task)
	})
	e.GET("/tasks/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		task, err := client.Task.FindUnique(
			db.Task.ID.Equals(id),
		).Exec(context.Background())
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusOK, task)
	})

	e.Static("/", "assets")
	e.Logger.Fatal(e.Start(":8989"))
}
