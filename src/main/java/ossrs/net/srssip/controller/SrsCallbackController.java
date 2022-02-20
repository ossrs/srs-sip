package ossrs.net.srssip.controller;

import cn.dev33.satoken.stp.StpUtil;
import com.alibaba.fastjson.JSON;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ossrs.net.srssip.dto.*;

/**
 * @ Description SRS服务回调接口
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:07
 */
@Slf4j
@RestController
@RequestMapping("/srscallback")
public class SrsCallbackController {

    /**
     * on_connect on_close 设备连接 断开连接 动作
     *
     * @param postClientsDTO
     * @return
     */
    @PostMapping("/clients")
    public int clients(@RequestBody SrsPostClientsDTO postClientsDTO) {
        log.info("Clients 事件 :{}", postClientsDTO.toString());
        return StpUtil.isLogin() ? 0 : 401;
    }

    /**
     * on_publish on_unpublish 推流 停止推流 动作
     *
     * @param postStreamsDTO
     * @return
     */
    @PostMapping("/streams")
    public int streams(@RequestBody SrsPostStreamsDTO postStreamsDTO) {
        log.info("Streams 事件：{}", postStreamsDTO);
        return StpUtil.isLogin() ? 0 : 401;
    }

    /**
     * on_play 播放流 on_stop 停止播放流
     *
     * @param postSessionsDTO
     * @return
     */
    @PostMapping("/sessions")
    public int sessions(@RequestBody SrsPostSessionsDTO postSessionsDTO) {
        log.info("Sessions 事件：{}", postSessionsDTO);
        if ("on_stop".equals(postSessionsDTO.getAction())) {

            log.info("接收到停止播放事件 {}", "streamId");
        } else if ("on_play".equals(postSessionsDTO.getAction())) {

        }
        return StpUtil.isLogin() ? 0 : 401;
    }

    @PostMapping("/dvrs")
    public int dvrs(@RequestBody SrsPostDvrsDTO postDvrsDTO) {
        log.info("dvrs {}", JSON.toJSONString(postDvrsDTO));
        return 0;
    }

    @PostMapping("/hls")
    public int hls(@RequestBody SrsPostHlsDTO postHlsDTO) {
        log.info("hls {}", JSON.toJSONString(postHlsDTO));
        return 0;
    }
}
