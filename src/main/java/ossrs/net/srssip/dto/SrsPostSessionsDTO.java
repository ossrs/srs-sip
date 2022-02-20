package ossrs.net.srssip.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @ Description ossrs.net.srssip.dto
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:16
 */
@NoArgsConstructor
@Data
public class SrsPostSessionsDTO {
    @JsonProperty("app")
    private String app;
    @JsonProperty("vhost")
    private String vhost;
    @JsonProperty("stream")
    private String stream;
    @JsonProperty("param")
    private String param;
    @JsonProperty("ip")
    private String ip;
    @JsonProperty("action")
    private String action;
    @JsonProperty("pageUrl")
    private String pageUrl;
    @JsonProperty("server_id")
    private String serverId;
    @JsonProperty("client_id")
    private String clientId;
}
