package ossrs.net.srssip.gb28181.listeners;

import com.alibaba.fastjson.JSON;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.gb28181.event.messageevent.MessageEventAbstract;
import ossrs.net.srssip.gb28181.event.request.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.listeners.factory.MessageEventFactory;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.message.Request;

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
        log.info(request.toString());
        MessageEventAbstract messageEventAbstract = MessageEventFactory.INSTANCE.
                getMessageEvent(requestEvent);
        applicationEventPublisher.publishEvent(messageEventAbstract);
    }

    @Override
    public void processResponse(ResponseEvent responseEvent) {
        log.info(responseEvent.toString());
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
