package ossrs.net.srssip.gb28181.transaction.factory;

import gov.nist.javax.sip.SipStackImpl;
import gov.nist.javax.sip.message.SIPRequest;
import gov.nist.javax.sip.stack.SIPServerTransaction;
import lombok.extern.slf4j.Slf4j;

import javax.sip.*;
import javax.sip.header.CallIdHeader;
import javax.sip.header.ViaHeader;
import javax.sip.message.Request;
import java.util.Optional;

/**
 * @ Description ossrs.net.srssip.gb28181.transaction.factory
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午2:35
 */
@Slf4j
public class ServerTransactionFactory {
    public static ServerTransactionFactory INSTANCE = new ServerTransactionFactory();

    public ServerTransaction getServerTransaction(SipProvider sipTcpProvider, SipProvider sipUdpProvider, Request request) {
        ViaHeader reqViaHeader = (ViaHeader) request.getHeader(ViaHeader.NAME);
        String transport = reqViaHeader.getTransport();
        if (ListeningPoint.TCP.equals(transport)) {
            return getServerTransaction(sipTcpProvider, request);
        } else if (ListeningPoint.UDP.equals(transport)) {
            return getServerTransaction(sipUdpProvider, request);
        } else {
            return null;
        }
    }

    private ServerTransaction getServerTransaction (SipProvider sipProvider, Request request) {
        SipStackImpl sipStack =  (SipStackImpl)sipProvider.getSipStack();
        return Optional.ofNullable((SIPServerTransaction) sipStack.findTransaction((SIPRequest) request, true)).orElseGet(() -> {
            try {
                return (SIPServerTransaction) sipProvider.getNewServerTransaction(request);
            } catch (TransactionAlreadyExistsException | TransactionUnavailableException e) {
                log.error("获取SipServerTransaction出错", e);
                return null;
            }
        });
    }

    public ClientTransaction getClientTransaction(SipProvider sipTcpProvider, SipProvider sipUdpProvider,
                                                  String transport, Request request) throws TransactionUnavailableException {
        if (ListeningPoint.TCP.equals(transport)) {
            return sipTcpProvider.getNewClientTransaction(request);
        }
        return sipUdpProvider.getNewClientTransaction(request);
    }

    public CallIdHeader getCallIdHeader(SipProvider sipTcpProvider, SipProvider sipUdpProvider, String transport) {
        if (ListeningPoint.TCP.equals(transport)) {
            return sipTcpProvider.getNewCallId();
        }
        return sipUdpProvider.getNewCallId();
    }
}
