package ossrs.net.srssip.gb28181.transaction.response;

import javax.sip.ClientTransaction;
import javax.sip.Dialog;
import javax.sip.RequestEvent;
import javax.sip.message.Request;
import javax.sip.message.Response;

/**
 * @ Description ossrs.net.srssip.gb28181.transaction.response
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午2:22
 */
public interface MessageResponseHandler {
    void sendResponse(RequestEvent requestEvent, Response response);

    ClientTransaction sendResponse(String transport, Request request);

    void sendDialog(Dialog dialog, Request request, String transport);
}
