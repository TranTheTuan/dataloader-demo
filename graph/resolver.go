package graph

import (
	"github.com/TranTheTuan/dataloader-demo/internal/core/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
// This mirrors the original project's graph/resolver.go pattern.

type Resolver struct {
	RedemptionSvc service.RedemptionSvc
}
