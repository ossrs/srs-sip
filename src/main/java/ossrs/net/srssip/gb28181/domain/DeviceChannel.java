package ossrs.net.srssip.gb28181.domain;

import javax.xml.bind.annotation.XmlElement;
import java.io.Serializable;

/**
 * @ Description ossrs.net.srssip.gb28181.domain
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 14/3/2022 下午11:09
 */
public class DeviceChannel implements Serializable {

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
    private String	Name;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Manufacturer;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Model;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Owner;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	CivilCode;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Block;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Address;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Parental;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	ParentID;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	SafetyWay;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	RegisterWay;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	CertNum;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Certifiable;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	ErrCode;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	EndTime;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Secrecy;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	IPAddress;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Port;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Password;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Status;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Longitude;

    /**
     * <pre>
     *
     * </pre>
     */
    private String	Latitude;


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

    public String getParental() {
        return this.Parental;
    }
    @XmlElement(name = "Parental")
    public void setParental(String Parental) {
        this.Parental = Parental;
    }

    public String getParentID() {
        return this.ParentID;
    }
    @XmlElement(name = "ParentID")
    public void setParentID(String ParentID) {
        this.ParentID = ParentID;
    }

    public String getSafetyWay() {
        return this.SafetyWay;
    }
    @XmlElement(name = "SafetyWay")
    public void setSafetyWay(String SafetyWay) {
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

    public String getCertifiable() {
        return this.Certifiable;
    }
    @XmlElement(name = "Certifiable")
    public void setCertifiable(String Certifiable) {
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

    public String getSecrecy() {
        return this.Secrecy;
    }
    @XmlElement(name = "Secrecy")
    public void setSecrecy(String Secrecy) {
        this.Secrecy = Secrecy;
    }

    public String getIPAddress() {
        return this.IPAddress;
    }
    @XmlElement(name = "IPAddress")
    public void setIPAddress(String IPAddress) {
        this.IPAddress = IPAddress;
    }

    public String getPort() {
        return this.Port;
    }
    @XmlElement(name = "Port")
    public void setPort(String Port) {
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

    public String getLongitude() {
        return this.Longitude;
    }
    @XmlElement(name = "Longitude")
    public void setLongitude(String Longitude) {
        this.Longitude = Longitude;
    }

    public String getLatitude() {
        return this.Latitude;
    }
    @XmlElement(name = "Latitude")
    public void setLatitude(String Latitude) {
        this.Latitude = Latitude;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }
}
