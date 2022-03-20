package ossrs.net.srssip.gb28181.listeners;

import gov.nist.javax.sip.address.AddressImpl;
import gov.nist.javax.sip.address.SipUri;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.gb28181.domain.Device;
import ossrs.net.srssip.gb28181.domain.DeviceChannel;
import ossrs.net.srssip.gb28181.domain.DeviceInfo;
import ossrs.net.srssip.gb28181.event.request.message.notify.KeepLiveMessageRequest;
import ossrs.net.srssip.gb28181.event.request.message.query.CatalogMessageQuery;
import ossrs.net.srssip.gb28181.event.request.message.response.CatalogMessageResponse;
import ossrs.net.srssip.gb28181.event.request.message.response.DeviceInfoMessageRequest;
import ossrs.net.srssip.gb28181.event.subscribe.SipCatalogResponseSubscribe;
import ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder;
import ossrs.net.srssip.gb28181.interfaces.IDeviceInterface;
import ossrs.net.srssip.gb28181.transaction.response.impl.SipMessageResponseHandler;

import javax.annotation.Resource;
import javax.sip.RequestEvent;
import javax.sip.header.FromHeader;
import javax.sip.message.MessageFactory;
import javax.sip.message.Request;
import javax.sip.message.Response;
import java.text.ParseException;
import java.util.List;

import static ossrs.net.srssip.gb28181.event.subscribe.SipResponseHolder.CALLBACK_CMD_CATALOG;

/**
 * @ Description ossrs.net.srssip.gb28181.listeners
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 19/3/2022 下午2:14
 */
@Slf4j
@Component
public class MessageListeners {
    @Resource
    private MessageFactory messageFactory;
    @Resource
    private SipResponseHolder sipResponseHolder;
    @Resource
    private IDeviceInterface deviceInterface;
    @Resource
    private SipMessageResponseHandler sipMessageResponseHandler;

    @EventListener
    public void deviceInfo(DeviceInfoMessageRequest deviceInfoMessageRequest) {
        DeviceInfo deviceInfo = deviceInfoMessageRequest.getDeviceInfo();
        Device device = deviceInterface.getById(deviceInfo.getCode());
        if (device != null) {
            device.setName(deviceInfo.getName());
            device.setManufacturer(deviceInfo.getManufacturer());
            deviceInterface.save(device);
        }
    }


    @EventListener
    public void keepLiveEvent(KeepLiveMessageRequest keepLiveMessageRequest) {
        RequestEvent requestEvent = keepLiveMessageRequest.getRequestEvent();
        Request request = requestEvent.getRequest();

        try {
            FromHeader fromHeader = (FromHeader) request.getHeader(FromHeader.NAME);
            AddressImpl address = (AddressImpl) fromHeader.getAddress();
            SipUri uri = (SipUri) address.getURI();
            String deviceId = uri.getUser();
            Device device = deviceInterface.getById(deviceId);
            Response response;
            if (device == null) {
                response = messageFactory.createResponse(Response.DOES_NOT_EXIST_ANYWHERE, request);
                sendResponse(requestEvent, response);
                return;
            }
            response = messageFactory.createResponse(Response.OK, request);
            sendResponse(requestEvent, response);
        } catch (ParseException e) {
            e.printStackTrace();
        }
    }

    @EventListener
    public void catologMessageResponse(CatalogMessageResponse catalogMessageResponse) {
        RequestEvent requestEvent = catalogMessageResponse.getRequestEvent();
        Request request = requestEvent.getRequest();

        try {
            Response response = messageFactory.createResponse(Response.OK, request);
            sendResponse(requestEvent, response);
        } catch (ParseException e) {
            e.printStackTrace();
        }
        List<DeviceChannel> deviceChannelList = catalogMessageResponse.getDeviceChannel();
        String key = CALLBACK_CMD_CATALOG + catalogMessageResponse.getDeviceId();
        String sn = catalogMessageResponse.getSn();
        SipCatalogResponseSubscribe catalogResponseSubscribe = (SipCatalogResponseSubscribe) sipResponseHolder.getCallable(key, sn);
        if (catalogResponseSubscribe.putChannelListVo(deviceChannelList) &&
                catalogResponseSubscribe.getChannelListVo().size() >= catalogMessageResponse.getSumNum()) {
            sipResponseHolder.complete(key, sn, catalogResponseSubscribe.getChannelListVo());
        }
        deviceInterface.saveDeviceChannel(deviceChannelList);
    }

    @EventListener()
    public void catalogMessageQuery(CatalogMessageQuery catalogMessageQuery) {
        RequestEvent requestEvent = catalogMessageQuery.getRequestEvent();
        Request request = requestEvent.getRequest();
        try {
            Response response = messageFactory.createResponse(Response.OK, request);
            sendResponse(requestEvent, response);
        } catch (ParseException e) {
            e.printStackTrace();
        }
    }

    private void sendResponse(RequestEvent requestEvent, Response response) {
        sipMessageResponseHandler.sendResponse(requestEvent, response);
    }
}
