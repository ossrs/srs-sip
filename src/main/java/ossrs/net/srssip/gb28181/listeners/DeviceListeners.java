package ossrs.net.srssip.gb28181.listeners;

import gov.nist.javax.sip.address.AddressImpl;
import gov.nist.javax.sip.address.SipUri;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.event.messageevent.MessageEvent;
import ossrs.net.srssip.gb28181.event.request.KeepLiveMessageRequest;
import ossrs.net.srssip.gb28181.event.messageevent.AckEvent;
import ossrs.net.srssip.gb28181.event.messageevent.RegisterEvent;
import ossrs.net.srssip.gb28181.event.request.MessageRequestAbstract;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import ossrs.net.srssip.gb28181.listeners.factory.MessageEventFactory;
import ossrs.net.srssip.gb28181.transaction.response.impl.SipMessageResponseHandler;
import ossrs.net.srssip.util.DigestServerAuthenticationHelper;

import javax.annotation.Resource;
import javax.sip.InvalidArgumentException;
import javax.sip.RequestEvent;
import javax.sip.SipException;
import javax.sip.header.ContactHeader;
import javax.sip.header.FromHeader;
import javax.sip.header.HeaderFactory;
import javax.sip.message.MessageFactory;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.security.NoSuchAlgorithmException;
import java.text.ParseException;
import java.util.Calendar;
import java.util.Locale;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:49
 */
@Slf4j
@Component
public class DeviceListeners {

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

    @EventListener
    public void deviceRegister(RegisterEvent registerEvent){
        RequestEvent requestEvent =  registerEvent.getRequestEvent();
        Request request = requestEvent.getRequest();
        Response response;
        try {
            FromHeader fromHeader = (FromHeader) request.getHeader(FromHeader.NAME);
            AddressImpl address = (AddressImpl) fromHeader.getAddress();
            SipUri uri = (SipUri) address.getURI();
            String deviceId = uri.getUser();
            Device device = deviceInterface.getById(deviceId);


            if (!registerEvent.doAuthenticatePlainTextPassword(sipConfig.getPassword())) {
                response = messageFactory.createResponse(Response.UNAUTHORIZED, request);
                new DigestServerAuthenticationHelper().generateChallenge(headerFactory, response, sipConfig.getRealm());
                sendResponse(requestEvent, response);
                return;
            }
            response = messageFactory.createResponse(Response.OK, requestEvent.getRequest());
            response.addHeader(headerFactory.createDateHeader(Calendar.getInstance(Locale.ENGLISH)));
            response.addHeader(requestEvent.getRequest().getHeader(ContactHeader.NAME));
            response.addHeader(requestEvent.getRequest().getExpires());
            sendResponse(requestEvent, response);
            log.info("接收到注册请求 {}",registerEvent.toString());
            if(registerEvent.getExpires()<=0){
                //注销设备
            }
            //检查设备是否已存在
        } catch (NoSuchAlgorithmException | ParseException e) {
            e.printStackTrace();
        }
    }

    @EventListener
    public void keepLiveEvent(KeepLiveMessageRequest keepLiveMessageRequest){
        RequestEvent requestEvent  = keepLiveMessageRequest.getRequestEvent();
        Request request = requestEvent.getRequest();

        try {
            Response response = messageFactory.createResponse(Response.OK, request);
            sendResponse(requestEvent,response);
        } catch (ParseException e) {
            e.printStackTrace();
        }
    }

    private void sendResponse(RequestEvent requestEvent, Response response) {
        sipMessageResponseHandler.sendResponse(requestEvent, response);
    }

    @EventListener
    public void ackMessage(AckEvent ackEvent){
        Request ackRequest = null;
        try {
            ackRequest = ackEvent.getDialog().createAck(ackEvent.getSeqNumber());
            ackEvent.getDialog().sendAck(ackRequest);
        } catch (InvalidArgumentException e) {
            e.printStackTrace();
        } catch (SipException e) {
            e.printStackTrace();
        }
    }

    @EventListener
    public void messageEvent(MessageEvent messageEvent){
        RequestEvent requestEvent = messageEvent.getRequestEvent();
        Response response;
        MessageRequestAbstract messageRequestAbstract = MessageEventFactory.INSTANCE.
                getMessageRequest(messageEvent.cmdType,requestEvent);
        applicationEventPublisher.publishEvent(messageRequestAbstract);
        try {
            response = messageFactory.createResponse(Response.OK, requestEvent.getRequest());
            sendResponse(requestEvent, response);
        } catch (ParseException e) {
            e.printStackTrace();
        }
    }
}
