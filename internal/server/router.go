package server

import (
	"demo_message/internal/service/api"

	http "github.com/go-kratos/kratos/v2/transport/http"

	campaign_v1 "demo_message/proto/api/campaign/v1"
	message_v1 "demo_message/proto/api/message/v1"
	system_v1 "demo_message/proto/api/system/v1"
)

type Router struct {
}

// InitHTTPRouter register HTTP router.
func InitHTTPRouter(srv *http.Server,
	campaignApi *api.CampaignApi,
	messageApi *api.MessageApi,
	systemApi *api.SystemApi,
) (r *Router) {

	// router
	campaign_v1.RegisterCampaignHTTPServer(srv, campaignApi)
	message_v1.RegisterMessageHTTPServer(srv, messageApi)
	system_v1.RegisterSystemHTTPServer(srv, systemApi)

	return
}
