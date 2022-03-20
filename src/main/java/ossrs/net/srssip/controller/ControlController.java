package ossrs.net.srssip.controller;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ossrs.net.srssip.gb28181.cmd.ISIPCommander;

import javax.annotation.Resource;

/**
 * @ Description ossrs.net.srssip.controller
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/3/2022 下午7:07
 */
@Slf4j
@RestController
@Api(value = "ControlController",tags = "设备控制")
@RequestMapping("/api/v1/control")
public class ControlController {

    @Resource
    private ISIPCommander sipCommander;

    @ApiOperation(value = "设备控制",notes = "设备控制 - 云台控制")
    @GetMapping("/ptz")
    public String ptz(@ApiParam(required = true,value = "设备编号") String serial,
                      @ApiParam(required = false,value = "通道序号\n" +
                              "\n" +
                              "默认值: 1") String channel,
                      @ApiParam(required = false,value = "通道编号,通过 /api/v1/device/channellist 获取的 ChannelList.ID, 该参数和 channel 二选一传递即可")
                              String code,
                      @ApiParam(required = true,value = "控制指令\n" +
                              "\n" +
                              "允许值: left, right, up, down, upleft, upright, downleft, downright, zoomin, zoomout, stop") String command,
                      @ApiParam(required = false,value = "速度(0~255)\n" +
                              "\n" +
                              "默认值: 129") Integer speed){
        return sipCommander.ptzAction(serial,channel,code,command,speed)?"ok":"failure";
    }
}
