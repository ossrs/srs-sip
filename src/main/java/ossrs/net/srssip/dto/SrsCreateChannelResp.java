package ossrs.net.srssip.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * @ Description ossrs.net.srssip.dto
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 11/3/2022 上午12:35
 */
@NoArgsConstructor
@Data
public class SrsCreateChannelResp {
    @JsonProperty("id")
    private String id;
    @JsonProperty("ip")
    private String ip;
    @JsonProperty("rtmp_port")
    private Integer rtmpPort;
    @JsonProperty("app")
    private String app;
    @JsonProperty("stream")
    private String stream;
    @JsonProperty("rtp_port")
    private Integer rtpPort;
    @JsonProperty("ssrc")
    private Integer ssrc;
    @JsonProperty("rtmp_url")
    private String rtmpUrl;
    @JsonProperty("port_mode")
    private String portMode;
    @JsonProperty("rtp_peer_port")
    private Integer rtpPeerPort;
    @JsonProperty("rtp_peer_ip")
    private String rtpPeerIp;
    @JsonProperty("recv_time")
    private Integer recvTime;
    @JsonProperty("recv_time_str")
    private LocalDateTime recvTimeStr;
    @JsonProperty("flv_url")
    private String flvUrl;
    @JsonProperty("hls_url")
    private String hlsUrl;
    @JsonProperty("webrtc_url")
    private String webrtcUrl;
}
