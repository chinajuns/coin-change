package route

import (
	"io"
	"log"
	"okc/app/controller/admin"
	"okc/app/controller/api"
	"okc/app/middleware"
	"okc/app/service"
	"okc/utils"
	"os"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

// 首页
func homeHandle(c *gin.Context) {
	c.String(200, "welcome")
}

func init() {
	// 开启日志写入
	f, err := utils.GetRootLogFile()
	if err != nil {
		log.Println(err)
		return
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 初始化gin
	r = gin.Default()

	// 加载日志中间件
	r.Use(middleware.LogMiddleware())
	r.Static("/storage/images", "./storage/images")
	// 首页
	r.GET("/", homeHandle)

	// api路由组
	apis := r.Group("/api")
	{
		// 发送验证码
		apis.POST("/sms_mail", api.SendEmailCode)

		// 用户路由组
		user := apis.Group("/user")
		{
			// 用户邮箱注册
			user.POST("/register", api.UserEmailRegister)
			// 用户登录
			user.POST("/login", api.UserLogin)
		}

		// 获取币种和行情
		apis.GET("/currency/quotation_new", api.NewQuotation)
		// 获取交易对详情
		apis.GET("/currency/deal", api.CurrencyMatchesFind)
		// 单图片上传
		apis.POST("/upload/one", api.UploadOne)
		// 获取轮播图
		apis.GET("/news/banner", api.QueryNewsBanner)

		// auth认证路由组
		apiAuth := apis.Group("", middleware.ApiMiddleware())
		{
			// 获取用户信息(基础认证)
			apiAuth.GET("/user/center", api.UserCenter)
			// 获取用户信息(高级认证)
			apiAuth.GET("/user/center2", api.UserCenter2)
			// 修改用户密码
			apiAuth.POST("/user/e_pwd", api.UpdatePwd)
			// 修改交易密码
			apiAuth.POST("/safe/update_password", api.UpdatePayPwd)
			// 获取用户提币记录
			apiAuth.GET("/charge_mention/log3", api.UserTiBiRecord)
			// 获取用户冲币记录
			apiAuth.GET("/charge_mention/log2", api.UserChongBiRecord)
			// 保存用户钱包地址
			apiAuth.POST("/wallet/addaddress", api.UserWalletSaveAddress)
			// 获取用户钱包地址
			apiAuth.POST("/wallet/get_address", api.UserWalletQueryAddress)
			// 获取用户钱包详情
			apiAuth.POST("/wallet/detail", api.UserWalletQueryDetail)
			// 获取用户充值地址
			apiAuth.POST("/wallet/get_in_address", api.UserRechargeAddress)
			// 用户充值
			apiAuth.POST("/wallet/charge_req", api.UserChargeReq)
			// 提币限制
			apiAuth.POST("/wallet/get_info", api.CurrencyRestrict)
			// 用户提币
			apiAuth.POST("/wallet/out", api.UserWalletOut)
			// 检查用户是否有交易密码
			apiAuth.POST("/user/has_password", api.UserCheckHasPayPassword)
			// 秒合约（期权）交易列表
			apiAuth.GET("/microtrade/lists", api.MicroOrderList)
			// 秒合约（期权）时间列表
			apiAuth.GET("/microtrade/seconds", api.MicroSecondList)
			// 秒合约（期权）可支付币种列表
			apiAuth.GET("/microtrade/payable_currencies", api.PayableCurrencies)
			// 秒合约（期权）提交订单
			apiAuth.GET("/microtrade/submit", api.MicroOrderSubmit)
			// 秒合约（合约）交易信息
			apiAuth.POST("/lever/deal", api.LeverDeal)
			// 秒合约（合约）我的交易
			apiAuth.POST("/lever/my_trade", api.LeverMyTrade)
			// 秒合约（合约）平仓
			apiAuth.POST("/lever/close", api.LeverClose)
			// 秒合约（合约）提交订单
			apiAuth.POST("/lever/submit", api.LeverSubmit)
			// 币币交易列表
			apiAuth.GET("/coin/list", api.CoinTradeList)
			// 币币交易订单提交
			apiAuth.POST("/coin/trade", api.CoinTradeSubmit)
			// 币币交易订单取消
			apiAuth.PUT("/coin/trade", api.CoinTradeCancel)
			// 用户钱包转划
			apiAuth.POST("/wallet/change", api.UserWalletChange)
			// 获取kline历史记录
			apiAuth.POST("/kline/history", api.KlineHistoryV2)
		}
	}

	// 管理员登录
	r.POST("/admin/login", admin.Login)
	// admin 路由组
	admins := r.Group("/admin", middleware.AdminMiddleware())
	{
		// 添加管理员
		admins.POST("/manager/add", new(admin.Admin).AddAdmin)
		// 编辑管理员
		admins.POST("/manager/edit", new(admin.Admin).EditAdmin)
		// 删除管理员
		admins.POST("/manager/delete", new(admin.Admin).DeleteAdmin)
		// 获取管理员列表
		admins.GET("/manager/list", new(admin.Admin).QueryAdminList)

		// 添加管理员角色
		admins.POST("/role/add", new(admin.AdminRole).AddAdminRole)
		// 编辑管理员角色
		admins.POST("/role/edit", new(admin.AdminRole).EditAdminRole)
		// 删除管理员角色
		admins.POST("/role/delete", new(admin.AdminRole).DeleteAdminRole)
		// 获取管理员角色列表
		admins.GET("/role/list", new(admin.AdminRole).QueryAdminRoleList)
		// 获取管理员所有角色
		admins.GET("/role/all", new(admin.AdminRole).QueryAdminRoleAll)

		// 获取币种列表
		admins.GET("/currency/lists", new(admin.AdminCurrency).CurrencyLists)
		// 添加币种
		admins.POST("/currency/add", new(admin.AdminCurrency).AddCurrency)
		// 删除币种
		admins.POST("/currency/delete", new(admin.AdminCurrency).DeleteCurrency)
		// 编辑币种
		admins.POST("/currency/edit", new(admin.AdminCurrency).EditCurrency)

		// 获取充值申请列表
		admins.GET("/charge/list", new(admin.ChargeReq).QueryChargeReqList)
		// 充值申请同意
		admins.POST("/charge/pass", new(admin.ChargeReq).ChargePass)
		// 充值申请拒绝
		admins.POST("/charge/refuse", new(admin.ChargeReq).ChargeReqRefuse)

		// 获取提币申请列表
		admins.GET("/walletOut/list", new(admin.UserWalletOut).QueryUserWalletOutList)
		// 获取提币申请信息
		admins.GET("/walletOut/info", new(admin.UserWalletOut).QueryUserWalletOutInfo)
		// 同意提币申请
		admins.POST("/walletOut/pass", new(admin.UserWalletOut).QueryUserWalletOutPass)

		// 获取用户日志列表
		admins.GET("/accountLog/list", new(admin.AccountLog).QueryAccountLogList)

		// 获取用户列表
		admins.GET("/users/list", new(admin.Users).QueryUserList)
		// 获取用户钱包列表
		admins.GET("/users/walletList", new(admin.Users).QueryUserWalletList)
	}

	// websocket
	//r.GET("/socket.io", api.SocketIO)
	//r.GET("/socket.io/kline", api.SocketIOKline)
	r.GET("/socket.io/*any", gin.WrapH(service.Server))
}

// Run 启动web服务
func Run(addr ...string) error {

	go service.Server.Serve()
	defer service.Server.Close()
	defer close(service.SocketIoManage.Register)
	defer close(service.SocketIoManage.UnRegister)
	go service.SocketIoManage.RegisterGo()
	go service.SocketIoManage.UnRegisterGo()

	go new(service.Candle1D).SubOpenKlineServer("wss://ws.okx.com:8443/ws/v5/public")
	go new(service.Tickers).SubOpenTickersServer("wss://ws.okx.com:8443/ws/v5/public")

	err := r.Run(addr...)
	if err != nil {
		return err
	}

	return nil
}
