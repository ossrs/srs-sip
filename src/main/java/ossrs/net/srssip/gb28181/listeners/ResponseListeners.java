package ossrs.net.srssip.gb28181.listeners;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.config.SmsConfig;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.response.InviteResponseEvent;
import ossrs.net.srssip.gb28181.event.response.ResponseEventAbstract;
import ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;

import javax.annotation.Resource;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;

import static ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder.CALLBACK_CMD_PLAY;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/3/2022 下午1:51
 */
@Slf4j
@Component
public class ResponseListeners {

    @Resource
    private IDeviceInterface deviceInterface;
    @Resource
    private SipResponseHolder sipResponseHolder;
    @Resource
    private SmsConfig smsConfig;

    @EventListener
    public void inviteResponseEvent(InviteResponseEvent inviteResponseEvent){
        inviteResponseEvent.process();
        if(inviteResponseEvent.isSuccess()){
            String subscrieId = inviteResponseEvent.getDeviceId()+"@"+inviteResponseEvent.getChannelId();
            SipStreamPlayResponseSubscribe streamPlayResponseSubscribe  = (SipStreamPlayResponseSubscribe) sipResponseHolder
                    .getCallable(CALLBACK_CMD_PLAY,subscrieId);
            if(streamPlayResponseSubscribe!=null){
                String app = "gb28181";

                streamPlayResponseSubscribe.complete(
                        StreamInfo.builder()
                                .streamid(inviteResponseEvent.getStreamCode())
                                .smsid(smsConfig.getSerial())
                                .deviceid(inviteResponseEvent.getDeviceId())
                                .channelid(inviteResponseEvent.getChannelId())
                                .channelname("")
                                .channelcustomname("")
                                .webrtc(String.format("webrtc://%s:%s/%s/%s", smsConfig.getHost(), smsConfig.getHttpPort(),app,subscrieId))
                                .flv(String.format("%s://%s:%s/%s/%s.flv",smsConfig.getHTTPSPort()==0?"http":"https",  smsConfig.getHost(),smsConfig.getPort() ,app, subscrieId))
//                                .wsflv(String.format("%s://%s:8080/%s/%s.flv",gbsProtocol , host, channevo.getName()))
                                .rtmp(String.format("rtmp://%s:%s/%s/%s", smsConfig.getHost(),smsConfig.getRTMPPort(),app, subscrieId))
                                .hls(String.format("%s://%s:%s/%s/%s.m3u8",smsConfig.getHTTPSPort()==0?"http":"https" ,smsConfig.getHost(),smsConfig.getPort(), app, subscrieId))
                                .rtsp(String.format("rtsp://%s:554/%s/%s", smsConfig.getHost(),app, subscrieId))
                                .cdn(String.format("rtmp://%s:%s/%s/%s", smsConfig.getHost(), smsConfig.getRTMPPort(),app, subscrieId))
//                        .snapurl(channel.getSnapUrl())
                                .transport("UDP")
//                        .startat(LocalDateTime.now().minus(streamDTO.getLiveMs(), ChronoUnit.MILLIS)
//                                .format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")) )
//                        .recordstartat(LocalDateTime.now().minus(streamDTO.getLiveMs(), ChronoUnit.MILLIS)
//                                .format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")))
//                        .duration(streamDTO.getLiveMs().intValue())
//                        .sourcevideocodecname(streamDTO.getVideo().getCodec())
//                        .sourcevideowidth(streamDTO.getVideo().getWidth())
//                        .sourcevideoheight(streamDTO.getVideo().getHeight())
//                        .sourcevideoframerate(streamDTO.getFrames())
//                        .sourceaudiocodecname(streamDTO.getAudio().getCodec())
//                        .sourceaudiosamplerate(streamDTO.getAudio().getSampleRate())
                                .rtpcount(0)
                                .rtplostcount(0)
                                .rtplostrate(0)
//                        .videoframecount(streamDTO.getFrames())
//                        .audioenable(streamDTO.getAudio()!=null?1:0)
//                        .ondemand(channel.getOndemand()?1:0)
//                        .inbytes(streamDTO.getRecvBytes())
//                        .inbitrate(0L)
//                        .outbytes(streamDTO.getSendBytes())
//                        .numoutputs(channel.getNumoutputs())
//                        .cascadesize(0)
//                        .relaysize(0)
//                        .channelptztype(channel.getPtzType())
                                .build()
                );
            }
            if(inviteResponseEvent.getResponseEvent().getClientTransaction()!=null){
                deviceInterface.putInviteResponseEvent(
                        inviteResponseEvent.getDeviceId()
                        ,inviteResponseEvent.getChannelId()
                        ,inviteResponseEvent);
            }
        }
    }
}
