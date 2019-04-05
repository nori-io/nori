package service

//
//import (
//	"github.com/cheebo/gorest"
//	"github.com/nori-io/nori/core/plugins"
//)
//
//type RESTService interface {
//	Plugins()
//}
//
//type PluginsService interface {
//	// list of plugins ready to start
//	Runnable() (rest.ListResp, error)
//	/// list of plugins ready to install
//	Installable() (rest.ListResp, error)
//	// already running plugins
//	Running() (rest.ListResp, error)
//	// start plugin
//}
//
//type service struct {
//	p pluginService
//}
//
//type pluginService struct {
//	pluginManager plugins.Manager
//}
//
//func NewRestService(pm plugins.Manager) RESTService {
//	return service{
//		p: pluginService{
//			pluginManager: pm,
//		},
//	}
//}
