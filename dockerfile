FROM euleros:2.2

MAINTAINER shoothzj@gmail.com

COPY bin /opt/shooothzj/

RUN chmod -R 777 /opt/shooothzj

CMD ["/opt/shooothzj/scripts/start.sh"]
ENTRYPOINT ["/bin/c-init"]