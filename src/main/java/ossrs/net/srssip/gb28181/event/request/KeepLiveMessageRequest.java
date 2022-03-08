package ossrs.net.srssip.gb28181.event.request;

import ossrs.net.srssip.gb28181.annotation.MessageRequestEventHandler;
import ossrs.net.srssip.gb28181.util.XmlUtil;

import javax.xml.bind.annotation.XmlElement;
import javax.xml.bind.annotation.XmlRootElement;

/**
 * @ Description ossrs.net.srssip.gb28181.event.message
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 上午2:25
 */
@XmlRootElement(name = "Notify")
@MessageRequestEventHandler("Keepalive")
public class KeepLiveMessageRequest extends MessageRequestAbstract {
    private String cmdType;

    private String serialNumber;

    private String deviceId;

    private String status;

    @Override
    public void process() {
        super.process();
        String content = this.content;
        this.messageRequestAbstract = (KeepLiveMessageRequest) XmlUtil.xmlToObject(content, this);
    }

    public String getCmdType() {
        return cmdType;
    }

    @XmlElement(name = "CmdType")
    public void setCmdType(String cmdType) {
        this.cmdType = cmdType;
    }

    @XmlElement(name = "SN")
    public String getSerialNumber() {
        return serialNumber;
    }

    public void setSerialNumber(String serialNumber) {
        this.serialNumber = serialNumber;
    }

    public String getDeviceId() {
        return deviceId;
    }

    @XmlElement(name = "DeviceID")
    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public String getStatus() {
        return status;
    }

    @XmlElement(name = "Status")
    public void setStatus(String status) {
        this.status = status;
    }
}
