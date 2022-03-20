package ossrs.net.srssip.gb28181.listeners;

import com.alibaba.fastjson.JSON;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.gb28181.event.request.RequestEventAbstract;
import ossrs.net.srssip.gb28181.event.response.ResponseEventAbstract;
import ossrs.net.srssip.gb28181.listeners.factory.MessageEventFactory;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.header.CSeqHeader;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.io.UnsupportedEncodingException;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午11:14
 */
@Slf4j
@Component
public class SipServerListeners implements SipListener {

    @Resource
    private ApplicationEventPublisher applicationEventPublisher;


    @Override
    public void processRequest(RequestEvent requestEvent) {
        Request request = requestEvent.getRequest();
        log.info("RequestEvent {}",request.toString());
        RequestEventAbstract requestEventAbstract = MessageEventFactory.INSTANCE.
                getRequestEvent(requestEvent);
        applicationEventPublisher.publishEvent(requestEventAbstract);
    }

    @Override
    public void processResponse(ResponseEvent responseEvent) {
        Response response = responseEvent.getResponse();
        try {
            log.info("ResponseEvent {}",response.toString() );
            log.info("ResponseContent {}",new String(response.getRawContent(),"gbk") );
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }
        int status = response.getStatusCode();
        if (status >= HttpStatus.OK.value() && status <= HttpStatus.MULTIPLE_CHOICES.value()) {
            CSeqHeader cseqHeader = (CSeqHeader) response.getHeader(CSeqHeader.NAME);
            String method = cseqHeader.getMethod();
            ResponseEventAbstract messageResponse = MessageEventFactory.INSTANCE.
                    getResponseEvent(method, responseEvent);
            applicationEventPublisher.publishEvent(messageResponse);
        }
    }

    @Override
    public void processTimeout(TimeoutEvent timeoutEvent) {
        log.error(timeoutEvent.toString());
    }

    @Override
    public void processIOException(IOExceptionEvent exceptionEvent) {
        log.error(exceptionEvent.toString());
    }

    @Override
    public void processTransactionTerminated(TransactionTerminatedEvent transactionTerminatedEvent) {
        log.error(JSON.toJSONString(transactionTerminatedEvent));
    }

    @Override
    public void processDialogTerminated(DialogTerminatedEvent dialogTerminatedEvent) {
        log.error(dialogTerminatedEvent.toString());
    }
}
