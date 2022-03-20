package ossrs.net.srssip.gb28181.cmd;

import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.subscribe.SipCatalogResponseSubscribe;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import reactor.core.publisher.Mono;

import java.util.List;

/**
 * @ Description ossrs.net.srssip.gb28181.cmd
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 下午10:35
 */
public interface ISIPCommander {
    Mono<StreamInfo> playStreamCmd(Device device, String channelId, String audio,
                                   String transport, String transport_mode, Integer timeout,
                                   SipStreamPlayResponseSubscribe streamPlayResponseSubscribe);
    void streamByeCmd(String deviceId, String channelId);

    Mono<List<DeviceChannel>> catalogQuery(Device device, Integer sn, Integer timeout, SipCatalogResponseSubscribe catalogResponseSubscribe);

    boolean ptzAction(String serial, String channel, String code, String command, Integer speed);
}
