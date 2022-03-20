package ossrs.net.srssip;

import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springdoc.core.GroupedOpenApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import ossrs.net.srssip.config.SRSStreamCasterConfig;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.config.SmsConfig;

@SpringBootApplication
@EnableConfigurationProperties({SipConfig.class, SmsConfig.class, SRSStreamCasterConfig.class})
public class SrsSipApplication {

    public static void main(String[] args) {
        SpringApplication.run(SrsSipApplication.class, args);
    }

}
