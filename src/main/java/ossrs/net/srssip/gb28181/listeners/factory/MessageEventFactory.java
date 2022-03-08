package ossrs.net.srssip.gb28181.listeners.factory;

import lombok.extern.slf4j.Slf4j;
import org.reflections.Reflections;
import ossrs.net.srssip.gb28181.annotation.MessageRequestEventHandler;
import ossrs.net.srssip.gb28181.annotation.MessageResponseEventHandler;
import ossrs.net.srssip.gb28181.annotation.MessageEventHandler;
import ossrs.net.srssip.gb28181.event.request.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.event.messageevent.MessageEventAbstract;
import ossrs.net.srssip.gb28181.event.response.MessageResponseAbstract;

import javax.sip.RequestEvent;
import javax.sip.ResponseEvent;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午12:35
 */
@Slf4j
public class MessageEventFactory {

    public static MessageEventFactory INSTANCE = new MessageEventFactory();

    private static final Map<String, MessageEventAbstract> messageEventMap = new ConcurrentHashMap<>();

    private static final Map<String, MessageRequestAbstract> messageRequestMap = new ConcurrentHashMap<>();

    private static final Map<String, MessageResponseAbstract> messageResponseMap = new ConcurrentHashMap<>();
    static {
        Reflections reflections = new Reflections("ossrs.net.srssip.gb28181.event.*");
        initMessageHandler(reflections);
        initMessageRequestHandler(reflections);
        initMessageResponseHandler(reflections);
    }

    private static void initMessageHandler(Reflections reflections){
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(MessageEventHandler.class);
        for (Class<?> c : classes) {
            try {
                Object o = c.newInstance();
                if(o instanceof MessageEventAbstract){
                    MessageEventHandler handler = c.getAnnotation(MessageEventHandler.class);
                    messageEventMap.put(handler.value(), (MessageEventAbstract) o);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init message handler failed");
            }
        }

    }

    private static void initMessageRequestHandler(Reflections reflections){
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(MessageRequestEventHandler.class);
        classes.forEach(aClass -> {
            try {
                Object o = aClass.newInstance();
                if(o instanceof MessageRequestAbstract){
                    MessageRequestEventHandler handler = aClass.getAnnotation(MessageRequestEventHandler.class);
                    messageRequestMap.put(handler.value(), (MessageRequestAbstract) o);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init message request handler failed");
            }
        });
    }

    private static void initMessageResponseHandler(Reflections reflections){
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(MessageResponseEventHandler.class);
        classes.forEach(aClass -> {
            try {
                Object o = aClass.newInstance();
                if(o instanceof MessageResponseAbstract){
                    MessageResponseEventHandler handler = aClass.getAnnotation(MessageResponseEventHandler.class);
                    messageResponseMap.put(handler.value(), (MessageResponseAbstract) o);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init message response handler failed");
            }
        });
    }

    public MessageEventAbstract getMessageEvent(RequestEvent requestEvent){
        String method = requestEvent.getRequest().getMethod();
        MessageEventAbstract messageEventAbstract = messageEventMap.get(method);
        messageEventAbstract.constructor(requestEvent);
        return messageEventAbstract;
    }

    public MessageRequestAbstract getMessageRequest(String requestType, RequestEvent requestEvent){
        MessageRequestAbstract messageRequestAbstract = messageRequestMap.get(requestType);
        messageRequestAbstract.constructor(requestEvent);
        return messageRequestAbstract;
    }

    public MessageResponseAbstract getMessageResponse(String method, ResponseEvent responseEvent){
        MessageResponseAbstract messageResponseAbstract = messageResponseMap.get(method);
        messageResponseAbstract.constructor(responseEvent);
        return messageResponseAbstract;
    }
}
