package ossrs.net.srssip.gb28181.cmd.impl;

import com.alibaba.fastjson.JSONObject;
import gov.nist.javax.sip.message.SIPRequest;
import gov.nist.javax.sip.stack.SIPDialog;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
import org.springframework.web.reactive.function.client.WebClient;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.config.SmsConfig;
import ossrs.net.srssip.dto.SrsCreateChannelResp;
import ossrs.net.srssip.gb28181.cmd.ISIPCommander;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.domain.StreamInfo;
import ossrs.net.srssip.gb28181.event.response.InviteResponseEvent;
import ossrs.net.srssip.gb28181.event.subscribe.SipCatalogResponseSubscribe;
import ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder;
import ossrs.net.srssip.gb28181.event.subscribe.SipStreamPlayResponseSubscribe;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import reactor.core.publisher.Mono;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.address.Address;
import javax.sip.address.SipURI;
import javax.sip.header.*;
import javax.sip.message.Request;

import java.text.ParseException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

import static ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder.CALLBACK_CMD_PLAY;

/**
 * @ Description ossrs.net.srssip.gb28181.cmd.impl
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 下午10:56
 */
@Slf4j
@Service
public class SIPCommander implements ISIPCommander {

    private final static String CREATE_CHANNEL_URL = "http://%s:%d/api/v1/gb28181?action=create_channel&id=%s&stream=%s&port_mode=fixed&app=%s";

    private final static String DELETE_CHANNEL_URL = "http://%s:%d/api/v1/gb28181?action=delete_channel&id=%s&chid=%s";

    @Resource
    private SipConfig sipConfig;
    @Resource
    private SmsConfig smsConfig;
    @Resource
    private SipResponseHolder sipResponseHolder;
    @Value(value = "${user-settings.senior_sdp}")
    private Boolean isSeniorSdp;
    @Resource
    private SipFactory sipFactory;
    @Resource
    private SipProvider sipTcpProvider;
    @Resource
    private SipProvider sipUdpProvider;
    @Resource
    private IDeviceInterface deviceInterface;

