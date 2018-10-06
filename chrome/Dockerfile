FROM debian:jessie-slim
MAINTAINER Pierre Tachoire <pierre@tch.re>

RUN apt-get -qqy update \
    && apt-get -qqy install apt-transport-https wget

RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb https://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list \
    && apt-get -qqy update \
    && apt-get -qqy install google-chrome-stable

RUN groupadd -r chrome && useradd -r -g chrome -G audio,video chrome \
    && mkdir -p /home/chrome && chown -R chrome:chrome /home/chrome \
    && mkdir -p /opt/google/chrome \
    && chown -R chrome:chrome /opt/google/chrome
USER chrome

EXPOSE 9222
ENTRYPOINT ["google-chrome-stable", "--headless", "--disable-gpu", "--no-sandbox"]
CMD ["--remote-debugging-address=0.0.0.0", "--remote-debugging-port=9222"]
