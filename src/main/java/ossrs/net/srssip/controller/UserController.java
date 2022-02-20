package ossrs.net.srssip.controller;

import cn.dev33.satoken.stp.StpUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

/**
 * @ Description ossrs.net.srssip.controller
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:26
 */
@Slf4j
@RestController()
@RequestMapping("/user/")
public class UserController {

    @GetMapping("login")
    public String doLogin(@RequestParam String username,@RequestParam String password) {
        if("admin".equals(username) && "admin".equals(password)) {
            StpUtil.login(1);
            log.info("user token:{}",StpUtil.getTokenValue());
            return "登录成功";
        }
        return "登录失败";
    }

    @RequestMapping("isLogin")
    public String isLogin() {
        return "是否登录：" + StpUtil.isLogin();
    }

    @RequestMapping("logout")
    public String logout() {
        StpUtil.logout();
        return "ok";
    }
}
