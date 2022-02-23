package ossrs.net.srssip.gb28181.event.messageevent;

import gov.nist.javax.sip.header.CSeq;
import lombok.Getter;
import ossrs.net.srssip.gb28181.annotation.MessageEventHandler;

import javax.sip.Dialog;
import javax.sip.message.Request;

/**
 * @ Description ossrs.net.srssip.gb28181.event.messageevent
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:27
 */
@MessageEventHandler("ACK")
public class AckEvent extends MessageEventAbstract{

    @Getter
    private long seqNumber;

    @Getter
    private Dialog dialog;

    @Override
    public void process() {
        Request request = requestEvent.getRequest();
        this.dialog = requestEvent.getDialog();
        CSeq csReq = (CSeq) request.getHeader(CSeq.NAME);
        seqNumber = csReq.getSeqNumber();
    }
}
