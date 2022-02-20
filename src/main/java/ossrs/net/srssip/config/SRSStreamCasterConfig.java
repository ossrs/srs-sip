package ossrs.net.srssip.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

/**
 * @ Description ossrs.net.srssip.config
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午11:37
 */
@Data
@Component
@ConfigurationProperties(prefix = "srs.stream-caster",ignoreInvalidFields = true)
public class SRSStreamCasterConfig {
    private Integer rtpMuxPort = 9000;
    private Integer rtpIdleTimeout = 30;
    private boolean audioEnable = false;
}
