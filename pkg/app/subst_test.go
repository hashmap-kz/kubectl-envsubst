package app

import (
	"os"
	"testing"
)

func TestSubstituteEnvs(t *testing.T) {
	// Set environment variables for testing
	_ = os.Setenv("FOO", "foo_value")
	_ = os.Setenv("BAR", "bar_value")
	_ = os.Setenv("EMPTY_VAR_VALUE", "")

	tests := []struct {
		name        string
		text        string
		allowedEnvs []string
		expected    string
	}{
		{
			name:        "Basic substitution with single variable",
			text:        "Hello $FOO!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value!",
		},
		{
			name:        "Basic substitution with multiple variables",
			text:        "Hello $FOO and ${BAR}!",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "Hello foo_value and bar_value!",
		},
		{
			name:        "Variable not in allowed list",
			text:        "Hello $FOO and $BAR!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value and $BAR!",
		},
		{
			name:        "No variables allowed",
			text:        "Hello $FOO!",
			allowedEnvs: []string{},
			expected:    "Hello $FOO!",
		},
		{
			name:        "Variable not set in environment",
			text:        "Hello $UNSET_VAR!",
			allowedEnvs: []string{"UNSET_VAR"},
			expected:    "Hello $UNSET_VAR!",
		},
		{
			name:        "No substitution needed",
			text:        "Hello world!",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "Hello world!",
		},
		{
			name:        "Empty text",
			text:        "",
			allowedEnvs: []string{"FOO", "BAR"},
			expected:    "",
		},
		{
			name:        "Partially valid and invalid variables",
			text:        "Hello ${FOO} and $BAZ!",
			allowedEnvs: []string{"FOO"},
			expected:    "Hello foo_value and $BAZ!",
		},

		{
			name:        "Edge cases",
			text:        "$var-name, $1VAR, ${var.name} $",
			allowedEnvs: []string{"var-name", "1VAR", "var.name"},
			expected:    "$var-name, $1VAR, ${var.name} $",
		},

		{
			name:        "Variable is set, but empty",
			text:        "Hello $EMPTY_VAR_VALUE!",
			allowedEnvs: []string{"EMPTY_VAR_VALUE"},
			expected:    "Hello !",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SubstituteEnvs(test.text, test.allowedEnvs)
			if result != test.expected {
				t.Errorf("Test %q failed: expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}

func TestComplexManifests(t *testing.T) {
	var complexMixedTest = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-https-config
  labels:
    app.kubernetes.io/name: nginx-https
    app.kubernetes.io/instance: nginx-https
data:
  nginx.conf: |
    user nginx;
    worker_processes 4;

    error_log /var/log/nginx/error.log warn;
    pid /var/run/nginx.pid;

    events {
        worker_connections 1024;
    }

    http {
        include       mime.types;
        default_type  application/octet-stream;

        log_format main '$remote_addr - $remote_user [$time_local] "$request" '
            '$status $body_bytes_sent "$http_referer" '
            '"$http_user_agent" "$http_x_forwarded_for"';

        log_format full	'$remote_addr - $host [$time_local] "$request" '
            'request_length=$request_length '
            'status=$status bytes_sent=$bytes_sent '
            'body_bytes_sent=$body_bytes_sent '
            'referer=$http_referer '
            'user_agent="$http_user_agent" '
            'upstream_status=$upstream_status '
            'request_time=$request_time '
            'upstream_response_time=$upstream_response_time '
            'upstream_connect_time=$upstream_connect_time '
            'upstream_header_time=$upstream_header_time';

        log_format json_combined escape=json
        '{'
            '"time_local":"$time_local",'
            '"remote_addr":"$remote_addr",'
            '"remote_user":"$remote_user",'
            '"request":"$request",'
            '"status": "$status",'
            '"body_bytes_sent":"$body_bytes_sent",'
            '"request_time":"$request_time",'
            '"http_referrer":"$http_referer",'
            '"http_user_agent":"$http_user_agent"'
        '}';

        log_format postdata '$remote_addr - $time_local - $request_body';

        access_log /var/log/nginx/access.log main;

        ######################################################################
        ## [Various settings]
        client_max_body_size 100M;
        client_body_buffer_size 512k;

        # copies data between one FD and other from within the kernel
        # faster than read() + write()
        sendfile on;

        # send headers in one piece, it is better than sending them one by one
        tcp_nopush on;

        server_tokens off;
        keepalive_timeout 65;
        types_hash_max_size 4096;


        ######################################################################
        ## [TLS settings]
        ssl_certificate     /etc/nginx/ssl/tls.crt;
        ssl_certificate_key /etc/nginx/ssl/tls.key;
        ssl_dhparam         /etc/nginx/dhparam/dhparam.pem; # Diffie-Hellman parameter for DHE ciphersuites, recommended 2048 bits

        # https://ssl-config.mozilla.org/
        ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA';
        ssl_prefer_server_ciphers on;

        # enable session resumption to improve https performance
        # http://vincent.bernat.im/en/blog/2011-ssl-session-reuse-rfc5077.html
        ssl_session_cache shared:SSL:50m;
        ssl_session_timeout 1d;
        ssl_session_tickets off;


        ######################################################################
        ## [Proxy settings]

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Server $host;
        proxy_read_timeout 5m;
        proxy_send_timeout 5m;
        proxy_connect_timeout 5m;
        #
        proxy_buffer_size 128k;
        proxy_buffers 4 256k;
        proxy_busy_buffers_size 256k;


        ######################################################################
        ## [Security settings]
        add_header X-Frame-Options SAMEORIGIN;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubdomains; preload";


        ######################################################################
        ## [Compressing settings]
        ##
        # reduce the data that needs to be sent over network -- for testing environment
        gzip on;
        # gzip_static on;
        gzip_min_length 10240;
        gzip_comp_level 1;
        gzip_vary on;
        gzip_disable msie6;
        gzip_proxied expired no-cache no-store private auth;
        gzip_types
            # text/html is always compressed by HttpGzipModule
            text/css
            text/javascript
            text/xml
            text/plain
            text/x-component
            application/javascript
            application/x-javascript
            application/json
            application/xml
            application/rss+xml
            application/atom+xml
            font/truetype
            font/opentype
            application/vnd.ms-fontobject
            image/svg+xml;

        include /etc/nginx/conf.d/*.conf;
    }

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: gateway-service-http
data:
  gateway-service-http.conf: |
    server {
        listen 80;
        server_name localhost;
        return 301 https://$server_name$request_uri;
        server_tokens off;
        access_log off;
        error_log off;
    }
    server {
        listen 443 ssl;
        server_name localhost;

        access_log /var/log/nginx/access.log json_combined;
        error_log  /var/log/nginx/error.log  warn;

        # Backend API endpoints
        #
        location /api/ {
          proxy_pass http://gateway-service-http:8080/api/;
        }

        # Swagger endpoints
        #
        location /swagger-ui/ {
            proxy_pass http://gateway-service-http:8080/swagger-ui/;
        }
        location /swagger-resources {
            proxy_pass http://gateway-service-http:8080/swagger-resources;
        }
        location /v3/api-docs {
            proxy_pass http://gateway-service-http:8080/v3/api-docs;
        }
    }

---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: nginx-tls
spec:
  refreshInterval: "15s"
  secretStoreRef:
    name: cluster-secret-store
    kind: ClusterSecretStore
  target:
    template:
      type: kubernetes.io/tls
      engineVersion: v2
  data:
    - secretKey: tls.crt
      remoteRef:
        key: secret/certs
        property: tls.crt
    - secretKey: tls.key
      remoteRef:
        key: secret/certs
        property: tls.key

---
apiVersion: v1
data:
  dhparam.pem: ""
kind: Secret
metadata:
  name: nginx-dhparam

---
apiVersion: v1
kind: Service
metadata:
  name: nginx-https
  labels:
    app.kubernetes.io/name: nginx-https
    app.kubernetes.io/instance: nginx-https
spec:
  type: NodePort
  ports:
    - name: https
      port: 443
      targetPort: 443
      nodePort: 30080
  selector:
    app.kubernetes.io/name: nginx-https
    app.kubernetes.io/instance: nginx-https

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nginx-https
  labels:
    app.kubernetes.io/name: nginx-https
    app.kubernetes.io/instance: nginx-https
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: nginx-https
      app.kubernetes.io/instance: nginx-https
  replicas: 1
  serviceName: nginx-https
  template:
    metadata:
      labels:
        app.kubernetes.io/name: nginx-https
        app.kubernetes.io/instance: nginx-https
    spec:
      containers:
        - name: nginx
          image: nginx:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 443
              name: https
          volumeMounts:
            # servers
            - mountPath: /etc/nginx/conf.d
              name: gateway-service-http
            # nginx.conf
            - mountPath: /etc/nginx/nginx.conf
              subPath: nginx.conf
              name: nginx-https-config
            # SSL certs
            - mountPath: /etc/nginx/ssl
              name: nginx-tls
            # DH param
            - name: nginx-dhparam
              mountPath: /etc/nginx/dhparam
          resources:
            requests:
              memory: "256Mi"
              cpu: "2m"
            limits:
              memory: "2Gi"
              cpu: "2"
      volumes:
        - name: nginx-tls
          secret:
            secretName: nginx-tls
        - name: nginx-https-config
          configMap:
            name: nginx-https-config
        - name: gateway-service-http
          configMap:
            name: gateway-service-http
            items:
              - key: gateway-service-http.conf
                path: gateway-service-http.conf
        - name: nginx-dhparam
          secret:
            secretName: nginx-dhparam

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-service-http
  labels:
    app: gateway-service-http
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gateway-service-http
  template:
    metadata:
      labels:
        app: gateway-service-http
    spec:
      containers:
        - image: my-registry/my-image:my-tag
          imagePullPolicy: IfNotPresent
          name: h8
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "256Mi"
              cpu: "2m"
            limits:
              memory: "2Gi"
              cpu: "2"

---
apiVersion: v1
kind: Service
metadata:
  name: gateway-service-http
  labels:
    app: gateway-service-http
spec:
  type: NodePort
  selector:
    app: gateway-service-http
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30081
`

	// Set environment variables for testing
	// especially we need to test when some nginx values are collide with env-vars,
	// they MUST not be expanded, unless explicitly set in allowed-vars
	_ = os.Setenv("request_uri", "1024")
	_ = os.Setenv("server_name", "2048")
	_ = os.Setenv("remote_addr", "2048")
	_ = os.Setenv("host", "2048")

	tests := []struct {
		name        string
		text        string
		allowedEnvs []string
		expected    string
	}{
		{
			name:        "Advanced substitution for a complex mixed manifest",
			text:        complexMixedTest,
			allowedEnvs: []string{},
			expected:    complexMixedTest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SubstituteEnvs(test.text, test.allowedEnvs)
			if result != test.expected {
				t.Errorf("Test %q failed: expected %q, got %q", test.name, test.expected, result)
			}
		})
	}
}
