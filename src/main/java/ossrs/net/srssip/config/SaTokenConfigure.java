package ossrs.net.srssip.config;

import cn.dev33.satoken.context.SaHolder;
import cn.dev33.satoken.reactor.filter.SaReactorFilter;
import cn.dev33.satoken.router.SaHttpMethod;
import cn.dev33.satoken.router.SaRouter;
import cn.dev33.satoken.stp.StpUtil;
import com.alibaba.fastjson.JSON;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

/**
 * @ Description ossrs.net.srssip.config
 * @ Author StormBirds
 * @ Email xbaojun@gmail.com
 * @ Date 20/2/2022 下午10:39
 */
@Slf4j
@Configuration
public class SaTokenConfigure {

        /**
         * 注册 [Sa-Token全局过滤器]
         */
        @Bean
        public SaReactorFilter getSaReactorFilter() {
            return new SaReactorFilter()
                    // 指定 [拦截路由]
                    .addInclude("/**")
                    // 指定 [放行路由]
                    .addExclude("/favicon.ico")
                    // 指定[认证函数]: 每次请求执行
                    .setAuth(obj -> {
                        // 登录认证 -- 拦截所有路由，并排除/user/doLogin 用于开放登录
                        SaRouter.match("/**", "/user/login", StpUtil::checkLogin);
                    })
                    .setExcludeList(Arrays.asList(
                            "/user/login",
                            "/user/isLogin",
                            "/user/logout"
                    ))
                    // 异常处理函数：每次认证函数发生异常时执行此函数
                    .setError(e -> {
                        Map<String,Object> result = new HashMap<>();
                        result.put("code",0);
                        result.put("msg","failed");
                        result.put("data",e.getMessage());
                        return JSON.toJSONString(result);
                    })
                    // 前置函数：在每次认证函数之前执行
                    .setBeforeAuth(obj -> {
                        // ---------- 设置跨域响应头 ----------
                        SaHolder.getResponse()
                                // 允许指定域访问跨域资源
                                .setHeader("Access-Control-Allow-Origin", "*")
                                // 允许所有请求方式
                                .setHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
                                // 有效时间
                                .setHeader("Access-Control-Max-Age", "3600")
                                // 允许的header参数
                                .setHeader("Access-Control-Allow-Headers", "*");

                        // 如果是预检请求，则立即返回到前端
                        SaRouter.match(SaHttpMethod.OPTIONS)
                                .free(r -> log.info("--------OPTIONS预检请求，不做处理"))
                                .back();
                    })
                    ;
    }

}
