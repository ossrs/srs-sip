package ossrs.net.srssip.gb28181.listeners.factory;

import ossrs.net.srssip.gb28181.event.messageevent.MessageEventAbstract;
import ossrs.net.srssip.gb28181.event.message.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.event.response.MessageResponseAbstract;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午12:35
 */
public class MessageEventFactory {

    public static MessageEventFactory INSTANCE = new MessageEventFactory();

    private static final Map<String, MessageEventAbstract> messageEventMap = new ConcurrentHashMap<>();

    private static final Map<String, MessageRequestAbstract> messageRequestMap = new ConcurrentHashMap<>();

    private static final Map<String, MessageResponseAbstract> messageResponseMap = new ConcurrentHashMap<>();
}
