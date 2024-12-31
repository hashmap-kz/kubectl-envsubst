package cmd

import (
	"os"
	"testing"
)

func TestEnvsubst(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		allowedVars     []string
		allowedPrefixes []string
		strict          bool
		env             map[string]string
		want            string
		expectError     bool
	}{
		{
			name:        "Basic substitution",
			input:       "Hello $USER",
			allowedVars: []string{"USER"},
			strict:      false,
			env:         map[string]string{"USER": "Alice"},
			want:        "Hello Alice",
			expectError: false,
		},
		{
			name:        "Unallowed variable",
			input:       "Hello $USER",
			allowedVars: []string{"NOT_USER"},
			strict:      false,
			env:         map[string]string{"USER": "Alice"},
			want:        "Hello $USER",
			expectError: false,
		},
		{
			name:        "Strict mode with undefined variable",
			input:       "Hello $USER",
			allowedVars: []string{},
			strict:      true,
			env:         map[string]string{},
			want:        "",
			expectError: true,
		},
		{
			name:            "Allowed prefix substitution",
			input:           "Hello $APP_USER",
			allowedVars:     []string{},
			allowedPrefixes: []string{"APP_"},
			strict:          false,
			env:             map[string]string{"APP_USER": "Bob"},
			want:            "Hello Bob",
			expectError:     false,
		},
		{
			name:        "Undefined variable with strict mode",
			input:       "Hello $USER",
			allowedVars: []string{"USER"},
			strict:      true,
			env:         map[string]string{},
			want:        "",
			expectError: true,
		},
		{
			name:        "Multiple variables",
			input:       "Hello $USER, welcome to $APP_ENV",
			allowedVars: []string{"USER", "APP_ENV"},
			strict:      false,
			env:         map[string]string{"USER": "Charlie", "APP_ENV": "production"},
			want:        "Hello Charlie, welcome to production",
			expectError: false,
		},
		{
			name:        "Partially resolved variables",
			input:       "Hello $USER, welcome to $APP_ENV",
			allowedVars: []string{"USER"},
			strict:      false,
			env:         map[string]string{"USER": "Charlie"},
			want:        "Hello Charlie, welcome to $APP_ENV",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			// Cleanup environment variables
			defer func() {
				for k := range tt.env {
					os.Unsetenv(k)
				}
			}()

			// Create Envsubst instance
			envsubst := NewEnvsubst(tt.allowedVars, tt.allowedPrefixes, tt.strict)
			result, err := envsubst.SubstituteEnvs(tt.input)

			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectError)
			}

			if result != tt.want {
				t.Errorf("unexpected result: got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestRegex(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"$VAR", []string{"VAR"}},
		{"${VAR}", []string{"VAR"}},
		{"$1VAR", nil},
		{"$(VAR)", nil},
		{"${VAR1_2}", []string{"VAR1_2"}},
		{"No variable here", nil},
	}

	for _, tc := range testCases {
		matches := envVarRegex.FindAllStringSubmatch(tc.input, -1)
		var result []string
		for _, match := range matches {
			if len(match) > 1 {
				result = append(result, match[1])
			}
		}

		if len(result) != len(tc.expected) {
			t.Errorf("For input '%s', expected %v, got %v", tc.input, tc.expected, result)
			continue
		}
		for i, v := range result {
			if v != tc.expected[i] {
				t.Errorf("For input '%s', expected %v, got %v", tc.input, tc.expected, result)
				break
			}
		}
	}
}

func TestStrictMode(t *testing.T) {
	os.Setenv("USER", "Alice")
	defer os.Unsetenv("USER")

	testCases := []struct {
		name        string
		input       string
		allowedVars []string
		strict      bool
		expected    string
		expectError bool
	}{
		{
			name:        "Strict mode with all variables resolved",
			input:       "Hello $USER",
			allowedVars: []string{"USER"},
			strict:      true,
			expected:    "Hello Alice",
			expectError: false,
		},
		{
			name:        "Strict mode with unresolved variable",
			input:       "Hello $UNKNOWN",
			allowedVars: []string{"USER"},
			strict:      true,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Strict mode with partially resolved variables",
			input:       "Hello $USER and $UNKNOWN",
			allowedVars: []string{"USER"},
			strict:      true,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Strict mode with no variables",
			input:       "Hello World",
			allowedVars: []string{},
			strict:      true,
			expected:    "Hello World",
			expectError: false,
		},
		{
			name:        "Non-strict mode with unresolved variable",
			input:       "Hello $UNKNOWN",
			allowedVars: []string{"USER"},
			strict:      false,
			expected:    "Hello $UNKNOWN",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			envsubst := NewEnvsubst(tc.allowedVars, nil, tc.strict)
			output, err := envsubst.SubstituteEnvs(tc.input)

			if (err != nil) != tc.expectError {
				t.Errorf("Unexpected error status for input '%s': got %v, want error=%v", tc.input, err, tc.expectError)
			}
			if output != tc.expected {
				t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, output)
			}
		})
	}
}

func TestIsFromAllowedList(t *testing.T) {
	envsubst := NewEnvsubst([]string{"USER"}, nil, false)

	if !envsubst.isFromAllowedList("USER") {
		t.Errorf("expected USER to be in the allowed list")
	}

	if envsubst.isFromAllowedList("OTHER") {
		t.Errorf("expected OTHER not to be in the allowed list")
	}
}

func TestIsFromPrefixList(t *testing.T) {
	envsubst := NewEnvsubst(nil, []string{"APP_"}, false)

	if !envsubst.isFromPrefixList("APP_USER") {
		t.Errorf("expected APP_USER to be matched by prefix APP_")
	}

	if envsubst.isFromPrefixList("USER") {
		t.Errorf("expected USER not to be matched by prefix APP_")
	}
}

func TestGetEnvValue(t *testing.T) {
	envsubst := NewEnvsubst([]string{"USER"}, []string{"APP_"}, false)
	os.Setenv("USER", "Alice")
	os.Setenv("APP_USER", "Bob")
	defer func() {
		os.Unsetenv("USER")
		os.Unsetenv("APP_USER")
	}()

	value, ok := envsubst.getEnvValue("USER")
	if !ok || value != "Alice" {
		t.Errorf("unexpected value for USER: got %q, want %q", value, "Alice")
	}

	value, ok = envsubst.getEnvValue("APP_USER")
	if !ok || value != "Bob" {
		t.Errorf("unexpected value for APP_USER: got %q, want %q", value, "Bob")
	}

	value, ok = envsubst.getEnvValue("UNSET_VAR")
	if ok {
		t.Errorf("expected UNSET_VAR to not be found")
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

	defer func() {
		for _, e := range []string{
			"request_uri",
			"server_name",
			"remote_addr",
			"host",
		} {
			os.Unsetenv(e)
		}
	}()

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
			envsubst := NewEnvsubst(test.allowedEnvs, []string{}, false)
			output, err := envsubst.SubstituteEnvs(test.text)

			if err != nil {
				t.Errorf("Unexpected error status for complex test")
			}

			if output != test.expected {
				t.Errorf("Test %q failed: expected %q, got %q", test.name, test.expected, output)
			}
		})
	}
}
