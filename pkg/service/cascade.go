package service

import "github.com/emiago/sipgo"

type Cascade struct {
	ua     *sipgo.UserAgent
	sipCli *sipgo.Client
	sipSvr *sipgo.Server
}
