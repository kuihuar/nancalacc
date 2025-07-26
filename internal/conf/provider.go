package conf

func ProvideDingtalkConfig(confServcie *Service) *Service_Auth_Dingtalk {
	return confServcie.GetAuth().GetDingtalk() // 使用Protobuf生成的Get方法
}

// 提供 Authenticator 配置的提取函数
func ProvideAuthConfig(confServcie *Service) *Service_Auth {
	return confServcie.GetAuth()
}

func ProvideBusinessConfig(confServcie *Service) *Service_Business {
	return confServcie.GetBusiness()
}
