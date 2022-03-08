package ossrs.net.srssip.gb28181.interfaces;

import ossrs.net.srssip.gb28181.domain.Device;

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
     * @return 设备列表
     */
    List<Device> list(int start, int limit, String q, boolean online);

    /**
     * 查询指定设备信息
     * @param serial 设备国标号
     * @return 设备信息
     */
    Device info(String serial);

    Device getById(String deviceId);
}
