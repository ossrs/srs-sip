package ossrs.net.srssip.controller;

import io.swagger.annotations.ApiParam;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ossrs.net.srssip.gb28181.cmd.impl.SIPCommander;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import reactor.core.publisher.Mono;

import javax.annotation.Resource;

import static ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder.CALLBACK_CMD_PLAY;

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
    @Resource
    private SipResponseHolder sipResponseHolder;

    @GetMapping("/start")
    public Mono<StreamInfo> start(@ApiParam(required = true) String serial,
                                  @ApiParam(required = false) Integer channel,
                                  @ApiParam(required = true) String code,
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
        streamPlayResponseSubscribe.setKey(CALLBACK_CMD_PLAY);
        streamPlayResponseSubscribe.setId(serial+"@"+code);
        sipResponseHolder.put(CALLBACK_CMD_PLAY,serial+"@"+code,streamPlayResponseSubscribe);
        return sipCommander.playStreamCmd(device,code,audio,transport,transport_mode,timeout,streamPlayResponseSubscribe);
    }

    @GetMapping("/stop")
    public String stop(@ApiParam String serial,
                       @ApiParam(required = false) Integer channel,
                       @ApiParam(required = false) String code,
                       @ApiParam(required = false) Boolean check_outputs) {
        log.info("关闭 {}", serial);
        sipCommander.streamByeCmd(serial,code);
        return "ok";
    }
}
