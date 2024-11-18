FROM golang:1.21.8

WORKDIR /usr/share/bash-completion/completions

RUN curl -O https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh
RUN curl -O https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash

RUN chmod a+x git*.*

RUN ls -l $PWD/git*.* | awk '{print "source "$9}' >> ~/.bashrc

RUN echo "GIT_PS1_SHOWDIRTYSTATE=true" >> ~/.bashrc
RUN echo "GIT_PS1_SHOWUNTRACKEDFILES=true" >> ~/.bashrc
RUN echo "GIT_PS1_SHOWUPSTREAM=auto" >> ~/.bashrc

WORKDIR /workspace

CMD ["bin/bash"]