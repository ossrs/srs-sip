package ossrs.net.srssip.service;

import org.springframework.stereotype.Service;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.event.response.InviteResponseEvent;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
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

    private ConcurrentHashMap<String, Map<String, InviteResponseEvent>> INVITE_RESPONSE_EVENT_MAP = new ConcurrentHashMap<>();

    @Override
    public List<Device> list(int start, int limit, String q, boolean online, String sort, String order) {
        return new ArrayList<>(DEVICE_LIST.values());
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

    @Override
    public void putInviteResponseEvent(String deviceId, String channelsId, InviteResponseEvent inviteResponseEvent) {
        INVITE_RESPONSE_EVENT_MAP
                .putIfAbsent(deviceId,new ConcurrentHashMap<String,
                        InviteResponseEvent>(){{put(channelsId,inviteResponseEvent);}});
        INVITE_RESPONSE_EVENT_MAP.get(deviceId).put(channelsId,inviteResponseEvent);
    }

    @Override
    public boolean removeInviteResponseEvent(String deviceId, String channelsId) {
        return INVITE_RESPONSE_EVENT_MAP.get(deviceId).remove(channelsId)!=null;
    }

    @Override
    public InviteResponseEvent getInviteResponseEvent(String deviceId, String channelsId) {
        return INVITE_RESPONSE_EVENT_MAP.get(deviceId).get(channelsId);
    }

    @Override
    public List<DeviceChannel> channellist(String serial, String code, String civilcode, String block, String channel_type, String dir_serial, Integer start, Integer limit, String q, Boolean online) {
        return new ArrayList<>(DEVICE_CHANNEL_LIST.values());
    }


}
