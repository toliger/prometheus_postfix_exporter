
build:
	docker build -t prometheus_postfix_exporter .
	docker run --name prometheus_postfix_exporter prometheus_postfix_exporter /bin/true 2> /dev/null || echo "Its ok"
	docker cp prometheus_postfix_exporter:/bin/prometheus_postfix_exporter .
	docker rm prometheus_postfix_exporter

clean:
	[ -f "prometheus_postfix_exporter" ] && rm prometheus_postfix_exporter

test:
	go test

