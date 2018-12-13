FROM iron/base

COPY sesam-ftps-csv-source /opt/service/

WORKDIR /opt/service

RUN chmod +x /opt/service/sesam-ftps-csv-source

EXPOSE 8080:8080

CMD /opt/service/sesam-ftps-csv-source