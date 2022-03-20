package ossrs.net.srssip.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * @ Description ossrs.net.srssip.config
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/3/2022 上午1:08
 */
@Configuration
public class SwaggerConfig {
    @Bean
    public OpenAPI customOpenAPI(@Value("1.0") String appVersion) {
        return new OpenAPI()
                .info(new Info()
                        .title("SRS-SIP API v1.0")
                        .description("SRS Sip Server")
                        .version(appVersion)
                        .license(new License()
                                .name("MIT")
                                .url("https://github.com/ossrs/srs-sip/blob/main/LICENSE"))
                        .contact(new Contact()
                                .name("stormbirds")
                                .email("xbaojun@gmail.com")
                                .url("https://github.com/ossrs/srs-sip")));
    }
}
