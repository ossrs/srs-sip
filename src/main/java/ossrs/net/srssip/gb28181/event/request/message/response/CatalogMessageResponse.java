package ossrs.net.srssip.gb28181.event.request.message.response;

import lombok.extern.slf4j.Slf4j;
import ossrs.net.srssip.gb28181.annotation.MessageRequestHandler;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.event.request.message.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.util.XmlUtil;

import javax.xml.bind.annotation.XmlAttribute;
import javax.xml.bind.annotation.XmlElement;
import javax.xml.bind.annotation.XmlRootElement;
import java.util.List;
import java.util.stream.Collectors;

/**
 * @ Description ossrs.net.srssip.gb28181.event.request
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 14/3/2022 下午10:42
 */
@Slf4j
@XmlRootElement(name = "Response")
@MessageRequestHandler(type = "Response", cmd = "Catalog")
public class CatalogMessageResponse extends MessageRequestAbstract {

    private String deviceId;

    private Integer sumNum;

    private String sn;

    private DeviceList deviceList;

    private List<DeviceChannel> deviceChannel;

    @Override
    public void process() {
        super.process();
        String content = this.content;
        CatalogMessageResponse catalogMessageResponse = (CatalogMessageResponse) XmlUtil.xmlToObject(content, this);
        if(catalogMessageResponse !=null){
            deviceChannel =  catalogMessageResponse.getDeviceList().getItem().stream().map(item ->
                    DeviceChannel.builder()
                            .Address(item.getAddress())
                            .CivilCode(item.getCivilCode())
                            .channelID(item.getDeviceId())
                            .deviceId(catalogMessageResponse.getDeviceId())
                            .Manufacturer(item.getManufacturer())
                            .Model(item.getModel())
                            .Name(item.getName())
                            .Owner(item.getOwner())
                            .RegisterWay(item.getRegisterWay())
                            .Secrecy(item.getSecrecy())
                            .Parental(item.getParental())
                            .ParentID(item.getParentID())
                            .SafetyWay(item.getSafetyWay())
                            .Status(item.getStatus())
                            .build()).collect(Collectors.toList());
        }
        this.messageRequestAbstract = catalogMessageResponse;
    }

    public String getDeviceId() {
        return deviceId;
    }

    @XmlElement(name = "DeviceID")
    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public Integer getSumNum() {
        return sumNum;
    }

    @XmlElement(name = "SumNum")
    public void setSumNum(Integer sumNum) {
        this.sumNum = sumNum;
    }

    public String getSn() {
        return sn;
    }

    @XmlElement(name = "SN")
    public void setSn(String sn) {
        this.sn = sn;
    }

    public DeviceList getDeviceList() {
        return deviceList;
    }

    @XmlElement(name = "DeviceList")
    public void setDeviceList(DeviceList deviceList) {
        this.deviceList = deviceList;
    }

    public List<DeviceChannel> getDeviceChannel() {
        return deviceChannel;
    }

    public void setDeviceChannel(List<DeviceChannel> deviceChannel) {
        this.deviceChannel = deviceChannel;
    }

    static class DeviceList {

        private Integer num;

        private List<DeviceChannel> item;

        public Integer getNum() {
            return num;
        }
        @XmlAttribute(name = "NUM")
        public void setNum(Integer num) {
            this.num = num;
        }

        public List<DeviceChannel> getItem() {
            return item;
        }
        @XmlElement(name = "Item")
        public void setItem(List<DeviceChannel> item) {
            this.item = item;
        }
    }
}
