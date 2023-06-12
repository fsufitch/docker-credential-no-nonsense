###########
# DEVCONTAINER
##########
FROM fedora:37 AS devcontainer

# Install stuff necessary for a reasonable CLI
COPY devcontainer-packages.txt /opt
RUN dnf install -y $(cat /opt/devcontainer-packages.txt) && \
    rm -rf /var/cache/dnf && \
    npm i -g node npm

# Set up the devcontainer user
RUN useradd -ms /bin/bash developer
RUN echo "developer ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/developer
WORKDIR /home/developer
USER developer

COPY ./.bashrc.d /home/developer/.bashrc.d

COPY ./.devcontainer-install-go-tools.sh /opt/go-tools.sh
RUN bash /opt/go-tools.sh

ENV PATH=/home/developer/go/bin:${PATH}

CMD ["echo", "devcontainer should have its command overridden by the IDE"]

