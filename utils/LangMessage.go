package utils

const (
	ParameterError          = iota // 参数错误
	PasswordHshError               // 两次密码不一致
	CodeError                      // 验证码错误
	AccountExistsError             // 账号已存在
	PasswordLenError               // 密码只能在6-16位之间
	PcodeError                     // 请填写正确的邀请码
	ServerError                    // 服务器异常
	SendCodeError                  // 验证码发送失败
	Success                        // 成功
	UserFindError                  // 用户未找到
	PasswordInputError             // 请输入密码
	UserNameInputError             // 请输入账号
	PasswordError                  // 密码错误
	UserStatusError                // 用户状态异常
	UserAuthError                  // 用户认证失败
	UpdateError                    // 更新失败
	UpdateSuccess                  // 更新成功
	AddressInputError              // 请输入地址
	Error                          // 失败
	CurrencyFindError              // 币种未找到
	UserWalletFindError            // 用户钱包未找到
	NumberLessTHanZeroError        // 数值不能小于零
	PayPasswordNilError            // 交易密码未设置
	NumberLessThanMinError         // 数值不能小于最小值
	BalanceNotEnoughError          // 余额不足
)

// GetLangMessage
// lang       语言[zh|en]
// constError 本包中的常量
// 返回对应语言错误消息
func GetLangMessage(lang string, constError int) (msg string) {
	messageMap := map[string][]string{
		"zh": {
			"参数错误",
			"两次密码不一致",
			"验证码错误",
			"账号已存在",
			"密码只能在6-16位之间",
			"请填写正确的邀请码",
			"服务器异常",
			"验证码发送失败",
			"成功",
			"用户未找到",
			"请输入密码",
			"请输入账号",
			"密码错误",
			"用户状态异常",
			"用户认证失败",
			"更新失败",
			"更新成功",
			"请输入地址",
			"失败",
			"币种未找到",
			"用户钱包未找到",
			"数值不能小于零",
			"交易密码未设置",
			"数值不能小于最小值",
			"余额不足",
		},
		"en": {
			"Parameter error",
			"The two passwords are inconsistent",
			"Verification code error",
			"Account already exists",
			"Password can only be between 6-16 digits",
			"Please fill in the correct invitation code",
			"Server Error",
			"Send Code Error",
			"Success",
			"User not found",
			"Please enter the password",
			"Please input Username",
			"wrong password",
			"User Status Abnormal",
			"User Auth Error",
			"Update Error",
			"Update Success",
			"Please input Address",
			"Error",
			"Currency not found",
			"UserWallet not found",
			"Please Number not less than zero",
			"PayPassword not set",
			"Please Number not less than min Number",
			"Sorry Balance not enough you",
		},
	}

	if _, ok := messageMap[lang]; !ok {
		msg = messageMap["en"][constError]
	} else {
		msg = messageMap[lang][constError]
	}

	return
}
