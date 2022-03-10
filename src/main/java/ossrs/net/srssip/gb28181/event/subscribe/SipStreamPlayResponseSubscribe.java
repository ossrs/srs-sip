package ossrs.net.srssip.gb28181.event.subscribe;

import ossrs.net.srssip.gb28181.domain.StreamInfo;

import java.util.concurrent.CompletableFuture;

/**
 * @ Description ossrs.net.srssip.gb28181.event.subscribe
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 11/3/2022 上午12:25
 */
public class SipStreamPlayResponseSubscribe extends CompletableFuture<StreamInfo> {
    private String key;
    private String id;
    private StreamInfo streamInfo = new StreamInfo();

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

    public StreamInfo getStreamInfo() {
        return streamInfo;
    }

    public void setStreamInfo(StreamInfo streamInfo) {
        this.streamInfo = streamInfo;
    }
}
