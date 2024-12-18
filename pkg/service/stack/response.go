package stack

import (
	"fmt"
	"time"

	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/srs-sip/pkg/utils"
)

const TIME_LAYOUT = "2024-01-01T00:00:00"
const EXPIRES_TIME = 3600

func newResponse(req *sip.Request, code sip.StatusCode, reason string) *sip.Response {
	resp := sip.NewResponseFromRequest(req, code, reason, nil)

	newTo := &sip.ToHeader{Address: resp.To().Address, Params: sip.NewParams()}
	newTo.Params.Add("tag", sip.GenerateTagN(10))

	resp.ReplaceHeader(newTo)
	resp.RemoveHeader("Allow")

	return resp
}

func NewRegisterResponse(req *sip.Request, code sip.StatusCode, reason string) *sip.Response {
	resp := newResponse(req, code, reason)

	expires := sip.ExpiresHeader(EXPIRES_TIME)
	resp.AppendHeader(&expires)
	resp.AppendHeader(sip.NewHeader("Date", time.Now().Format(TIME_LAYOUT)))

	return resp
}

func NewUnauthorizedResponse(req *sip.Request, code sip.StatusCode, reason string, realm string) *sip.Response {
	resp := newResponse(req, code, reason)

	nonce := utils.GenerateNonce()
	resp.AppendHeader(sip.NewHeader("WWW-Authenticate", fmt.Sprintf(`Digest realm="%s",nonce="%s",algorithm=MD5`, realm, nonce)))

	return resp
}
