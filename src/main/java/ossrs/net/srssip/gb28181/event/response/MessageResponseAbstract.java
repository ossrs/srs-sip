package ossrs.net.srssip.gb28181.event.response;

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
public abstract class MessageResponseAbstract {
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

        try {
            this.content = new String(response.getRawContent(), "GBK");
            this.reqAck = dialog.createAck(cseq.getSeqNumber());
            SipURI requestURI = (SipURI) reqAck.getRequestURI();
            ViaHeader viaHeader = (ViaHeader) response.getHeader(ViaHeader.NAME);
            requestURI.setHost(viaHeader.getHost());
            requestURI.setPort(viaHeader.getPort());
            reqAck.setRequestURI(requestURI);
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        } catch (ParseException e) {
            e.printStackTrace();
        } catch (InvalidArgumentException e) {
            e.printStackTrace();
        } catch (SipException e) {
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
}
