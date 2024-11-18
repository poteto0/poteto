#!/bin/bash

cat <<EOF > ~/.bashrc

parse_git_branch() {
    git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/ (\1)/'
}

export PS1="\u@\h:\w\[\033[34m\]\$(parse_git_branch)\[\033[00m\] $ "
EOF