    @Override
    public Mono<StreamInfo> playStreamCmd(Device device, String channelId, String audio, String transport,
                                          String transport_mode, Integer timeout,
                                          SipStreamPlayResponseSubscribe streamPlayResponseSubscribe) {
        if (device == null) return Mono.justOrEmpty(new StreamInfo());
        String key = CALLBACK_CMD_PLAY + device.getId() + channelId;
        String url = String.format(CREATE_CHANNEL_URL,
                smsConfig.getHost(), smsConfig.getHttpPort(),
                device.getId().concat("@").concat(channelId),
                device.getId().concat("@").concat(channelId), "gb28181");
        return WebClient.create(url)
                .get()
                .retrieve()
                .bodyToMono(JSONObject.class)
                .flatMap(jsonObject -> {
                    if (jsonObject.getInteger("code") == 0) {
                        SrsCreateChannelResp srsCreateChannelResp = jsonObject
                                .getJSONObject("data")
                                .getObject("query", SrsCreateChannelResp.class);
                        sipResponseHolder.put(key, String.valueOf(srsCreateChannelResp.getSsrc()), streamPlayResponseSubscribe);
                        String streamMode;
                        if ("UDP".equals(transport)) {
                            streamMode = transport;
                        } else if ("TCP".equals(transport)) {
                            streamMode = "active".equalsIgnoreCase(transport_mode) ? "TCP-ACTIVE" : "TCP-PASSIVE";
                        } else {
                            streamMode = device.getMediaTransport().toUpperCase();
                        }

                        StringBuilder content = new StringBuilder(256);
                        content.append("v=0\r\n");
                        content.append("o=")
                                .append(channelId)
                                .append(" 0 0 IN IP4 ")
                                .append(sipConfig.getIp())
                                .append("\r\n");
                        content.append("s=Play\r\n");
                        content.append("c=IN IP4 ")
                                .append(sipConfig.getIp())
                                .append("\r\n");
                        content.append("t=0 0\r\n");

                        // tcp被动模式
                        // tcp主动模式
                        if (isSeniorSdp) {
                            switch (streamMode) {
                                case "TCP-PASSIVE":
                                case "TCP-ACTIVE":
                                    content.append("m=video ")
                                            .append(smsConfig.getRtp().getMux_port())
                                            .append(" TCP/RTP/AVP 96 126 125 99 34 98 97\r\n");
                                    break;
                                case "UDP":
                                    content.append("m=video ")
                                            .append(smsConfig.getRtp().getMux_port())
                                            .append(" RTP/AVP 96 126 125 99 34 98 97\r\n");
                                    break;
                            }
                            content.append("a=recvonly\r\n");
                            content.append("a=rtpmap:96 PS/90000\r\n");
                            content.append("a=fmtp:126 profile-level-id=42e01e\r\n");
                            content.append("a=rtpmap:126 H264/90000\r\n");
                            content.append("a=rtpmap:125 H264S/90000\r\n");
                            content.append("a=fmtp:125 profile-level-id=42e01e\r\n");
                            content.append("a=rtpmap:99 MP4V-ES/90000\r\n");
                            content.append("a=fmtp:99 profile-level-id=3\r\n");
                        } else {
                            switch (streamMode) {
                                case "TCP-PASSIVE":
                                case "TCP-ACTIVE":
                                    content.append("m=video ")
                                            .append(smsConfig.getRtp().getMux_port())
                                            .append(" TCP/RTP/AVP 96 98 97\r\n");
                                    break;
                                case "UDP":
                                    content.append("m=video ")
                                            .append(smsConfig.getRtp().getMux_port())
                                            .append(" RTP/AVP 96 98 97\r\n");
                                    break;
                            }
                            content.append("a=recvonly\r\n");
                            content.append("a=rtpmap:96 PS/90000\r\n");
                        }
                        content.append("a=rtpmap:98 H264/90000\r\n");
                        content.append("a=rtpmap:97 MPEG4/90000\r\n");
                        if ("TCP-PASSIVE".equals(streamMode)) { // tcp被动模式
                            content.append("a=setup:passive\r\n");
                            content.append("a=connection:new\r\n");
                        } else if ("TCP-ACTIVE".equals(streamMode)) { // tcp主动模式
                            content.append("a=setup:active\r\n");
                            content.append("a=connection:new\r\n");
                        }
                        content.append("y=")
                                .append(srsCreateChannelResp.getSsrc())
                                .append("\r\n");//ssrc
                        String tm = Long.toString(System.currentTimeMillis());

                        CallIdHeader callIdHeader = "TCP".equals(streamMode) ? sipTcpProvider.getNewCallId()
                                : sipUdpProvider.getNewCallId();
                        Request request;

                        try {
                            SipURI sipURI = sipFactory
                                    .createAddressFactory()
                                    .createSipURI(channelId,
                                            device.getRemoteIP()+":"+device.getRemotePort());

                            ArrayList<ViaHeader> viaHeaders = new ArrayList<>();
                            ViaHeader viaHeader = sipFactory.createHeaderFactory()
                                    .createViaHeader(sipConfig.getIp(), sipConfig.getPort(), streamMode, null);
                            viaHeader.setRPort();
                            viaHeaders.add(viaHeader);

                            SipURI fromSipURI = sipFactory
                                    .createAddressFactory()
                                    .createSipURI(sipConfig.getSerial(), sipConfig.getIp()+":"+ sipConfig.getPort());
                            Address fromAddress = sipFactory
                                    .createAddressFactory()
                                    .createAddress(fromSipURI);
                            FromHeader fromHeader = sipFactory
                                    .createHeaderFactory()
                                    .createFromHeader(fromAddress, "FromInvt" + tm);

                            SipURI toSipURI = sipFactory
                                    .createAddressFactory()
                                    .createSipURI(channelId, sipConfig.getRealm());
                            Address toAddress = sipFactory
                                    .createAddressFactory()
                                    .createAddress(toSipURI);
                            ToHeader toHeader = sipFactory
                                    .createHeaderFactory()
                                    .createToHeader(toAddress, null);

                            MaxForwardsHeader maxForwards = sipFactory
                                    .createHeaderFactory()
                                    .createMaxForwardsHeader(70);

                            CSeqHeader cSeqHeader = sipFactory
                                    .createHeaderFactory()
                                    .createCSeqHeader(1L, Request.INVITE);

                            ContentTypeHeader contentTypeHeader = sipFactory
                                    .createHeaderFactory()
                                    .createContentTypeHeader("APPLICATION", "SDP");

                            Address contactAddress = sipFactory
                                    .createAddressFactory()
                                    .createAddress(sipFactory
                                            .createAddressFactory()
                                            .createSipURI(sipConfig.getSerial(), sipConfig.getIp() + ":" + sipConfig.getPort())
                                    );

                            request = sipFactory
                                    .createMessageFactory()
                                    .createRequest(sipURI, Request.INVITE, callIdHeader,
                                            cSeqHeader, fromHeader, toHeader, viaHeaders,
                                            maxForwards, contentTypeHeader, content.toString());
                            request.addHeader(sipFactory.createHeaderFactory().createContactHeader(contactAddress));

                            SubjectHeader subjectHeader = sipFactory
                                    .createHeaderFactory()
                                    .createSubjectHeader(String.format("%s:%s,%s:%s", channelId,
                                            srsCreateChannelResp.getSsrc(), sipConfig.getSerial(), 0));
                            request.addHeader(subjectHeader);


                            ClientTransaction clientTransaction = null;
                            if ("TCP".equals(streamMode)) {
                                clientTransaction = sipTcpProvider.getNewClientTransaction(request);
                            } else if ("UDP".equals(streamMode)) {
                                clientTransaction = sipUdpProvider.getNewClientTransaction(request);
                            }
                            if (clientTransaction != null){
                                clientTransaction.sendRequest();
                                log.info("invite {} request：\n{}", streamMode, request.toString());
                            }
                        } catch (SipException | ParseException | InvalidArgumentException e) {
                            log.error("Failed to send invite request", e);
                            return Mono.error(e);
                        }
                        return Mono.fromFuture(streamPlayResponseSubscribe)
                                .timeout(Duration.ofSeconds(timeout == null ? 15 : timeout))
                                .doFinally(signalType -> sipResponseHolder.remove(streamPlayResponseSubscribe.getKey(),
                                        streamPlayResponseSubscribe.getId()));
                    } else return Mono.error(new Exception(jsonObject.toJSONString()));
                });
    }

