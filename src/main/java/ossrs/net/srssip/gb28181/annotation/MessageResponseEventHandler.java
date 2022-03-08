package ossrs.net.srssip.gb28181.annotation;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * @ Description ossrs.net.srssip.gb28181.annotation
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:06
 */
@Retention(RetentionPolicy.RUNTIME)
public @interface MessageResponseEventHandler {
    String value();
}
