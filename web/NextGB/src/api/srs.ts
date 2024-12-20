// @ts-nocheck
export function SrsRtcPlayer() {
  const self = {
    pc: new RTCPeerConnection(null),
    
    async play(url: string) {
      const conf = this.__internal.prepareUrl(url);
      this.pc.addTransceiver("audio", { direction: "recvonly" });
      this.pc.addTransceiver("video", { direction: "recvonly" });

      const offer = await this.pc.createOffer();
      await this.pc.setLocalDescription(offer);
      
      const session = await fetch(conf.apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          api: conf.apiUrl,
          streamurl: conf.streamUrl,
          clientip: null,
          sdp: offer.sdp
        })
      }).then(res => res.json());

      if (session.code) {
        throw session;
      }

      await this.pc.setRemoteDescription(
        new RTCSessionDescription({ type: 'answer', sdp: session.sdp })
      );
      return session;
    },

    close() {
      this.pc.close();
    },

    onaddstream: null,

    __internal: {
      defaultPath: '/rtc/v1/play/',
      
      prepareUrl(webrtcUrl: string) {
        const urlObject = this.parse(webrtcUrl);
        const schema = urlObject.user_query.schema 
          ? urlObject.user_query.schema + ':' 
          : window.location.protocol;

        let port = urlObject.port || 1985;
        if (schema === 'https:') {
          port = urlObject.port || 443;
        }

        let api = urlObject.user_query.play || this.defaultPath;
        if (api.lastIndexOf('/') !== api.length - 1) {
          api += '/';
        }

        let apiUrl = schema + '//' + urlObject.server + ':' + port + api;
        for (const key in urlObject.user_query) {
          if (key !== 'api' && key !== 'play') {
            apiUrl += '&' + key + '=' + urlObject.user_query[key];
          }
        }
        apiUrl = apiUrl.replace(api + '&', api + '?');

        return { 
          apiUrl, 
          streamUrl: urlObject.url,
          schema,
          urlObject,
          port
        };
      },

      parse(url: string) {
        const a = document.createElement("a");
        a.href = url.replace("rtmp://", "http://")
          .replace("webrtc://", "http://")
          .replace("rtc://", "http://");

        let vhost = a.hostname;
        let app = a.pathname.substring(1, a.pathname.lastIndexOf("/"));
        const stream = a.pathname.slice(a.pathname.lastIndexOf("/") + 1);

        app = app.replace("...vhost...", "?vhost=");
        if (app.indexOf("?") >= 0) {
          const params = app.slice(app.indexOf("?"));
          app = app.slice(0, app.indexOf("?"));

          if (params.indexOf("vhost=") > 0) {
            vhost = params.slice(params.indexOf("vhost=") + "vhost=".length);
            if (vhost.indexOf("&") > 0) {
              vhost = vhost.slice(0, vhost.indexOf("&"));
            }
          }
        }

        if (a.hostname === vhost) {
          const re = /^(\d+)\.(\d+)\.(\d+)\.(\d+)$/;
          if (re.test(a.hostname)) {
            vhost = "__defaultVhost__";
          }
        }

        let schema = "rtmp";
        if (url.indexOf("://") > 0) {
          schema = url.slice(0, url.indexOf("://"));
        }

        let port = parseInt(a.port);
        if (!port) {
          if (schema === 'http') {
            port = 80;
          } else if (schema === 'https') {
            port = 443;
          } else if (schema === 'rtmp') {
            port = 1935;
          }
        }

        const ret = {
          url,
          schema,
          server: a.hostname,
          port,
          vhost,
          app,
          stream,
          user_query: {}
        };

        this.fill_query(a.search, ret);
        return ret;
      },

      fill_query(query_string: string, obj: any) {
        if (query_string.length === 0) {
          return;
        }

        if (query_string.indexOf("?") >= 0) {
          query_string = query_string.split("?")[1];
        }

        const queries = query_string.split("&");
        for (const elem of queries) {
          const query = elem.split("=");
          obj[query[0]] = query[1];
          obj.user_query[query[0]] = query[1];
        }

        if (obj.domain) {
          obj.vhost = obj.domain;
        }
      }
    }
  };

  self.pc.onaddstream = (event) => {
    if (self.onaddstream) {
      self.onaddstream(event);
    }
  };

  return self;
}