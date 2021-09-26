FROM scratch

# Success apt
RUN apt install -y \
  bzr \
  cvs \
  git \
  mercurial \
  subversion
