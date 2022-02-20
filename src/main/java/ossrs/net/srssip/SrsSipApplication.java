package ossrs.net.srssip;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import ossrs.net.srssip.config.SRSStreamCasterConfig;
import ossrs.net.srssip.config.SipConfig;

@SpringBootApplication
@EnableConfigurationProperties({SipConfig.class, SRSStreamCasterConfig.class})
public class SrsSipApplication {

    public static void main(String[] args) {
        SpringApplication.run(SrsSipApplication.class, args);
    }

}
