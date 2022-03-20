package ossrs.net.srssip.gb28181.event.request.message.query;

import lombok.extern.slf4j.Slf4j;
import ossrs.net.srssip.gb28181.annotation.MessageRequestHandler;
import ossrs.net.srssip.gb28181.event.request.message.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.util.XmlUtil;

import javax.xml.bind.annotation.XmlElement;
import javax.xml.bind.annotation.XmlRootElement;

/**
 * @ Description ossrs.net.srssip.gb28181.event
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 19/3/2022 下午12:49
 */
@Slf4j
@XmlRootElement(name = "Query")
@MessageRequestHandler(type="Query", cmd = "Catalog")
public class CatalogMessageQuery extends MessageRequestAbstract {
    private String deviceId;

    private String sn;
    @Override
    public void process() {
        super.process();
        this.messageRequestAbstract = (CatalogMessageQuery) XmlUtil.xmlToObject(content, this);
    }

    public String getDeviceId() {
        return deviceId;
    }

    @XmlElement(name = "DeviceID")
    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public String getSn() {
        return sn;
    }

    @XmlElement(name = "SN")
    public void setSn(String sn) {
        this.sn = sn;
    }
}
