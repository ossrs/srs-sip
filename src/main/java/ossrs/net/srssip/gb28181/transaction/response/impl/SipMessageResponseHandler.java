package ossrs.net.srssip.gb28181.transaction.response.impl;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import ossrs.net.srssip.gb28181.transaction.factory.ServerTransactionFactory;
import ossrs.net.srssip.gb28181.transaction.response.MessageResponseHandler;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.util.Optional;

/**
 * @ Description ossrs.net.srssip.gb28181.transaction.response.impl
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午2:24
 */
@Slf4j
@Service
public class SipMessageResponseHandler implements MessageResponseHandler {

    @Resource(name = "sipTcpProvider")
    private SipProvider sipTcpProvider;
    @Resource(name = "sipUdpProvider")
    private SipProvider sipUdpProvider;

    @Override
    public void sendResponse(RequestEvent requestEvent, Response response) {
        ServerTransaction serverTransaction = Optional.ofNullable(requestEvent.getServerTransaction())
                .orElseGet(() ->
                        ServerTransactionFactory.INSTANCE.getServerTransaction(sipTcpProvider, sipUdpProvider,
                                requestEvent.getRequest()));
        try {
            serverTransaction.sendResponse(response);
        } catch (SipException | InvalidArgumentException e) {
            log.error("发送消息失败", e);
        }
    }

    @Override
    public ClientTransaction sendResponse(String transport, Request request) {
        return null;
    }

    @Override
    public void sendDialog(Dialog dialog, Request request, String transport) {

    }
}
