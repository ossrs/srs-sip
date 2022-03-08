package ossrs.net.srssip.gb28181.event.messageevent;

import org.dom4j.DocumentException;
import org.dom4j.Element;
import ossrs.net.srssip.gb28181.annotation.MessageEventHandler;
import ossrs.net.srssip.gb28181.util.XmlUtil;

/**
 * @ Description ossrs.net.srssip.gb28181.event.messageevent
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 上午3:01
 */
@MessageEventHandler("MESSAGE")
public class MessageEvent extends MessageEventAbstract{

    private final static String CMD_TYPE = "CmdType";
    public String cmdType;

    @Override
    public void process() {
        Element element = null;
        try {
            element = XmlUtil.getRootElement(requestEvent.getRequest().getRawContent());
        } catch (DocumentException e) {
            e.printStackTrace();
        }
        this.cmdType = XmlUtil.getText(element,CMD_TYPE);
    }
}
