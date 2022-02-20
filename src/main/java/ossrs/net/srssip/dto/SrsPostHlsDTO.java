package ossrs.net.srssip.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @ Description ossrs.net.srssip.dto
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:18
 */
@NoArgsConstructor
@Data
public class SrsPostHlsDTO {

    @JsonProperty("server_id")
    private String serverId;
    @JsonProperty("action")
    private String action;
    @JsonProperty("client_id")
    private String clientId;
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
    @JsonProperty("duration")
    private Double duration;
    @JsonProperty("cwd")
    private String cwd;
    @JsonProperty("file")
    private String file;
    @JsonProperty("url")
    private String url;
    @JsonProperty("m3u8")
    private String m3u8;
    @JsonProperty("m3u8_url")
    private String m3u8Url;
    @JsonProperty("seq_no")
    private Integer seqNo;
}
