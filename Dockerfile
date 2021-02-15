FROM centos:centos7

COPY ./mlcp /mlcp
ENTRYPOINT["/mlcp"]