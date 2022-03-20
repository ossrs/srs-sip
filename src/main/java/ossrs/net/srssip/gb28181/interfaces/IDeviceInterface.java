package ossrs.net.srssip.gb28181.interfaces;

import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.event.response.InviteResponseEvent;

import java.util.List;

/**
 * @ Description ossrs.net.srssip.gb28181.interfaces
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午9:55
 */
public interface IDeviceInterface {
    /**
     * 查询设备列表
     * @param start 分页开始，从零开始
     * @param limit 分页大小
     * @param q 搜索关键字
     * @param online 在线状态\n
     *               允许值：true，false
     * @param sort
     * @param order
     * @return 设备列表
     */
    List<Device> list(int start, int limit, String q, boolean online, String sort, String order);

    /**
     * 查询指定设备信息
     * @param deviceId 设备国标号
     * @return 设备信息
     */
    Device getById(String deviceId);

    boolean save(Device device);

    void saveDeviceChannel(List<DeviceChannel> deviceChannels);

    List<DeviceChannel> getDeviceChannels(String deviceId);

    void putInviteResponseEvent(String deviceId, String channelsId, InviteResponseEvent inviteResponseEvent);

    boolean removeInviteResponseEvent(String deviceId, String channelsId);

    InviteResponseEvent getInviteResponseEvent(String deviceId, String channelsId);

    List<DeviceChannel> channellist(String serial, String code, String civilcode, String block, String channel_type, String dir_serial, Integer start, Integer limit, String q, Boolean online);
}
