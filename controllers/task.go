package controllers

import (
	"github.com/Temirlan-Temirbaev/tz-golang/config"
	"github.com/Temirlan-Temirbaev/tz-golang/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InitTaskRoutes(e *echo.Echo) {
	e.POST("/tasks", func(c echo.Context) error {
		var task models.Task
		if err := c.Bind(&task); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if err := config.DB.Create(&task).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create task"})
		}

		return c.JSON(http.StatusCreated, task)
	})
	e.GET("/tasks", func(c echo.Context) error {
		var tasks []models.Task
		if err := config.DB.Find(&tasks).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tasks"})
		}

		return c.JSON(http.StatusOK, tasks)
	})
	e.POST("/tasks/:task_id/:telegram_id", func(c echo.Context) error {
		taskID := c.Param("task_id")
		telegramID := c.Param("telegram_id")
		var task models.Task
		if err := config.DB.First(&task, taskID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
		}
		var user models.User
		if err := config.DB.First(&user, "telegram_id = ?", telegramID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err := config.DB.Model(&task).Association("Users").Append(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add user to task"})
		}

		if err := config.DB.Model(&user).Association("Tasks").Append(&task); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add task to user"})
		}

		if err := config.DB.Save(&task).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update task"})
		}

		return c.JSON(http.StatusOK, task)
	})
}
