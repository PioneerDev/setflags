FROM gitpod/workspace-postgres
USER root
RUN apt-get update -y
RUN apt-get install -y graphviz
USER gitpod
