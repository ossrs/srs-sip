package ossrs.net.srssip.gb28181.annotation;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

/**
 * @ Description ossrs.net.srssip.gb28181.annotation
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 24/2/2022 上午1:04
 */
@Retention(RetentionPolicy.RUNTIME)
public @interface RequestEventHandler {
    String value();
}
