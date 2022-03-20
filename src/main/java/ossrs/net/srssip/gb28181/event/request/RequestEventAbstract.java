package ossrs.net.srssip.gb28181.event.request;

import lombok.Getter;

import javax.sip.RequestEvent;

/**
 * @ Description ossrs.net.srssip.gb28181.event
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午12:00
 */
public abstract class RequestEventAbstract {
    @Getter
    public RequestEvent requestEvent;

    public abstract void process();

    public void constructor(RequestEvent requestEvent) {
        this.requestEvent = requestEvent;
        this.process();
    }
}
