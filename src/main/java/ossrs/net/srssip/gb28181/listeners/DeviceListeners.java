package ossrs.net.srssip.gb28181.listeners;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.gb28181.event.messageevent.RegisterEvent;
import ossrs.net.srssip.gb28181.transaction.response.impl.SipMessageResponseHandler;

import javax.annotation.Resource;
import javax.sip.RequestEvent;
import javax.sip.header.HeaderFactory;
import javax.sip.header.WWWAuthenticateHeader;
import javax.sip.message.MessageFactory;
import javax.sip.message.Response;
import java.security.NoSuchAlgorithmException;
import java.text.ParseException;

import static gov.nist.javax.sip.clientauthutils.DigestServerAuthenticationHelper.DEFAULT_ALGORITHM;
import static gov.nist.javax.sip.clientauthutils.DigestServerAuthenticationHelper.DEFAULT_SCHEME;

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
    @EventListener
    public void deviceRegister(RegisterEvent registerEvent){
        RequestEvent requestEvent =  registerEvent.getRequestEvent();
        Response response;
        try {
            if (!registerEvent.doAuthenticatePlainTextPassword(sipConfig.getPassword())) {
                response = messageFactory.createResponse(Response.UNAUTHORIZED, requestEvent.getRequest());

                return;
            }
        } catch (NoSuchAlgorithmException | ParseException e) {
            e.printStackTrace();
        }
    }

    private void sendResponse(RequestEvent requestEvent, Response response) {
        sipMessageResponseHandler.sendResponse(requestEvent, response);
    }
}
