package ossrs.net.srssip.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

import javax.validation.constraints.NotEmpty;

/**
 * @ Description ossrs.net.srssip.config
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 10/3/2022 下午11:52
 */
@Data
@Component
@ConfigurationProperties(prefix = "srs.sms", ignoreInvalidFields = true)
public class SmsConfig {
    @NotEmpty(message = "serial 不能为空")
    private String serial;
    @NotEmpty(message = "realm 不能为空")
    private String realm;

    private String ip;
    private Integer httpPort;
    private String wanIp;
    private Boolean use_wan_ip_recv_stream;

    private String host;
    private Integer Port;

    private Integer RTMPPort;
    private Integer HTTPSPort;

    private Rtp rtp;

    @Data
    public static class Rtp {
        private Integer mux_port;
        private String tcp_port_range;
        private String udp_port_range;
        private Integer idle_timeout_seconds;
    }
}
