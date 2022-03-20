package ossrs.net.srssip.gb28181.util;

import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.io.SAXReader;
import org.springframework.util.StringUtils;

import javax.xml.bind.JAXBContext;
import javax.xml.bind.JAXBException;
import javax.xml.bind.Unmarshaller;
import java.io.ByteArrayInputStream;
import java.io.StringReader;
import java.util.HashMap;
import java.util.Map;

/**
 * @ Description ossrs.net.srssip.gb28181.util
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 9/3/2022 上午2:29
 */
public class XmlUtil {

    private final static String RN = "\r\n";

    public static Object xmlToObject(String xml, Object obj){
        try {
            JAXBContext context = JAXBContext.newInstance(obj.getClass());
            Unmarshaller unmarshaller = context.createUnmarshaller();
            StringReader sr = new StringReader(xml);
            return unmarshaller.unmarshal(sr);
        } catch (JAXBException e) {
            e.printStackTrace();
        }
        return null;
    }

    public static String getText(Element element, String tag) {
        if (null == element) {
            return null;
        }
        Element e = element.element(tag);
        //
        return null == e ? null : e.getText();
    }

    public static Element getRootElement(byte[] rawContent) throws DocumentException {
        SAXReader reader = new SAXReader();
        reader.setEncoding("gb2312");
        Document xml = reader.read(new ByteArrayInputStream(rawContent));
        return xml.getRootElement();
    }

    public static Map<String, String> convertStreamCode(String content) {
        if (!StringUtils.hasLength(content)) {
            return null;
        }
        Map<String, String> data = new HashMap<>();
        String [] values = content.split(RN);
        for (String value: values) {
            String[] fields = value.split("=");
            if (fields.length != 2) {
                continue;
            }
            data.put(fields[0], fields[1]);
        }
        return data;
    }
}
