package ossrs.net.srssip.gb28181.event.messageevent;

import gov.nist.javax.sip.address.SipUri;
import gov.nist.javax.sip.header.Expires;
import lombok.Getter;
import ossrs.net.srssip.gb28181.annotation.MessageEventHandler;
import ossrs.net.srssip.util.DigestServerAuthenticationHelper;

import javax.sip.header.AuthorizationHeader;
import javax.sip.header.ExpiresHeader;
import javax.sip.header.FromHeader;
import javax.sip.header.ViaHeader;
import javax.sip.message.Request;
import java.security.NoSuchAlgorithmException;
import java.util.Optional;

/**
 * @ Description ossrs.net.srssip.gb28181.event.messageevent
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:29
 */
@Getter
@MessageEventHandler(value = "REGISTER")
public class RegisterEvent extends MessageEventAbstract{

    private AuthorizationHeader authHeader;

    private String deviceId;

    private String host;

    private Integer port;

    private String transport;

    private Integer expires = 0;

    @Override
    public void process() {
        Request request = requestEvent.getRequest();
        this.authHeader = (AuthorizationHeader) request.getHeader(AuthorizationHeader.NAME);
        FromHeader fromHeader = (FromHeader) request.getHeader(FromHeader.NAME);
        ViaHeader viaHeader = (ViaHeader) request.getHeader(ViaHeader.NAME);
        SipUri uri = (SipUri)  fromHeader.getAddress().getURI();
        this.deviceId = uri.getUser();
        this.host = viaHeader.getHost();
        this.port = viaHeader.getPort();
        this.transport = viaHeader.getTransport();

        ExpiresHeader expiresHeader = (ExpiresHeader) request.getHeader(Expires.NAME);
        this.expires = Optional.of(expiresHeader.getExpires()).orElse(1);
    }

    public boolean doAuthenticatePlainTextPassword(String password) throws NoSuchAlgorithmException {
        Request request = requestEvent.getRequest();
        return new DigestServerAuthenticationHelper().doAuthenticatePlainTextPassword(request,password);
    }
}
