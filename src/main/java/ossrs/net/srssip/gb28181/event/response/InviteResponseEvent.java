package ossrs.net.srssip.gb28181.event.response;

import gov.nist.javax.sip.header.AddressParametersHeader;
import lombok.extern.slf4j.Slf4j;
import ossrs.net.srssip.gb28181.annotation.MessageResponseEventHandler;

import javax.sip.address.SipURI;
import java.util.Map;
import java.util.Objects;

/**
 * @ Description ossrs.net.srssip.gb28181.event.response
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:01
 */
@Slf4j
@MessageResponseEventHandler("INVITE")
public class InviteResponseEvent extends MessageResponseAbstract{
    @Override
    public void process() {
        responseAck();
        log.info("接收到拉流响应： {}", getContent());
        AddressParametersHeader header = (AddressParametersHeader) response.getHeader("To");
        SipURI sipURI = (SipURI)header.getAddress().getURI();
        String user = sipURI.getUser();
        log.info("目标客户端 {}", user);

    }
}
