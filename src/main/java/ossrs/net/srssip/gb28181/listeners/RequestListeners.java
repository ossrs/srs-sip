package ossrs.net.srssip.gb28181.listeners;

import gov.nist.javax.sip.address.AddressImpl;
import gov.nist.javax.sip.address.SipUri;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.gb28181.cmd.ISIPCommander;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.event.request.AckEvent;
import ossrs.net.srssip.gb28181.event.request.MessageEvent;
import ossrs.net.srssip.gb28181.event.request.RegisterEvent;
import ossrs.net.srssip.gb28181.event.request.message.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.event.response.InviteResponseEvent;
import ossrs.net.srssip.gb28181.event.subscribe.SipCatalogResponseSubscribe;
import ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import ossrs.net.srssip.gb28181.listeners.factory.MessageEventFactory;
import ossrs.net.srssip.gb28181.transaction.response.impl.SipMessageResponseHandler;
import ossrs.net.srssip.util.DigestServerAuthenticationHelper;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.header.ContactHeader;
import javax.sip.header.FromHeader;
import javax.sip.header.HeaderFactory;
import javax.sip.header.ViaHeader;
import javax.sip.message.MessageFactory;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.security.NoSuchAlgorithmException;
import java.text.ParseException;
import java.time.LocalDateTime;
import java.util.Calendar;
import java.util.Locale;

import static ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder.CALLBACK_CMD_CATALOG;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:49
 */
@Slf4j
@Component
public class RequestListeners {

    @Resource
    private SipConfig sipConfig;
    @Resource
    private MessageFactory messageFactory;
    @Resource
    private HeaderFactory headerFactory;
    @Resource
    private SipMessageResponseHandler sipMessageResponseHandler;
    @Resource
    private IDeviceInterface deviceInterface;
    @Resource
    private ApplicationEventPublisher applicationEventPublisher;
    @Resource
    private ISIPCommander sipCommander;
    @Resource
    private SipResponseHolder sipResponseHolder;

    @EventListener
    public void deviceRegister(RegisterEvent registerEvent) {
        RequestEvent requestEvent = registerEvent.getRequestEvent();
        Request request = requestEvent.getRequest();
        Response response;
        try {
            FromHeader fromHeader = (FromHeader) request.getHeader(FromHeader.NAME);
            AddressImpl address = (AddressImpl) fromHeader.getAddress();
            SipUri uri = (SipUri) address.getURI();
            String deviceId = uri.getUser();
            if (!registerEvent.doAuthenticatePlainTextPassword(sipConfig.getPassword())) {
                response = messageFactory.createResponse(Response.UNAUTHORIZED, request);
                new DigestServerAuthenticationHelper().generateChallenge(headerFactory, response, sipConfig.getRealm());
                sendResponse(requestEvent, response);
                return;
            }
            Device device = deviceInterface.getById(deviceId);
            LocalDateTime registerTime = LocalDateTime.now();
            ViaHeader viaHeader = (ViaHeader) request.getHeader(ViaHeader.NAME);
            String remoteIp = viaHeader.getReceived();
            int remotePort = viaHeader.getRPort();
            if (StringUtils.hasText(remoteIp) || remotePort == -1) {
                remoteIp = viaHeader.getHost();
                remotePort = viaHeader.getPort();
            }
            //检查设备是否已存在
            if (device == null) {
                device = new Device();
                device.setId(deviceId);
                device.setId(deviceId);
                device.setType("GB");
                device.setCatalogInterval(3600);
                device.setSubscribeInterval(0);
                device.setCatalogSubscribe(false);
                device.setAlarmSubscribe(false);
                device.setPositionSubscribe(false);
                device.setOnline(true);
                device.setPassword(sipConfig.getPassword());
                device.setCommandTransport(viaHeader.getTransport());
                device.setMediaTransport(viaHeader.getTransport());
                device.setRemoteIP(remoteIp);
                device.setRemotePort(remotePort);
                device.setLastRegisterAt(registerTime);
                device.setLastKeepaliveAt(registerTime);
                device.setUpdatedAt(registerTime);
                device.setCreatedAt(registerTime);
                deviceInterface.save(device);

            } else {
                device.setOnline(true);
                device.setRemoteIP(viaHeader.getHost());
                device.setRemotePort(viaHeader.getPort());
                device.setLastRegisterAt(registerTime);
                device.setLastKeepaliveAt(registerTime);
                device.setUpdatedAt(registerTime);
                deviceInterface.save(device);
            }

            response = messageFactory.createResponse(Response.OK, requestEvent.getRequest());
            response.addHeader(headerFactory.createDateHeader(Calendar.getInstance(Locale.ENGLISH)));
            response.addHeader(requestEvent.getRequest().getHeader(ContactHeader.NAME));
            response.addHeader(requestEvent.getRequest().getExpires());
            sendResponse(requestEvent, response);
            log.info("接收到注册请求 {}", registerEvent.toString());
            if (registerEvent.getExpires() <= 0) {
                //注销设备
                device.setOnline(false);
                device.setUpdatedAt(registerTime);
                deviceInterface.save(device);
            }else{
                int sn = (int) ((Math.random() * 9 + 1) * 10000000);
                String key = CALLBACK_CMD_CATALOG + deviceId;
                SipCatalogResponseSubscribe catalogResponseSubscribe = new SipCatalogResponseSubscribe(key, Integer.toString(sn));
                sipResponseHolder.put(key, Integer.toString(sn), catalogResponseSubscribe);
                sipCommander.catalogQuery(device, sn, sipConfig.getAckTimeout(), catalogResponseSubscribe);
            }
        } catch (NoSuchAlgorithmException | ParseException e) {
            e.printStackTrace();
        }
    }

    @EventListener
    public void ackMessage(AckEvent ackEvent) {
        RequestEvent requestEvent = ackEvent.getRequestEvent();
        Dialog dialog = requestEvent.getDialog();
        if (dialog == null) return;


    }

    @EventListener
    public void messageEvent(MessageEvent messageEvent) {
        RequestEvent requestEvent = messageEvent.getRequestEvent();
        MessageRequestAbstract messageRequestHandle = MessageEventFactory.INSTANCE.
                getMessageRequest(messageEvent.messageType, messageEvent.cmdType, requestEvent);
        applicationEventPublisher.publishEvent(messageRequestHandle);
    }

    private void sendResponse(RequestEvent requestEvent, Response response) {
        sipMessageResponseHandler.sendResponse(requestEvent, response);
    }
}
