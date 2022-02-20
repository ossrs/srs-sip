package ossrs.net.srssip.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @ Description ossrs.net.srssip.dto
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:17
 */
@NoArgsConstructor
@Data
public class SrsPostDvrsDTO {

    @JsonProperty("action")
    private String action;
    @JsonProperty("client_id")
    private Integer clientId;
    @JsonProperty("ip")
    private String ip;
    @JsonProperty("vhost")
    private String vhost;
    @JsonProperty("app")
    private String app;
    @JsonProperty("stream")
    private String stream;
    @JsonProperty("param")
    private String param;
    @JsonProperty("cwd")
    private String cwd;
    @JsonProperty("file")
    private String file;
}
