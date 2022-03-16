package ossrs.net.srssip.gb28181.event.request;

import lombok.extern.slf4j.Slf4j;
import ossrs.net.srssip.gb28181.annotation.MessageRequestEventHandler;
import ossrs.net.srssip.gb28181.domain.DeviceInfo;
import ossrs.net.srssip.gb28181.util.XmlUtil;

import javax.xml.bind.annotation.XmlElement;
import javax.xml.bind.annotation.XmlRootElement;

/**
 * @ Description ossrs.net.srssip.gb28181.event.request
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 14/3/2022 上午1:20
 */
@Slf4j
@XmlRootElement(name = "Response")
@MessageRequestEventHandler(value = "DeviceInfo")
public class DeviceInfoMessageRequest extends MessageRequestAbstract{

    private String deviceId;

    private String deviceName;

    private String sn;

    private String result;

    private String manufacturer;

    private String model;

    private String firmware;

    private Integer channel;

    private DeviceInfo deviceInfo;

    @Override
    public void process() {
        super.process();
        String content = this.content;
        DeviceInfoMessageRequest deviceInfoMessageRequest = (DeviceInfoMessageRequest) XmlUtil.xmlToObject(content,this);
        if(deviceInfoMessageRequest !=null){
            this.deviceInfo = DeviceInfo.builder()
                    .code(deviceInfoMessageRequest.deviceId)
                    .channelNum(deviceInfoMessageRequest.channel)
                    .name(deviceInfoMessageRequest.deviceName)
                    .firmware(deviceInfoMessageRequest.firmware)
                    .manufacturer(deviceInfoMessageRequest.manufacturer)
                    .model(deviceInfoMessageRequest.model)

                    .build();
        }
    }

    public String getDeviceId() {
        return deviceId;
    }

    @XmlElement(name = "DeviceID")
    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public String getDeviceName() {
        return deviceName;
    }

    @XmlElement(name = "DeviceName")
    public void setDeviceName(String deviceName) {
        this.deviceName = deviceName;
    }

    public String getSn() {
        return sn;
    }

    @XmlElement(name = "SN")
    public void setSn(String sn) {
        this.sn = sn;
    }

    public String getResult() {
        return result;
    }

    @XmlElement(name = "Result")
    public void setResult(String result) {
        this.result = result;
    }

    public String getManufacturer() {
        return manufacturer;
    }

    @XmlElement(name = "Manufacturer")
    public void setManufacturer(String manufacturer) {
        this.manufacturer = manufacturer;
    }

    public String getModel() {
        return model;
    }

    @XmlElement(name = "Model")
    public void setModel(String model) {
        this.model = model;
    }

    public String getFirmware() {
        return firmware;
    }

    @XmlElement(name = "Firmware")
    public void setFirmware(String firmware) {
        this.firmware = firmware;
    }

    public Integer getChannel() {
        return channel;
    }

    @XmlElement(name = "Channel")
    public void setChannel(Integer channel) {
        this.channel = channel;
    }

    public DeviceInfo getDeviceInfo() {
        return deviceInfo;
    }

    public void setDeviceInfo(DeviceInfo deviceInfo) {
        this.deviceInfo = deviceInfo;
    }
}
