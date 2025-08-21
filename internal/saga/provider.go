package saga

import "github.com/google/wire"

// ProviderSet is saga providers.
var ProviderSet = wire.NewSet(
	NewCoordinator,
	// 这里可以添加其他 Saga 相关的 Provider
)
