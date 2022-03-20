package ossrs.net.srssip.gb28181.listeners.factory;

import lombok.extern.slf4j.Slf4j;
import org.reflections.Reflections;
import ossrs.net.srssip.gb28181.annotation.RequestEventHandler;
import ossrs.net.srssip.gb28181.annotation.MessageRequestHandler;
import ossrs.net.srssip.gb28181.annotation.ResponseEventHandler;
import ossrs.net.srssip.gb28181.event.request.RequestEventAbstract;
import ossrs.net.srssip.gb28181.event.request.message.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.event.response.ResponseEventAbstract;

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

    private static final Map<String, RequestEventAbstract> REQUEST_EVENT_MAP = new ConcurrentHashMap<>();

    private static final Map<String,Map<String,MessageRequestAbstract> > MESSAGE_REQUEST_MAP = new ConcurrentHashMap<>();

    private static final Map<String, ResponseEventAbstract> RESPONSE_EVENT_MAP = new ConcurrentHashMap<>();
    static {
        Reflections reflections = new Reflections("ossrs.net.srssip.gb28181.event.*");
        initMessageRequestEventHandler(reflections);
        initMessageRequest(reflections);
        initMessageResponseHandler(reflections);
    }

    private static void initMessageRequest(Reflections reflections) {
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(MessageRequestHandler.class);
        for (Class<?> c : classes) {
            Object bean = null;
            try {
                bean = c.newInstance();
                if (bean instanceof MessageRequestAbstract) {
                    MessageRequestHandler annotation = c.getAnnotation(MessageRequestHandler.class);
                    MESSAGE_REQUEST_MAP.putIfAbsent(annotation.type(),new ConcurrentHashMap<>());
                    MESSAGE_REQUEST_MAP.get(annotation.type())
                            .put(annotation.cmd(), (MessageRequestAbstract)bean);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init message request handler failed");
            }
        }
    }


    private static void initMessageRequestEventHandler(Reflections reflections){
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(RequestEventHandler.class);
        classes.forEach(aClass -> {
            try {
                Object o = aClass.newInstance();
                if(o instanceof RequestEventAbstract){
                    RequestEventHandler handler = aClass.getAnnotation(RequestEventHandler.class);
                    REQUEST_EVENT_MAP.put(handler.value(), (RequestEventAbstract) o);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init request event handler failed");
            }
        });
    }

    private static void initMessageResponseHandler(Reflections reflections){
        Set<Class<?>> classes = reflections.getTypesAnnotatedWith(ResponseEventHandler.class);
        classes.forEach(aClass -> {
            try {
                Object o = aClass.newInstance();
                if(o instanceof ResponseEventAbstract){
                    ResponseEventHandler handler = aClass.getAnnotation(ResponseEventHandler.class);
                    RESPONSE_EVENT_MAP.put(handler.value(), (ResponseEventAbstract) o);
                }
            } catch (InstantiationException | IllegalAccessException e) {
                log.error("init response event handler failed");
            }
        });
    }

    public RequestEventAbstract getRequestEvent(RequestEvent requestEvent){
        String method = requestEvent.getRequest().getMethod();
        RequestEventAbstract requestEventAbstract = REQUEST_EVENT_MAP.get(method);
        requestEventAbstract.constructor(requestEvent);
        return requestEventAbstract;
    }


    public ResponseEventAbstract getResponseEvent(String method, ResponseEvent responseEvent){
        ResponseEventAbstract responseEventAbstract = RESPONSE_EVENT_MAP.get(method);
        responseEventAbstract.constructor(responseEvent);
        return responseEventAbstract;
    }

    public MessageRequestAbstract getMessageRequest(String messageType, String cmdType, RequestEvent requestEvent) {
        MessageRequestAbstract requestAbstract = MESSAGE_REQUEST_MAP.get(messageType).get( cmdType);
        requestAbstract.constructor(requestEvent);
        return requestAbstract;
    }
}
