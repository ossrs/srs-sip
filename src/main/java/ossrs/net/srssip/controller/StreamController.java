package ossrs.net.srssip.controller;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
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
@Api(tags = "流信息")
@RestController
@RequestMapping("/api/v1/stream")
public class StreamController {
    @Resource
    private SIPCommander sipCommander;
    @Resource
    private IDeviceInterface deviceInterface;
    @Resource
    private SipResponseHolder sipResponseHolder;

    @ApiOperation(value = "开始直播")
    @GetMapping("/start")
    public Mono<StreamInfo> start(@ApiParam(required = true,value = "设备编号") String serial,
                                  @ApiParam(required = false,value = "通道序号\n" +
                                          "\n" +
                                          "默认值: 1") Integer channel,
                                  @ApiParam(required = true,value = "通道编号,通过 /api/v1/device/channellist 获取的 ChannelList.ID, \n" +
                                          "该参数和 channel 二选一传递即可") String code,
                                  @ApiParam(required = false,value = "指定SMS，默认取设备配置") String sms_id,
                                  @ApiParam(required = false,value = "指定SMS分组，默认取设备配置") String sms_group_id,
                                  @ApiParam(required = false,value = "转推 CDN 地址, 形如: [rtmp|rtsp]://xxx, encodeURIComponent") String cdn,
                                  @ApiParam(required = false,value = "是否开启音频, 默认 config 表示 读取通道音频开关配置\n" +
                                          "\n" +
                                          "默认值: config\n" +
                                          "\n" +
                                          "允许值: true, false, config") String audio,
                                  @ApiParam(required = false,value = "流传输模式， 默认 config 表示 读取设备流传输模式配置\n" +
                                          "\n" +
                                          "默认值: config\n" +
                                          "\n" +
                                          "允许值: TCP, UDP, config") String transport,
                                  @ApiParam(required = false,value = "当 transport=TCP 时有效, 指示流传输主被动模式, 默认被动\n" +
                                          "\n" +
                                          "默认值: passive\n" +
                                          "\n" +
                                          "允许值: active, passive") String transport_mode,
                                  @ApiParam(required = false,value = "是否检查通道状态, 默认 false, 表示 拉流前不检查通道状态是否在线\n" +
                                          "\n" +
                                          "默认值: false\n" +
                                          "\n" +
                                          "允许值: true, false") Boolean check_channel_status,
                                  @ApiParam(required = false,value = "拉流超时(秒), 默认使用 srs > sip > ack_timeout") Integer timeout) {
        Device device = deviceInterface.getById(serial);

        SipStreamPlayResponseSubscribe streamPlayResponseSubscribe  = new SipStreamPlayResponseSubscribe();
        streamPlayResponseSubscribe.setKey(CALLBACK_CMD_PLAY);
        streamPlayResponseSubscribe.setId(serial+"@"+code);
        sipResponseHolder.put(CALLBACK_CMD_PLAY,serial+"@"+code,streamPlayResponseSubscribe);
        return sipCommander.playStreamCmd(device,code,audio,transport,transport_mode,timeout,streamPlayResponseSubscribe);
    }

    @ApiOperation("直播流停止")
    @GetMapping("/stop")
    public String stop(@ApiParam(value = "设备编号") String serial,
                       @ApiParam(required = false,value = "通道序号\n" +
                               "\n" +
                               "默认值: 1") Integer channel,
                       @ApiParam(required = false,value = "通道编号,通过 /api/v1/device/channellist 获取的 ChannelList.ID, " +
                               "该参数和 channel 二选一传递即可") String code,
                       @ApiParam(required = false,value = "是否检查通道在线人数, 默认 false, 表示 停止前不检查通道是否有客户端正在播放\n" +
                               "\n" +
                               "默认值: false\n" +
                               "\n" +
                               "允许值: true, false") Boolean check_outputs) {
        log.info("关闭 {}", serial);
        sipCommander.streamByeCmd(serial,code);
        return "ok";
    }
}
