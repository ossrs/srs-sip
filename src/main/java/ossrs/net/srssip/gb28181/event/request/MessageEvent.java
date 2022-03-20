package ossrs.net.srssip.gb28181.event.request;

import org.dom4j.DocumentException;
import org.dom4j.Element;
import ossrs.net.srssip.gb28181.annotation.RequestEventHandler;
import ossrs.net.srssip.gb28181.util.XmlUtil;

/**
 * @ Description ossrs.net.srssip.gb28181.event.request.messageevent
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 上午3:01
 */
@RequestEventHandler("MESSAGE")
public class MessageEvent extends RequestEventAbstract {

    private final static String CMD_TYPE = "CmdType";
    public String messageType;
    public String cmdType;

    @Override
    public void process() {
        Element element = null;
        try {
            element = XmlUtil.getRootElement(requestEvent.getRequest().getRawContent());
            this.messageType = element.getName();
            this.cmdType = XmlUtil.getText(element,CMD_TYPE);
        } catch (DocumentException e) {
            e.printStackTrace();
        }
    }
}
