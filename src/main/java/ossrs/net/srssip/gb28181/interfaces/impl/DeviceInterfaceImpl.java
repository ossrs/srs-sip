package ossrs.net.srssip.gb28181.interfaces.impl;

import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;

import java.util.List;

/**
 * @ Description ossrs.net.srssip.gb28181.interfaces.impl
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午10:22
 */
public class DeviceInterfaceImpl implements IDeviceInterface {
    @Override
    public List<Device> list(int start, int limit, String q, boolean online) {
        return null;
    }

    @Override
    public Device info(String serial) {
        return null;
    }
}
