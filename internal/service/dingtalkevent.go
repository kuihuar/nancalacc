package service

import (
	"context"
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	clientV2 "github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
)

type DingTalkEventService struct {
	conf             *conf.Data
	log              *log.Helper
	accounterUsecase *biz.AccounterUsecase
	//client clientV2.OpenDingTalkClient
}

func NewDingTalkEventService(conf *conf.Data, logger log.Logger, accounterUsecase *biz.AccounterUsecase) *DingTalkEventService {
	return &DingTalkEventService{conf: conf, log: log.NewHelper(logger), accounterUsecase: accounterUsecase}
}

func (es *DingTalkEventService) Start() {
	log.Info(es.conf.Dingtalk)
	log.Info(es.conf.ServiceConf)

	cred := &clientV2.AuthClientCredential{
		ClientId:     es.conf.Dingtalk.AppKey,
		ClientSecret: es.conf.Dingtalk.AppSecret,
	}

	e := clientV2.
		NewBuilder().
		Credential(cred).
		//监听开放平台事件
		RegisterAllEventHandler(es.HandleEvent).
		Build().
		Start(context.Background())

	if e != nil {
		log.Error("failed to start stream client", e.Error())
		return
	}
	log.Info("=====DingTalkEventService.Start, please commit sub event===")
	select {}
}
func (es *DingTalkEventService) Stop() {}

func (es *DingTalkEventService) HandleEvent(event *clientV2.GenericOpenDingTalkEvent) clientV2.EventStatus {
	println("HandleEvent ", event.Data)

	switch event.EventType {
	case "org_dept_create":
		es.log.Infof("org_dept_create: %v", event.Data)
		es.OrgDeptCreate(context.Background(), event)
	case "org_dept_modify":
		es.log.Infof("org_dept_modify: %v", event.Data)
		es.OrgDeptModify(context.Background(), event)
	case "org_dept_remove":
		es.log.Infof("org_dept_remove: %v", event.Data)
		es.OrgDeptRemove(context.Background(), event)
	case "user_add_org":
		es.log.Infof("user_add_org: %v", event.Data)
		es.UserAddOrg(context.Background(), event)
	case "user_modify_org":
		es.log.Infof("user_modify_org: %v", event.Data)
	case "user_leave_org":
		es.UserLeaveOrg(context.Background(), event)
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
