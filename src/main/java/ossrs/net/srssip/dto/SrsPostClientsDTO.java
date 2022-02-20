package ossrs.net.srssip.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @ Description ossrs.net.srssip.dto
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:14
 */
@NoArgsConstructor
@Data
public class SrsPostClientsDTO {
    @JsonProperty("app")
    private String app;
    @JsonProperty("tcUrl")
    private String tcUrl;
    @JsonProperty("vhost")
    private String vhost;
    @JsonProperty("ip")
    private String ip;
    @JsonProperty("action")
    private String action;
    @JsonProperty("pageUrl")
    private String pageUrl;
    @JsonProperty("client_id")
    private String clientId;
}