    @Override
    public void streamByeCmd(String deviceId, String channelId) {
        InviteResponseEvent inviteResponseEvent = deviceInterface.getInviteResponseEvent(deviceId,channelId);
        ClientTransaction transaction = inviteResponseEvent.getResponseEvent().getClientTransaction();
        SIPDialog dialog = (SIPDialog) inviteResponseEvent.getDialog();
        if (transaction == null) {
            log.info("transaction lost on streamByeCmd deviceId:{},channelId:{}",deviceId,channelId);
            return;
        }
        if(dialog == null){
            log.info("dialog lost on streamByeCmd deviceId:{},channelId:{}",deviceId,channelId);
            return;
        }

        try {
            Request byeRequest = dialog.createRequest(Request.BYE);
            SipURI byeURI = (SipURI) byeRequest.getRequestURI();
            SIPRequest request = (SIPRequest)transaction.getRequest();
            byeURI.setHost(request.getRemoteAddress().getHostName());
            byeURI.setPort(request.getRemotePort());
            ViaHeader viaHeader = (ViaHeader) byeRequest.getHeader(ViaHeader.NAME);
            String protocol = viaHeader.getTransport().toUpperCase();
            ClientTransaction clientTransaction = null;
            if("TCP".equals(protocol)) {
                clientTransaction = sipTcpProvider.getNewClientTransaction(byeRequest);
            } else if("UDP".equals(protocol)) {
                clientTransaction = sipUdpProvider.getNewClientTransaction(byeRequest);
            }
            dialog.sendRequest(clientTransaction);

            String url = String.format(DELETE_CHANNEL_URL,
                    smsConfig.getHost(), smsConfig.getHttpPort(),
                    deviceId,channelId);
            WebClient.create(url)
                    .get()
                    .retrieve()
                    .bodyToMono(String.class)
                    .subscribe(s -> log.info("Notify SRS to delete media channel result : {}",s));
        } catch (SipException | ParseException e) {
            log.error("stop stream error",e);
        }
    }

