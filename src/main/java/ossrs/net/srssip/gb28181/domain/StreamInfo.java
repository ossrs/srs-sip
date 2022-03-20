package ossrs.net.srssip.gb28181.domain;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModel;
import lombok.*;

import java.io.Serializable;

/**
 * @ Description ossrs.net.srssip.gb28181.domain
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 下午10:50
 */
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Data
@EqualsAndHashCode(callSuper = false)
@ApiModel("媒体流信息")
public class StreamInfo implements Serializable {

    private static final long serialVersionUID = -7132538169020897684L;
    @JsonProperty("AudioEnable")
    private Integer audioenable;

    @JsonProperty("CDN")
    private String cdn;

    @JsonProperty("CascadeSize")
    private Integer cascadesize;

    @JsonProperty("ChannelCustomName")
    private String channelcustomname;

    @JsonProperty("ChannelID")
    private String channelid;

    @JsonProperty("ChannelName")
    private String channelname;

    @JsonProperty("ChannelPTZType")
    private Integer channelptztype;

    @JsonProperty("DeviceID")
    private String deviceid;

    @JsonProperty("Duration")
    private Integer duration;

    @JsonProperty("FLV")
    private String flv;

    @JsonProperty("HLS")
    private String hls;

    @JsonProperty("InBitRate")
    private Long inbitrate;

    @JsonProperty("InBytes")
    private Long inbytes;

    @JsonProperty("NumOutputs")
    private Integer numoutputs;

    @JsonProperty("Ondemand")
    private Integer ondemand;

    @JsonProperty("OutBytes")
    private Long outbytes;

    @JsonProperty("RTMP")
    private String rtmp;

    @JsonProperty("RTPCount")
    private Integer rtpcount;

    @JsonProperty("RTPLostCount")
    private Integer rtplostcount;

    @JsonProperty("RTPLostRate")
    private Integer rtplostrate;

    @JsonProperty("RTSP")
    private String rtsp;

    @JsonProperty("RecordStartAt")
    private String recordstartat;

    @JsonProperty("RelaySize")
    private Integer relaysize;

    @JsonProperty("SMSID")
    private String smsid;

    @JsonProperty("SnapURL")
    private String snapurl;

    @JsonProperty("SourceAudioCodecName")
    private String sourceaudiocodecname;

    @JsonProperty("SourceAudioSampleRate")
    private Integer sourceaudiosamplerate;

    @JsonProperty("SourceVideoCodecName")
    private String sourcevideocodecname;

    @JsonProperty("SourceVideoFrameRate")
    private Long sourcevideoframerate;

    @JsonProperty("SourceVideoHeight")
    private Integer sourcevideoheight;

    @JsonProperty("SourceVideoWidth")
    private Integer sourcevideowidth;

    @JsonProperty("StartAt")
    private String startat;

    @JsonProperty("StreamID")
    private String streamid;

    @JsonProperty("Transport")
    private String transport;

    @JsonProperty("VideoFrameCount")
    private Long videoframecount;

    @JsonProperty("WEBRTC")
    private String webrtc;

    @JsonProperty("WS_FLV")
    private String wsflv;

    @JsonProperty("PlaybackDuration")
    private Integer playbackduration;

    @JsonProperty("PlaybackFileURL")
    private String playbackfileurl;

    @JsonProperty("Progress")
    private Float progress;

    @JsonProperty("TimestampSec")
    private Integer timestampsec;


}
