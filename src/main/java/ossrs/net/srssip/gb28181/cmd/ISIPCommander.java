package ossrs.net.srssip.gb28181.cmd;

import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import reactor.core.publisher.Mono;

/**
 * @ Description ossrs.net.srssip.gb28181.cmd
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 下午10:35
 */
public interface ISIPCommander {
    Mono<StreamInfo> playStreamCmd(Device device, String channelId,
                                   String transport, String transport_mode, Integer timeout,
                                   SipStreamPlayResponseSubscribe streamPlayResponseSubscribe);
    void streamByeCmd(String deviceId, String channelId);
}
