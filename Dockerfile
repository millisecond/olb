FROM scratch
ADD build/ca-certificates.crt /etc/ssl/certs/
ADD olb.properties /etc/olb/olb.properties
ADD olb /
EXPOSE 9998 9999
CMD ["/olb", "-cfg", "/etc/olb/olb.properties"]
