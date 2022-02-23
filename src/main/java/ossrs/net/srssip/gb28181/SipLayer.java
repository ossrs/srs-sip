package ossrs.net.srssip.gb28181;

import gov.nist.javax.sip.SipStackImpl;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.DependsOn;
import org.springframework.stereotype.Component;
import ossrs.net.srssip.config.SipConfig;
import ossrs.net.srssip.gb28181.listeners.SipServerListeners;

import javax.annotation.Resource;
import javax.sip.*;
import javax.sip.address.AddressFactory;
import javax.sip.header.HeaderFactory;
import javax.sip.message.MessageFactory;
import java.util.Properties;
import java.util.TooManyListenersException;

/**
 * @ Description ossrs.net.srssip.gb28181
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 23/2/2022 下午10:34
 */
@Slf4j
@Configuration
public class SipLayer {

    @Resource
    private SipConfig sipConfig;

    @Resource
    private SipServerListeners sipServerListeners;

    @Bean
    public SipFactory sipFactory() {
        SipFactory sipFactory = SipFactory.getInstance();
        sipFactory.setPathName("gov.nist");
        return sipFactory;
    }

    @Bean
    @DependsOn({"sipFactory"})
    public HeaderFactory headerFactory(SipFactory sipFactory) throws PeerUnavailableException {
        return sipFactory.createHeaderFactory();
    }

    @Bean
    public AddressFactory addressFactory(SipFactory sipFactory) throws PeerUnavailableException {
        return sipFactory.createAddressFactory();
    }

    @Bean
    public MessageFactory messageFactory(SipFactory sipFactory) throws PeerUnavailableException {
        return sipFactory.createMessageFactory();
    }

    @Bean
    @DependsOn({"sipFactory"})
    public SipStack sipStack(SipFactory sipFactory) throws PeerUnavailableException {
        Properties properties = new Properties();
        properties.setProperty("javax.sip.STACK_NAME", sipConfig.getName());
        properties.setProperty("javax.sip.IP_ADDRESS", sipConfig.getIp());
        properties.setProperty("gov.nist.javax.sip.LOG_MESSAGE_CONTENT", sipConfig.getLogMessageContent());
        properties.setProperty("gov.nist.javax.sip.TRACE_LEVEL", sipConfig.getTraceLevel());
        properties.setProperty("gov.nist.javax.sip.SERVER_LOG", sipConfig.getName() + "_server_log");
        properties.setProperty("gov.nist.javax.sip.DEBUG_LOG", sipConfig.getName() + "_debug_log");
        return sipFactory.createSipStack(properties);
    }

    @Bean
    @DependsOn("sipStack")
    public SipProvider sipTcpProvider(SipStack sipStack) {
        ListeningPoint listeningPoint = null;
        SipProvider sipProvider = null;
        try {
            listeningPoint = sipStack.createListeningPoint(sipConfig.getIp(), sipConfig.getPort(), ListeningPoint.TCP);
            sipProvider = sipStack.createSipProvider(listeningPoint);
            sipProvider.addSipListener(sipServerListeners);
            log.info("tcp server {} is running on port {}.", listeningPoint.getIPAddress(), listeningPoint.getPort());
        } catch (TransportNotSupportedException e) {
            e.printStackTrace();
        } catch (InvalidArgumentException | ObjectInUseException e) {
            e.printStackTrace();
        } catch (TooManyListenersException e) {
            e.printStackTrace();
        }
        return sipProvider;
    }

    @Bean
    @DependsOn("sipStack")
    public SipProvider sipUdpProvider(SipStack sipStack) {
        ListeningPoint listeningPoint = null;
        SipProvider sipProvider = null;
        try {
            listeningPoint = sipStack.createListeningPoint(sipConfig.getIp(), sipConfig.getPort(), ListeningPoint.UDP);
            sipProvider = sipStack.createSipProvider(listeningPoint);
            sipProvider.addSipListener(sipServerListeners);
            log.info("udp server {} is running on port {}.", listeningPoint.getIPAddress(), listeningPoint.getPort());
        } catch (TransportNotSupportedException | InvalidArgumentException | ObjectInUseException | TooManyListenersException e) {
            e.printStackTrace();
        }
        return sipProvider;
    }
}
