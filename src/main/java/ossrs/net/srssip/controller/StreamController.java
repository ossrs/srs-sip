package ossrs.net.srssip.controller;

import io.swagger.annotations.ApiParam;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ossrs.net.srssip.gb28181.cmd.impl.SIPCommander;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import reactor.core.publisher.Mono;

import javax.annotation.Resource;

/**
 * @ Description ossrs.net.srssip.controller
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 13/3/2022 下午11:55
 */
@Slf4j
@RestController
@RequestMapping("/api/v1/stream")
public class StreamController {
    @Resource
    private SIPCommander sipCommander;
    @Resource
    private IDeviceInterface deviceInterface;

    @GetMapping("/start")
    public Mono<StreamInfo> start(@ApiParam String serial,
                                  @ApiParam(required = false) Integer channel,
                                  @ApiParam(required = false) String code,
                                  @ApiParam(required = false) String sms_id,
                                  @ApiParam(required = false) String sms_group_id,
                                  @ApiParam(required = false) String cdn,
                                  @ApiParam(required = false) String audio,
                                  @ApiParam(required = false) String transport,
                                  @ApiParam(required = false) String transport_mode,
                                  @ApiParam(required = false) Boolean check_channel_status,
                                  @ApiParam(required = false) Integer timeout) {
        Device device = deviceInterface.getById(serial);

        SipStreamPlayResponseSubscribe streamPlayResponseSubscribe  = new SipStreamPlayResponseSubscribe();

        return sipCommander.playStreamCmd(device,code,audio,transport,transport_mode,timeout,streamPlayResponseSubscribe);
    }
}
