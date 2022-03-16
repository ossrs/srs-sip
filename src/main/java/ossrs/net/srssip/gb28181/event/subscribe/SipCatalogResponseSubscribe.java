package ossrs.net.srssip.gb28181.event.subscribe;

import ossrs.net.srssip.gb28181.domain.DeviceChannel;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;

/**
 * @ Description ossrs.net.srssip.gb28181.event.subscribe
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 17/3/2022 上午12:24
 */
public class SipCatalogResponseSubscribe extends CompletableFuture<List<DeviceChannel>> {
    private String key;
    private String id;
    private List<DeviceChannel> channelListVo = new ArrayList<>();

    public SipCatalogResponseSubscribe() {
    }

    public SipCatalogResponseSubscribe(String key, String id) {
        this.key = key;
        this.id = id;
    }

    public boolean putChannelListVo(List<DeviceChannel> channelLists){
        return this.channelListVo.addAll(channelLists);
    }

    public List<DeviceChannel> getChannelListVo(){return this.channelListVo;}

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }
}
