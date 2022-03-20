package ossrs.net.srssip.gb28181.domain;

import lombok.Data;

import java.time.LocalDateTime;

/**
 * @ Description ossrs.net.srssip.gb28181.domain
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午10:04
 */
@Data
public class Device {
    /**
     * 国标编号
     */
    private String id;
    /**
     * 设备名称
     */
    private String name;
    /**
     * 设备厂商
     */
    private String manufacturer;
    /**
     * 设备类型,固定为国标设备"GB"
     */
    private String type;
    /**
     * 通道数
     */
    private Integer channelCount;
    /**
     * 指定到 SRS ID, 为空时, 自动选择负载最小的 SRS
     */
    private String srsId;
    /**
     * 指定到 SRS 分组
     */
    private String srsGroupId;
    /**
     * 通道目录抓取周期(秒), 0 表示后台不周期抓取目录
     * 默认值: 3600
     */
    private Integer catalogInterval;
    /**
     * 订阅周期(秒), 0 表示后台不周期订阅
     * 默认值: 0
     */
    private Integer subscribeInterval;
    /**
     * 目录订阅是否开启
     */
    private boolean catalogSubscribe;
    /**
     * 报警订阅是否开启
     */
    private boolean alarmSubscribe;
    /**
     * 位置订阅是否开启
     */
    private boolean positionSubscribe;
    /**
     * 是否在线
     */
    private boolean online;
    /**
     * 接入密码, 为空时默认使用统一接入密码 SipConfig.password
     */
    private String password;
    /**
     * 信令传输模式
     * 允许值: UDP, TCP
     */
    private String commandTransport;
    /**
     * 流传输模式
     * 允许值: UDP, TCP
     */
    private String mediaTransport;
    /**
     * 出口IP
     */
    private String remoteIP;
    /**
     * 端口
     */
    private Integer remotePort;
    /**
     * 经度 默认值：0
     */
    private float longitude;
    /**
     * 维度 默认值：0
     */
    private float latitude;
    /**
     * 最近注册
     */
    private LocalDateTime lastRegisterAt;
    /**
     * 最近心跳
     */
    private LocalDateTime lastKeepaliveAt;
    /**
     * 更新时间
     */
    private LocalDateTime updatedAt;
    /**
     * 创建时间
     */
    private LocalDateTime createdAt;
}
