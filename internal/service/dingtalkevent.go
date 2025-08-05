package service

import (
	"context"
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	clientV2 "github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
)

type DingTalkEventService struct {
	confService      *conf.Service
	log              *log.Helper
	accounterUsecase *biz.AccounterIncreUsecase
	running          atomic.Bool
	cancel           context.CancelFunc
	//client clientV2.OpenDingTalkClient
}

func NewDingTalkEventService(confService *conf.Service, logger log.Logger, accounterUsecase *biz.AccounterIncreUsecase) *DingTalkEventService {
	return &DingTalkEventService{confService: confService, log: log.NewHelper(logger), accounterUsecase: accounterUsecase}
}

func (es *DingTalkEventService) Start() {
	log.Info(es.confService.Auth.Dingtalk)

	cred := &clientV2.AuthClientCredential{
		ClientId:     es.confService.Auth.Dingtalk.AppKey,
		ClientSecret: es.confService.Auth.Dingtalk.AppSecret,
	}

	ctx, cancel := context.WithCancel(context.Background())
	es.cancel = cancel
	es.running.Store(true)
	go func() {
		defer es.running.Store(false)
		e := clientV2.
			NewBuilder().
			Credential(cred).
			//监听开放平台事件
			RegisterAllEventHandler(es.HandleEvent).
			Build().
			Start(ctx)
		if e != nil {
			log.Error("DingTalkEventService.Start failed", e.Error())
		}
		log.Info("DingTalkEventService Start")
	}()
	log.Info("DingTalkEventService.Starting...")

}
func (es *DingTalkEventService) Stop() {
	if !es.running.Load() {
		return
	}
	log.Info("=====DingTalkEventService.Stop===")
	es.cancel()
}
func (es *DingTalkEventService) Running() bool {
	return es.running.Load()

}

func (es *DingTalkEventService) HandleEvent(event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {
	println("HandleEvent ", event.Data)

	ctx := context.Background()
	switch event.EventType {
	case "org_dept_create":
		es.log.Infof("org_dept_create: %v", event.Data)
		es.OrgDeptCreate(ctx, event)
	case "org_dept_modify":
		es.log.Infof("org_dept_modify: %v", event.Data)
		es.OrgDeptModify(ctx, event)
	case "org_dept_remove":
		es.log.Infof("org_dept_remove: %v", event.Data)
		es.OrgDeptRemove(ctx, event)
	case "user_add_org":
		es.log.Infof("user_add_org: %v", event.Data)
		es.UserAddOrg(ctx, event)
	case "user_modify_org":
		es.log.Infof("user_modify_org: %v", event.Data)
	case "user_leave_org":
		es.UserLeaveOrg(ctx, event)
		es.log.Infof("user_leave_org: %v", event.Data)
	default:
		es.log.Infof("unknown event: %v", event.Data)
		return clientV2.EventStatusSuccess
	}
	return clientV2.EventStatusSuccess
}

func (es *DingTalkEventService) OrgDeptCreate(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {
	es.log.Infof("OrgDeptCreate: %v", event.Data)
	err := es.accounterUsecase.OrgDeptCreate(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}

	return clientV2.EventStatusSuccess
}
func (es *DingTalkEventService) OrgDeptModify(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {

	err := es.accounterUsecase.OrgDeptModify(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}
	return clientV2.EventStatusSuccess
}
func (es *DingTalkEventService) OrgDeptRemove(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {

	err := es.accounterUsecase.OrgDeptRemove(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}
	return clientV2.EventStatusSuccess
}
func (es *DingTalkEventService) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {

	err := es.accounterUsecase.UserAddOrg(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}
	return clientV2.EventStatusSuccess
}
func (es *DingTalkEventService) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {

	err := es.accounterUsecase.UserModifyOrg(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}
	return clientV2.EventStatusSuccess
}
func (es *DingTalkEventService) UserLeaveOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {

	err := es.accounterUsecase.UserLeaveOrg(ctx, event)
	if err != nil {
		return clientV2.EventStatusLater
	}
	return clientV2.EventStatusSuccess
}
