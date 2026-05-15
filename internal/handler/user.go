package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"user/internal/model"
)

// UserHandler 负责处理用户相关的 HTTP 请求。
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler 创建用户处理器，并注入数据库连接。
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUsers 查询所有用户
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []model.User

	// 查询 users 表中的所有用户记录。
	err := h.db.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "查询用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "查询成功",
		"data": users,
	})
}

// GetUserByID 查询单个用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	// 按路径参数 id 查询用户，查不到时返回 404。
	err := h.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "查询用户失败",
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "查询成功",
		"data": user,
	})
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User

	// 将请求体中的 JSON 数据绑定到用户模型。
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "请求参数错误",
			"error": err.Error(),
		})
		return
	}

	// 做最基础的必填字段校验，避免空数据入库。
	if user.Name == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "名称和邮箱不能为空",
		})
		return
	}

	// 保存新用户记录。
	err := h.db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "创建用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": user,
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	// 更新前先确认用户存在，避免对不存在的记录执行更新。
	if err := h.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "查询用户失败",
				"error": err.Error(),
			})
		}
		return
	}

	// 将新的用户信息绑定到已有模型上。
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "请求参数错误",
			"error": err.Error(),
		})
		return
	}

	// 按 id 更新用户字段。
	err := h.db.Model(&user).Where("id = ?", id).Updates(user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "更新用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新成功",
		"data": user,
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	// 删除前先查一次，方便区分“用户不存在”和“删除失败”。
	if err := h.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "查询用户失败",
				"error": err.Error(),
			})
		}
		return
	}

	// 按 id 删除用户记录。
	err := h.db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "删除用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
