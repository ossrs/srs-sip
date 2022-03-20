package ossrs.net.srssip.gb28181.event.response;

import gov.nist.javax.sip.ResponseEventExt;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import javax.sip.*;
import javax.sip.address.SipURI;
import javax.sip.header.CSeqHeader;
import javax.sip.header.ViaHeader;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.io.UnsupportedEncodingException;
import java.text.ParseException;

/**
 * @ Description ossrs.net.srssip.gb28181.event.response
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午12:56
 */
@Slf4j
public abstract class ResponseEventAbstract {
    @Getter
    private ResponseEvent responseEvent;

    @Getter
    private Dialog dialog;

    @Getter
    private Request reqAck;

    @Getter
    private String content;

    public Response response;

    public abstract void process();

    public void constructor(ResponseEvent responseEvent) {
        this.responseEvent = responseEvent;
        this.dialog = responseEvent.getDialog();
        this.response = responseEvent.getResponse();
        CSeqHeader cseq = (CSeqHeader) response.getHeader(CSeqHeader.NAME);
        ResponseEventExt event = (ResponseEventExt)responseEvent;

        try {
            this.content = new String(response.getRawContent(), "GBK");
            this.reqAck = dialog.createAck(cseq.getSeqNumber());
            SipURI requestURI = (SipURI) reqAck.getRequestURI();
            requestURI.setHost(event.getRemoteIpAddress());
            requestURI.setPort(event.getRemotePort());
            reqAck.setRequestURI(requestURI);
        } catch (UnsupportedEncodingException | ParseException | InvalidArgumentException | SipException e) {
            e.printStackTrace();
        }
    }

    protected void responseAck(){
        try {
            this.dialog.sendAck(reqAck);
        } catch (SipException e) {
            log.error("响应消息发送失败",e);
        }
    }

    protected boolean isError() {
        return response.getStatusCode() / 100 >= 4;
    }

    protected boolean isProvisional() {
        return response.getStatusCode() / 100 == 1;
    }

    protected boolean isFinal() {
        return response.getStatusCode() >= 200;
    }

    protected boolean isSuccess(){
        return response.getStatusCode()/ 100 == 2;
    }

    protected boolean isRedirect() {
        return response.getStatusCode() / 100 == 3;
    }

    protected boolean isClientError() {
        return response.getStatusCode() / 100 == 4;
    }

    protected boolean isServerError() {
        return response.getStatusCode() / 100 == 5;
    }

    protected boolean isGlobalError() {
        return response.getStatusCode() / 100 == 6;
    }

    protected boolean is100Trying() {
        return response.getStatusCode() == 100;
    }

    protected boolean isRinging() {
        return response.getStatusCode() == 180 || response.getStatusCode() == 183;
    }
}
