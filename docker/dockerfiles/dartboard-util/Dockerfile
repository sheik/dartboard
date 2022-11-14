FROM ubuntu:22.10
RUN apt-get update && apt-get dist-upgrade -y && apt-get install ruby3.0 ruby3.0-dev ruby-rubygems gcc make rpm -y
RUN gem install fpm
WORKDIR /app
CMD ["/bin/bash"]
