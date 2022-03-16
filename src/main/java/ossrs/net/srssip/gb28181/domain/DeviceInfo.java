package ossrs.net.srssip.gb28181.domain;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * @ Description ossrs.net.srssip.gb28181.domain
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 14/3/2022 上午1:23
 */
@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class DeviceInfo implements Serializable {

    private String code;

    private String name;

    private String manufacturer;

    private String model;

    private String firmware;

    private Integer channelNum;

    private String host;

    private String transport;

    private Integer port;

    private String domain;

}
