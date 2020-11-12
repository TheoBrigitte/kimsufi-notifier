FROM chromedp/headless-shell:latest

RUN apt-get update && apt-get install -y ca-certificates
ENTRYPOINT []

ADD ./kimsufi-notifier /app/kimsufi-notifier

CMD ["/app/kimsufi-notifier"]
