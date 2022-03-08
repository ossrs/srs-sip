package ossrs.net.srssip.gb28181.event.request;

import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import javax.sip.RequestEvent;
import javax.sip.message.Request;
import java.io.UnsupportedEncodingException;

/**
 * @ Description ossrs.net.srssip.gb28181.event.message
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午12:50
 */
@Slf4j
public abstract class MessageRequestAbstract {

    @Getter
    public RequestEvent requestEvent;

    @Getter
    public String content;

    @Getter
    public MessageRequestAbstract messageRequestAbstract;

    public void process(){
        Request request = requestEvent.getRequest();
        try {
            this.content = new String(request.getRawContent(),"GBK");
        } catch (UnsupportedEncodingException e) {
            log.error("解析消息编码出错",e);
        }
    }

    public void constructor(RequestEvent requestEvent) {
        this.requestEvent = requestEvent;
        this.process();
    }
}
