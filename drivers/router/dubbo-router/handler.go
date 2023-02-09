package dubbo_router

import (
	"github.com/eolinker/apinto/router"
	"github.com/eolinker/apinto/service"
	"github.com/eolinker/eosc/eocontext"
	dubbo_context "github.com/eolinker/eosc/eocontext/dubbo-context"
)

var _ router.IRouterHandler = (*dubboHandler)(nil)

type dubboHandler struct {
	completeHandler eocontext.CompleteHandler
	routerName      string
	routerId        string
	serviceName     string
	disable         bool
	service         service.IService
}

func (d *dubboHandler) ServeHTTP(ctx eocontext.EoContext) {

	_, err := dubbo_context.Assert(ctx)
	if err != nil {
		return
	}
	if d.disable {
		//httpContext.Response().SetStatus(http.StatusNotFound, "")
		//httpContext.Response().SetBody([]byte("router disable"))
		//httpContext.FastFinish()
		return
	}

	//Set Label
	ctx.SetLabel("api", d.routerName)
	ctx.SetLabel("api_id", d.routerId)
	ctx.SetLabel("service", d.serviceName)
	ctx.SetLabel("service_id", d.service.Id())
	//ctx.SetLabel("ip", dubboCtx.Request().ReadIP())

	ctx.SetCompleteHandler(d.completeHandler)
	ctx.SetApp(d.service)
	ctx.SetBalance(d.service)
	ctx.SetUpstreamHostHandler(d.service)

}
