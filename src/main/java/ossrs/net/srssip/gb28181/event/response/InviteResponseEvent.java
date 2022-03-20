package ossrs.net.srssip.gb28181.event.response;

import gov.nist.javax.sip.header.AddressParametersHeader;
import lombok.extern.slf4j.Slf4j;
import ossrs.net.srssip.gb28181.annotation.ResponseEventHandler;
import ossrs.net.srssip.gb28181.util.XmlUtil;

import javax.sip.address.SipURI;
import java.util.Map;

/**
 * @ Description ossrs.net.srssip.gb28181.event.response
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:01
 */
@Slf4j
@ResponseEventHandler("INVITE")
public class InviteResponseEvent extends ResponseEventAbstract {

    private String deviceId;

    private String channelId;

    private String streamCode;

    @Override
    public void process() {
        if(isSuccess()){

            responseAck();

            log.info("接收到拉流响应： {}", getContent());
            AddressParametersHeader header = (AddressParametersHeader) response.getHeader("To");
            SipURI sipURI = (SipURI)header.getAddress().getURI();
            this.channelId = sipURI.getUser();
            Map<String, String> contentMap = XmlUtil.convertStreamCode(getContent());
            if(contentMap!=null && contentMap.size()>0){
                String streamField = "y";
                this.streamCode = contentMap.get(streamField);
                String deviceField = "o";
                this.deviceId = contentMap.get(deviceField).substring(0,20);
            }
            log.info("目标客户端 {}", this.channelId);
        }

    }

    @Override
    public boolean isSuccess() {
        return super.isSuccess();
    }

    public String getDeviceId() {
        return deviceId;
    }

    public String getChannelId() {
        return channelId;
    }

    public String getStreamCode() {
        return streamCode;
    }
}
