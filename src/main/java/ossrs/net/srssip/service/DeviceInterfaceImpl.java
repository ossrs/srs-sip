package ossrs.net.srssip.service;

import org.springframework.stereotype.Service;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;

import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * @ Description ossrs.net.srssip.gb28181.interfaces.impl
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午10:22
 */
@Service
public class DeviceInterfaceImpl implements IDeviceInterface {

    private ConcurrentHashMap<String, Device> DEVICE_LIST = new ConcurrentHashMap<>();

    private ConcurrentHashMap<String, DeviceChannel> DEVICE_CHANNEL_LIST = new ConcurrentHashMap<>();

    @Override
    public List<Device> list(int start, int limit, String q, boolean online) {
        return null;
    }

    @Override
    public Device info(String serial) {
        return null;
    }

    @Override
    public Device getById(String deviceId) {
        return DEVICE_LIST.get(deviceId);
    }

    @Override
    public boolean save(Device device) {
        return DEVICE_LIST.put(device.getId(), device) != null;
    }

    @Override
    public void saveDeviceChannel(List<DeviceChannel> deviceChannels) {
        DEVICE_CHANNEL_LIST.putAll(deviceChannels.stream().collect(Collectors.toMap(DeviceChannel::getChannelID, Function.identity())));
    }

    @Override
    public List<DeviceChannel> getDeviceChannels(String deviceId) {
        return DEVICE_CHANNEL_LIST.values()
                .stream()
                .filter(deviceChannel -> deviceChannel.getDeviceId().equals(deviceId))
                .collect(Collectors.toList());
    }
}
