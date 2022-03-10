package ossrs.net.srssip.gb28181.event.subscribe;

import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @ Description ossrs.net.srssip.gb28181.event.subscribe
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 11/3/2022 上午12:31
 */
@Component
public class SipResponseHolder {
    public static final String CALLBACK_CMD_DEVICESTATUS = "CALLBACK_DEVICESTATUS";

    public static final String CALLBACK_CMD_DEVICEINFO = "CALLBACK_DEVICEINFO";

    public static final String CALLBACK_CMD_DEVICECONTROL = "CALLBACK_DEVICECONTROL";

    public static final String CALLBACK_CMD_DEVICECONFIG = "CALLBACK_DEVICECONFIG";

    public static final String CALLBACK_CMD_CONFIGDOWNLOAD = "CALLBACK_CONFIGDOWNLOAD";

    public static final String CALLBACK_CMD_CATALOG = "CALLBACK_CATALOG";

    public static final String CALLBACK_CMD_RECORDINFO = "CALLBACK_RECORDINFO";

    public static final String CALLBACK_CMD_PLAY = "CALLBACK_PLAY";

    public static final String CALLBACK_CMD_PLAYBACK = "CALLBACK_PLAY";

    public static final String CALLBACK_CMD_DOWNLOAD = "CALLBACK_DOWNLOAD";

    public static final String CALLBACK_CMD_STOP = "CALLBACK_STOP";

    public static final String CALLBACK_CMD_MOBILEPOSITION = "CALLBACK_MOBILEPOSITION";

    public static final String CALLBACK_CMD_PRESETQUERY = "CALLBACK_PRESETQUERY";

    public static final String CALLBACK_CMD_ALARM = "CALLBACK_ALARM";

    public static final String CALLBACK_CMD_BROADCAST = "CALLBACK_BROADCAST";

    private Map<String, Map<String, CompletableFuture>> map = new ConcurrentHashMap<>();

    public CompletableFuture getCallable(String key, String id) {
        Map<String, CompletableFuture> callableMap = map.get(key);
        if(callableMap==null) return null;
        return map.get(key).get(id);
    }

    public void put(String key, String id, CompletableFuture future) {
        Map<String, CompletableFuture> callableMap = map.get(key);
        if(callableMap == null){
            callableMap = new ConcurrentHashMap<>();
            map.put(key,callableMap);
        }
        callableMap.put(id,future);
    }

    public void remove(String key, String id) {

    }

    public void complete(String key, String id, Object message) {
        Map<String, CompletableFuture> callableMap = map.get(key);
        if(callableMap == null){
            return;
        }
        CompletableFuture completableFuture = callableMap.get(id);
        if(completableFuture==null){
            return;
        }
        completableFuture.complete(message);
        callableMap.remove(id);
        if(callableMap.size()==0){
            map.remove(key);
        }
    }
}
