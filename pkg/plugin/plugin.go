package plugin

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/wyretech/autotask-datasource/pkg/datasource"
)

// NewDatasourceInstance creates a new instance of the Autotask datasource
func NewDatasourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return datasource.NewAutotaskDataSource(settings)
}

// PluginID is the unique identifier for the Autotask datasource plugin
const PluginID = "wyretech-autotask-datasource"

// PluginName is the display name for the Autotask datasource plugin
const PluginName = "Autotask"

// PluginType is the type identifier for the Autotask datasource plugin
const PluginType = "autotask"
