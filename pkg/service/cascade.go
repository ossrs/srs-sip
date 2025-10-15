package service

import (
	"context"

	"github.com/emiago/sipgo"
	"github.com/ossrs/srs-sip/pkg/config"
)

type Cascade struct {
	ua     *sipgo.UserAgent
	sipCli *sipgo.Client
	sipSvr *sipgo.Server

	ctx  context.Context
	conf *config.MainConfig
}
