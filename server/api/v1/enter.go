package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/brandtrekin"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup      system.ApiGroup
	ExampleApiGroup     example.ApiGroup
	BrandtrekinApiGroup brandtrekin.ApiGroup
}
