package ossrs.net.srssip.gb28181.event.request;

import gov.nist.javax.sip.header.CSeq;
import lombok.Getter;
import ossrs.net.srssip.gb28181.annotation.RequestEventHandler;

import javax.sip.Dialog;
import javax.sip.header.CSeqHeader;
import javax.sip.message.Request;

/**
 * @ Description ossrs.net.srssip.gb28181.event.request.messageevent
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:27
 */
@RequestEventHandler(value = "ACK")
public class AckEvent extends RequestEventAbstract {

    @Getter
    private long seqNumber;

    @Getter
    private Dialog dialog;

    @Override
    public void process() {
        Request request = requestEvent.getRequest();
        this.dialog = requestEvent.getDialog();
        CSeqHeader cSeqHeader = (CSeqHeader) request.getHeader(CSeq.NAME);
        CSeq csReq = (CSeq) request.getHeader(CSeq.NAME);
        seqNumber = csReq.getSeqNumber();
    }
}
