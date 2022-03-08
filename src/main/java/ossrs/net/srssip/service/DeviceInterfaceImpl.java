package ossrs.net.srssip.service;

import org.springframework.stereotype.Service;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;

import java.util.List;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @ Description ossrs.net.srssip.gb28181.interfaces.impl
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午10:22
 */
@Service
public class DeviceInterfaceImpl implements IDeviceInterface {

    private ConcurrentHashMap<String,Device> DEVICE_LIST  = new ConcurrentHashMap<>();

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
}
