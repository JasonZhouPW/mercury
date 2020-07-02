package rest

import (
	"git.ont.io/ontid/otf/config"
	"git.ont.io/ontid/otf/message"
	"git.ont.io/ontid/otf/middleware"
	"git.ont.io/ontid/otf/service"
	"git.ont.io/ontid/otf/service/controller"
	"git.ont.io/ontid/otf/store"
	"git.ont.io/ontid/otf/utils"
	"github.com/gin-gonic/gin"
	sdk "github.com/ontio/ontology-go-sdk"
)

var (
	Svr *service.Service
)

func NewService(acct *sdk.Account, cfg *config.Cfg, db store.Store, msgSvr *service.MsgService) {
	Svr = service.NewService()
	Svr.RegisterController(controller.NewSyscontroller(acct, cfg, db, msgSvr))
	Svr.RegisterController(controller.NewCredentialController(acct, cfg, db, msgSvr))
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.LoggerToFile())
	r.Use(gin.Recovery())
	v := r.Group(utils.Group_Api_V1)
	{
		v.POST(utils.Invite_Api, Invite)
		//v.POST(utils.SendConnectionReq_Api, SendConnectionReq)
		v.POST(utils.ConnectRequest_Api, ConnectRequest)
		v.POST(utils.ConnectResponse_Api, ConnectResponse)
		v.POST(utils.ConnectAck_Api, ConnectAck)

		v.POST(utils.SendProposalCredentialReq_Api, SendProposalCredentialReq)
		v.POST(utils.OfferCredential_Api, OfferCredential)
		v.POST(utils.ProposalCredentialReq_Api, ProposalCredentialReq)
		v.POST(utils.SendRequestCredential_Api, SendRequestCredential)
		v.POST(utils.RequestCredential_Api, RequestCredential)
		v.POST(utils.IssueCredential_Api, IssueCredential)
		v.POST(utils.CredentialAckInfo_Api, CredentialAckInfo)

		v.POST(utils.SendRequestPresentation_Api, SendRequestPresentation)
		v.POST(utils.RequestPresentation_Api, RequestPresentation)
		v.POST(utils.Presentation_Api, Presentation)
		v.POST(utils.PresentationAckInfo, PresentationAckInfo)

		v.POST(utils.SendGeneralMsg, SendGeneralMsg)
		v.POST(utils.ReceiveGeneralMsg, ReceiveGeneralMsg)
	}
	return r
}

func SendMsg(msgType message.MessageType, data interface{}) (interface{}, error) {
	msg := message.Message{MessageType: msgType, Content: data}
	resp, err := Svr.Serv(msg)
	if err != nil {
		middleware.Log.Errorf("err:%s", err)
		return nil, err
	}
	return resp.GetMessage()
}
