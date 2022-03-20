package ossrs.net.srssip.controller;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import reactor.core.publisher.Mono;

import javax.annotation.Resource;
import java.util.List;

/**
 * @ Description ossrs.net.srssip.controller
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/3/2022 下午8:46
 */
@Slf4j
@Api(tags = "设备信息")
@RestController
@RequestMapping("/api/v1/device")
public class DevicesController {

    @Resource
    private IDeviceInterface deviceInterface;

    /**
     * 设备信息 - 查询设备列表
     *
     * @param start  分页开始,从零开始
     * @param limit  分页大小
     * @param q      搜索关键字
     * @param online 在线状态
     *               允许值: true, false
     * @return
     */
    @ApiOperation(value = "设备信息 - 查询设备列表")
    @GetMapping("/list")
    public List<Device> list(@ApiParam(required = false,value = "分页开始,从零开始") Integer start,
                             @ApiParam(required = false,value = "分页大小") Integer limit,
                             @ApiParam(required = false,value = "搜索关键字") String q,
                             @ApiParam(required = false,value = "在线状态\n" +
                                     "\n" +
                                     "允许值: true, false") Boolean online,
                             @ApiParam(required = false,value = "") String sort,
                             @ApiParam(required = false,value = "") String order) {

        if(limit==null) limit=10;
        if(start==null) start =1;

        return deviceInterface.list(start,limit,q,online,sort,order);
    }

    /**
     * 设备信息 - 查询单条设备信息
     *
     * @param serial 设备编号
     * @return
     */
    @ApiOperation(value = "设备信息 - 查询单条设备信息")
    @GetMapping("/info")
    public Device info(@ApiParam(value = "设备编号\n" +
            "\n",required = true) String serial) {
        return deviceInterface.getById(serial);
    }

    /**
     * 获取国标设备通道列表
     * @param serial 设备编号
     * @param code 通道编号
     * @param civilcode
     * @param block
     * @param channel_type
     * @param dir_serial
     * @param start
     * @param limit
     * @param q
     * @param online
     * @return
     */
    @GetMapping("/channellist")
    public List<DeviceChannel> channellist(@ApiParam(required = false) String serial,
                                           @ApiParam(required = false) String code,
                                           @ApiParam(required = false) String civilcode,
                                           @ApiParam(required = false) String block,
                                           @ApiParam(required = false) String channel_type,
                                           @ApiParam(required = false) String dir_serial,
                                           @ApiParam(required = false) Integer start,
                                           @ApiParam(required = false) Integer limit,
                                           @ApiParam(required = false) String q,
                                           @ApiParam(required = false) Boolean online) {

        return deviceInterface.channellist(serial, code, civilcode, block, channel_type, dir_serial, start, limit, q, online);
    }
}