    @Override
    public Mono<List<DeviceChannel>> catalogQuery(Device device, Integer sn, Integer timeout, SipCatalogResponseSubscribe catalogResponseSubscribe) {
        StringBuilder catalogXml = new StringBuilder(200);
        catalogXml.append("<?xml version=\"1.0\" encoding=\"GB2312\"?>\r\n");
        catalogXml.append("<Query>\r\n");
        catalogXml.append("<CmdType>Catalog</CmdType>\r\n");
        catalogXml.append("<SN>")
                .append(sn)
                .append("</SN>\r\n");
        catalogXml.append("<DeviceID>")
                .append(device.getId())
                .append("</DeviceID>\r\n");
        catalogXml.append("</Query>\r\n");

        String tm = Long.toString(System.currentTimeMillis());
        CallIdHeader callIdHeader = "TCP".equals(device.getMediaTransport()) ? sipTcpProvider.getNewCallId()
                : sipUdpProvider.getNewCallId();
        try {
            SipURI requestURI = sipFactory.createAddressFactory()
                    .createSipURI(device.getId(), device.getRemoteIP()+":"+device.getRemotePort());

            ArrayList<ViaHeader> viaHeaders = new ArrayList<>();
            ViaHeader viaHeader = sipFactory.createHeaderFactory()
                    .createViaHeader(
                            sipConfig.getIp(),
                            sipConfig.getPort(),
                            device.getMediaTransport(),
                            "z9hG4bK" + tm);
            viaHeader.setRPort();
            viaHeaders.add(viaHeader);

            SipURI fromSipURI = sipFactory.createAddressFactory().createSipURI(sipConfig.getSerial(),
                    sipConfig.getRealm());
            Address fromAddress = sipFactory.createAddressFactory().createAddress(fromSipURI);
            FromHeader fromHeader = sipFactory.createHeaderFactory().createFromHeader(fromAddress, "FromCat" + tm);

            SipURI toSipURI = sipFactory.createAddressFactory().createSipURI(device.getId(), sipConfig.getRealm());
            Address toAddress = sipFactory.createAddressFactory().createAddress(toSipURI);
            ToHeader toHeader = sipFactory.createHeaderFactory().createToHeader(toAddress, null);

            MaxForwardsHeader maxForwards = sipFactory.createHeaderFactory().createMaxForwardsHeader(70);

            CSeqHeader cSeqHeader = sipFactory.createHeaderFactory().createCSeqHeader(1L, Request.MESSAGE);

            Request request = sipFactory.createMessageFactory().createRequest(requestURI, Request.MESSAGE, callIdHeader, cSeqHeader, fromHeader,
                    toHeader, viaHeaders, maxForwards);
            ContentTypeHeader contentTypeHeader = sipFactory.createHeaderFactory().createContentTypeHeader("Application", "MANSCDP+xml");
            request.setContent(catalogXml.toString(), contentTypeHeader);


            ClientTransaction clientTransaction = null;
            if ("TCP".equals(device.getCommandTransport())) {
                clientTransaction = sipTcpProvider.getNewClientTransaction(request);
            } else if ("UDP".equals(device.getCommandTransport())) {
                clientTransaction = sipUdpProvider.getNewClientTransaction(request);
            }
            if (clientTransaction != null){
                clientTransaction.sendRequest();
                log.info("catalogQuery request：\n{}", request.toString());
            }
        } catch (InvalidArgumentException | ParseException | SipException e) {
            e.printStackTrace();
        }

        return Mono.fromFuture(catalogResponseSubscribe)
                .timeout(Duration.ofSeconds(timeout == null ? 15 : timeout))
                .doFinally(signalType -> sipResponseHolder.remove(catalogResponseSubscribe.getKey(), catalogResponseSubscribe.getId()));
    }
}
