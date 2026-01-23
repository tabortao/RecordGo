package router

import (
	"path/filepath"
	"recordgo/internal/config"
	"recordgo/internal/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 中文注释：初始化 Gin 路由，注册健康检查与未来的模块路由
func New(cfg *config.Config, lg *zap.Logger) *gin.Engine {
	if cfg.Env == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allow := ""
		for _, o := range cfg.AllowedOrigins {
			if o == "*" {
				allow = "*"
				break
			}
			if origin == o {
				allow = origin
				break
			}
		}
		if allow == "" && len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
			allow = "*"
		}
		if allow != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allow)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 中文注释：静态文件服务（上传图片），映射到 /api/uploads
	r.Static("/api/uploads", filepath.Join(cfg.StorageRoot, "uploads"))

	api := r.Group("/api")
	{
		// 中文注释：认证路由（登录/注册）
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)
		// 中文注释：子账号令牌登录（免密码，仅限子账号）
		api.POST("/auth/token-login", handlers.TokenLogin)
		// 中文注释：账号安全与资料
		api.POST("/auth/change-password", handlers.ChangePassword)
		api.GET("/user/profile", handlers.GetProfile)
		api.PUT("/user/profile", handlers.UpdateProfile)

		api.GET("/health", handlers.Health)
		// 中文注释：默认心愿占位接口，可用于前端初始化
		api.GET("/wishes/default", handlers.DefaultWishes)

		// 中文注释：心愿管理与兑换记录相关接口
		api.GET("/wishes", handlers.ListWishes)
		api.POST("/wishes", handlers.CreateWish)
		api.GET("/wishes/:id", handlers.GetWish)
		api.PUT("/wishes/:id", handlers.UpdateWish)
		api.DELETE("/wishes/:id", handlers.DeleteWish)
		api.POST("/wishes/:id/exchange", handlers.ExchangeWish)
		api.GET("/wishes/records", handlers.ListWishRecords)
		api.POST("/upload/wish-icon", handlers.UploadWishIcon)
		// 中文注释：任务图片与音频上传接口
		api.POST("/upload/task-image", handlers.UploadTaskImage)
		api.POST("/upload/task-audio", handlers.UploadTaskAudio)
		// 中文注释：用户头像上传
		api.POST("/upload/avatar", handlers.UploadAvatar)
		api.POST("/upload/avatar/object", handlers.UploadAvatarObject)
		api.DELETE("/tasks/:id/image", handlers.DeleteTaskImage)
		api.GET("/storage/info", handlers.StorageInfo)
		api.POST("/storage/presign", handlers.StoragePresign)
		api.GET("/storage/presign-view", handlers.StoragePresignView)

		// 中文注释：任务管理 RESTful 路由
		api.POST("/tasks", handlers.CreateTask)
		api.POST("/tasks/batch", handlers.CreateTasksBatch)
		// 中文注释：AI 智能创建任务识别
		api.POST("/ai/parse-task", handlers.ParseTaskByAI)
		// 中文注释：OCR 代理接口（解决前端跨域问题）
		api.POST("/proxy/ocr", handlers.HandleOCRProxy)
		api.GET("/tasks", handlers.ListTasks)
		api.GET("/tasks/:id", handlers.GetTask)
		api.PUT("/tasks/:id", handlers.UpdateTask)
		api.PATCH("/tasks/:id/status", handlers.UpdateStatus)
		api.DELETE("/tasks/:id", handlers.DeleteTask)
		api.DELETE("/tasks", handlers.BatchDelete) // ?ids=1,2,3
		api.GET("/tasks/recycle-bin", handlers.ListRecycleBin)
		api.POST("/tasks/recycle-bin/restore", handlers.RestoreTasks) // ?ids=1,2
		api.POST("/tasks/:id/tomato/complete", handlers.CompleteTomato)

		// 中文注释：子账号管理（仅家长或具有权限的账号可操作）
		api.GET("/subaccounts", handlers.ListSubAccounts)
		api.POST("/subaccounts", handlers.CreateSubAccount)
		api.PUT("/subaccounts/:id", handlers.UpdateSubAccount)
		api.DELETE("/subaccounts/:id", handlers.DeleteSubAccount)
		api.POST("/subaccounts/:id/token", handlers.GenerateChildToken)
		// 中文注释：修复历史子账号未绑定父账号的情况（仅父账号/管理员）
		api.POST("/subaccounts/:id/rebind", handlers.RebindSubAccount)

		// 中文注释：金币查询（返回当前登录用户的有效金币，父子同步时返回父账号金币）
		api.GET("/coins", handlers.GetCoins)
		api.POST("/coins/set", handlers.SetCoins)

		api.GET("/task-categories", handlers.ListTaskCategories)
		api.PUT("/task-categories", handlers.PutTaskCategories)
		api.PATCH("/task-categories/order", handlers.PatchTaskCategoryOrder)

		api.GET("/tasks/occurrences", handlers.ListTaskOccurrences)
		api.POST("/tasks/:id/occurrences/complete", handlers.CompleteTaskOccurrence)
		api.POST("/tasks/:id/occurrences/uncomplete", handlers.UncompleteTaskOccurrence)
		api.POST("/tasks/:id/occurrences/delete", handlers.DeleteTaskOccurrence)

		api.POST("/data/clear", handlers.ClearData)

		// 管理员后台接口（仅用户ID=1允许）
		api.GET("/admin/overview", handlers.AdminUsersOverview)
		api.GET("/admin/users", handlers.AdminListUsers)
		api.POST("/admin/users/:id/vip", handlers.AdminUpdateVIP)
		api.POST("/admin/users/:id/reset-password", handlers.AdminResetPassword)
		api.POST("/admin/users/:id/disabled", handlers.AdminSetDisabled)
		api.DELETE("/admin/users/:id", handlers.AdminDeleteUser)

		// 听写大师模块
		dictation := api.Group("/dictation")
		{
			dictation.GET("/wordbanks", handlers.ListWordBanks)
			dictation.POST("/wordbanks", handlers.CreateWordBank)
			dictation.PUT("/wordbanks/:id", handlers.UpdateWordBank)
			dictation.DELETE("/wordbanks/:id", handlers.DeleteWordBank)

			dictation.GET("/settings", handlers.GetDictationSettings)
			dictation.PUT("/settings", handlers.UpdateDictationSettings)

			dictation.GET("/mistakes", handlers.ListMistakes)
			dictation.POST("/mistakes", handlers.AddMistake)
			dictation.DELETE("/mistakes/:id", handlers.DeleteMistake)

			dictation.POST("/history", handlers.AddDictationHistory)
			dictation.GET("/history", handlers.ListDictationHistory)
			dictation.GET("/stats", handlers.GetDictationStats)
		}

		// 小成长模块
		littleGrowth := api.Group("/little-growth")
		{
			littleGrowth.GET("/records", handlers.ListGrowthRecords)
			littleGrowth.GET("/records/:id", handlers.GetGrowthRecord) // New
			littleGrowth.POST("/records", handlers.CreateGrowthRecord)
			littleGrowth.PUT("/records/:id", handlers.UpdateGrowthRecord)
			littleGrowth.DELETE("/records/:id", handlers.DeleteGrowthRecord)
			littleGrowth.PATCH("/records/:id/pin", handlers.TogglePinGrowthRecord)
			littleGrowth.PATCH("/records/:id/favorite", handlers.ToggleFavoriteGrowthRecord)
			littleGrowth.POST("/records/:id/comments", handlers.AddGrowthComment)

			// Tags
			littleGrowth.GET("/tags", handlers.ListGrowthTags)
			littleGrowth.POST("/tags", handlers.CreateGrowthTag)
			littleGrowth.PUT("/tags/:id", handlers.UpdateGrowthTag)
			littleGrowth.DELETE("/tags/:id", handlers.DeleteGrowthTag)
		}
		api.POST("/upload/growth-file", handlers.UploadGrowthFile)

		growth := api.Group("/growth")
		{
			growth.GET("/records", handlers.ListGrowthMetricRecords)
			growth.POST("/records", handlers.CreateGrowthMetricRecord)
			growth.DELETE("/records/:id", handlers.DeleteGrowthMetricRecord)
		}

		// 课表模块
		timetable := api.Group("/timetable")
		{
			timetable.GET("/config", handlers.GetTimetableConfig)
			timetable.PUT("/config", handlers.UpdateTimetableConfig)
			timetable.GET("/courses", handlers.GetCourses)
			timetable.POST("/courses", handlers.AddCourse)
			timetable.PUT("/courses/:id", handlers.UpdateCourse)
			timetable.DELETE("/courses/:id", handlers.DeleteCourse)
			timetable.GET("/data", handlers.GetTimetable)
			timetable.POST("/data", handlers.SaveTimetable)
		}

	}

	r.PUT("/api/storage/put", handlers.StoragePut)
	return r
}
