package ossrs.net.srssip.gb28181.domain;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

import javax.xml.bind.annotation.XmlElement;
import java.io.Serializable;

/**
 * @ Description ossrs.net.srssip.gb28181.domain
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 14/3/2022 下午11:09
 */
@AllArgsConstructor
@NoArgsConstructor
@Builder
@ApiModel(value = "设备通道对象", description = "")
public class DeviceChannel implements Serializable {

    @ApiModelProperty(value = "设备国标编号")
    private String deviceId;
    /**
     * <pre>
     *
     * </pre>
     */
    private String channelID;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "通道名称")
    private String Name;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "设备厂家")
    private String Manufacturer;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "设备型号")
    private String Model;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "设备归属")
    private String Owner;

    /**
     * <pre>
     * 行政区号
     * </pre>
     */
    @ApiModelProperty(value = "行政区域")
    private String CivilCode;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "警区")
    private String Block;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "安装地址")
    private String Address;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "当为设备时, 是否有子设备, 1-有,0-没有")
    private Integer Parental;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "父设备/区域/系统ID")
    private String ParentID;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "安全通道 0否 1是")
    private Integer SafetyWay;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "注册方式, 缺省为1, 1-IETF RFC3261, 2-基于口令的双向认证, 3-基于数字证书的双向认证\n" +
            "\n" +
            "允许值: 1, 2, 3")
    private String RegisterWay;

    /**
     * <pre>
     *
     * </pre>
     */

    private String CertNum;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "SSL安全认证 0否 1是")
    private Integer Certifiable;

    /**
     * <pre>
     *
     * </pre>
     */
    private String ErrCode;

    /**
     * <pre>
     *
     * </pre>
     */
    private String EndTime;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "保密属性, 缺省为0, 0-不涉密, 1-涉密\n" +
            "\n" +
            "允许值: 0, 1")
    private Integer Secrecy;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "设备/区域/系统IP地址")
    private String IPAddress;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "设备/区域/系统端口")
    private Integer Port;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "通道密码")
    private String Password;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "在线状态\n" +
            "\n" +
            "允许值: ON, OFF")
    private String Status;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "经度\n" +
            "\n" +
            "默认值: 0")
    private float Longitude;

    /**
     * <pre>
     *
     * </pre>
     */
    @ApiModelProperty(value = "纬度\n" +
            "\n" +
            "默认值: 0")
    private float Latitude;


    public String getChannelID() {
        return this.channelID;
    }

    @XmlElement(name = "DeviceID")
    public void setChannelID(String DeviceID) {
        this.channelID = DeviceID;
    }

    public String getName() {
        return this.Name;
    }

    @XmlElement(name = "Name")
    public void setName(String Name) {
        this.Name = Name;
    }

    public String getManufacturer() {
        return this.Manufacturer;
    }

    @XmlElement(name = "Manufacturer")
    public void setManufacturer(String Manufacturer) {
        this.Manufacturer = Manufacturer;
    }

    public String getModel() {
        return this.Model;
    }

    @XmlElement(name = "Model")
    public void setModel(String Model) {
        this.Model = Model;
    }

    public String getOwner() {
        return this.Owner;
    }

    @XmlElement(name = "Owner")
    public void setOwner(String Owner) {
        this.Owner = Owner;
    }

    public String getCivilCode() {
        return this.CivilCode;
    }

    @XmlElement(name = "CivilCode")
    public void setCivilCode(String CivilCode) {
        this.CivilCode = CivilCode;
    }

    public String getBlock() {
        return this.Block;
    }

    @XmlElement(name = "Block")
    public void setBlock(String Block) {
        this.Block = Block;
    }

    public String getAddress() {
        return this.Address;
    }

    @XmlElement(name = "Address")
    public void setAddress(String Address) {
        this.Address = Address;
    }

    public Integer getParental() {
        return this.Parental;
    }

    @XmlElement(name = "Parental")
    public void setParental(Integer Parental) {
        this.Parental = Parental;
    }

    public String getParentID() {
        return this.ParentID;
    }

    @XmlElement(name = "ParentID")
    public void setParentID(String ParentID) {
        this.ParentID = ParentID;
    }

    public Integer getSafetyWay() {
        return this.SafetyWay;
    }

    @XmlElement(name = "SafetyWay")
    public void setSafetyWay(Integer SafetyWay) {
        this.SafetyWay = SafetyWay;
    }

    public String getRegisterWay() {
        return this.RegisterWay;
    }

    @XmlElement(name = "RegisterWay")
    public void setRegisterWay(String RegisterWay) {
        this.RegisterWay = RegisterWay;
    }

    public String getCertNum() {
        return this.CertNum;
    }

    @XmlElement(name = "CertNum")
    public void setCertNum(String CertNum) {
        this.CertNum = CertNum;
    }

    public Integer getCertifiable() {
        return this.Certifiable;
    }

    @XmlElement(name = "Certifiable")
    public void setCertifiable(Integer Certifiable) {
        this.Certifiable = Certifiable;
    }

    public String getErrCode() {
        return this.ErrCode;
    }

    @XmlElement(name = "ErrCode")
    public void setErrCode(String ErrCode) {
        this.ErrCode = ErrCode;
    }

    public String getEndTime() {
        return this.EndTime;
    }

    @XmlElement(name = "EndTime")
    public void setEndTime(String EndTime) {
        this.EndTime = EndTime;
    }

    public Integer getSecrecy() {
        return this.Secrecy;
    }

    @XmlElement(name = "Secrecy")
    public void setSecrecy(Integer Secrecy) {
        this.Secrecy = Secrecy;
    }

    public String getIPAddress() {
        return this.IPAddress;
    }

    @XmlElement(name = "IPAddress")
    public void setIPAddress(String IPAddress) {
        this.IPAddress = IPAddress;
    }

    public Integer getPort() {
        return this.Port;
    }

    @XmlElement(name = "Port")
    public void setPort(Integer Port) {
        this.Port = Port;
    }

    public String getPassword() {
        return this.Password;
    }

    @XmlElement(name = "Password")
    public void setPassword(String Password) {
        this.Password = Password;
    }

    public String getStatus() {
        return this.Status;
    }

    @XmlElement(name = "Status")
    public void setStatus(String Status) {
        this.Status = Status;
    }

    public float getLongitude() {
        return this.Longitude;
    }

    @XmlElement(name = "Longitude")
    public void setLongitude(float Longitude) {
        this.Longitude = Longitude;
    }

    public float getLatitude() {
        return this.Latitude;
    }

    @XmlElement(name = "Latitude")
    public void setLatitude(float Latitude) {
        this.Latitude = Latitude;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }
}
